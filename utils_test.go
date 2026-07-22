package maprenderer

import (
	"testing"

	"github.com/minetest-go/types"
)

func TestNormalizeIsoView(t *testing.T) {
	testcases := map[string]string{
		"NE":      "ne",
		"nw":      "nw",
		"Se":      "se",
		"sw":      "sw",
		"invalid": "ne",
	}

	for input, want := range testcases {
		if got := NormalizeIsoView(input); got != want {
			t.Fatalf("NormalizeIsoView(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestGetIsoViewConfigProbeDirections(t *testing.T) {
	testcases := map[string]struct {
		probeX int
		probeZ int
	}{
		"ne": {probeX: 1, probeZ: 1},
		"nw": {probeX: 1, probeZ: -1},
		"se": {probeX: -1, probeZ: 1},
		"sw": {probeX: -1, probeZ: -1},
	}

	for view, want := range testcases {
		cfg := getIsoViewConfig(view)
		if cfg.probeX != want.probeX || cfg.probeZ != want.probeZ {
			t.Fatalf("getIsoViewConfig(%q) = (%d,%d), want (%d,%d)", view, cfg.probeX, cfg.probeZ, want.probeX, want.probeZ)
		}
	}
}

func TestGetIsoCubePositionMatchesLegacyNEProjection(t *testing.T) {
	size := types.NewPos(5, 4, 3)
	vsize := GetVirtualSize(size, "ne")
	centerX, centerY := GetIsoCenterCubeOffset(vsize, 16)
	relPos := types.NewPos(4, 1, 2)
	vpos := GetVirtualPos(relPos, size, "ne")

	gotX, gotY := GetIsoCubePosition(centerX, centerY, 16, vpos)
	wantX := centerX - (relPos.Z() * 16 / 2) + (relPos.X() * 16 / 2)
	wantY := centerY - (relPos.X() * 16 / 4) - (relPos.Y() * 16 / 2) - (relPos.Z() * 16 / 4)

	if gotX != wantX || gotY != wantY {
		t.Fatalf("GetIsoCubePosition(ne) = (%d,%d), want (%d,%d)", gotX, gotY, wantX, wantY)
	}
}

func TestGetIsoNodeOrderForViewMatchesVirtualTransform(t *testing.T) {
	size := types.NewPos(6, 5, 4)
	points := []*types.Pos{
		types.NewPos(0, 0, 0),
		types.NewPos(1, 2, 0),
		types.NewPos(2, 1, 3),
		types.NewPos(5, 4, 2),
	}

	compare := func(a, b int) int {
		switch {
		case a < b:
			return -1
		case a > b:
			return 1
		default:
			return 0
		}
	}

	legacyOrder := func(p *types.Pos) int {
		return (64000 - p.X()) + p.Y() + (64000 - p.Z())
	}

	views := []string{"ne", "nw", "se", "sw"}
	for _, view := range views {
		for i := range points {
			for j := range points {
				got := compare(GetIsoNodeOrderForView(points[i], view), GetIsoNodeOrderForView(points[j], view))
				want := compare(
					legacyOrder(GetVirtualPos(points[i], size, view)),
					legacyOrder(GetVirtualPos(points[j], size, view)),
				)

				if got != want {
					t.Fatalf("view %q ordering mismatch for %v and %v: got %d want %d", view, points[i], points[j], got, want)
				}
			}
		}
	}
}
