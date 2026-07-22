package maprenderer

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/minetest-go/types"
)

type isoViewConfig struct {
	name         string
	probeX       int
	probeZ       int
	orderXWeight int
	orderZWeight int
}

func NormalizeIsoView(view string) string {
	switch strings.ToLower(view) {
	case "ne", "nw", "se", "sw":
		return strings.ToLower(view)
	default:
		return "ne"
	}
}

func getIsoViewConfig(view string) isoViewConfig {
	switch NormalizeIsoView(view) {
	case "nw":
		return isoViewConfig{
			name:         "nw",
			probeX:       1,
			probeZ:       -1,
			orderXWeight: -1,
			orderZWeight: 1,
		}
	case "sw":
		return isoViewConfig{
			name:         "sw",
			probeX:       -1,
			probeZ:       -1,
			orderXWeight: 1,
			orderZWeight: 1,
		}
	case "se":
		return isoViewConfig{
			name:         "se",
			probeX:       -1,
			probeZ:       1,
			orderXWeight: 1,
			orderZWeight: -1,
		}
	case "ne":
		fallthrough
	default:
		return isoViewConfig{
			name:         "ne",
			probeX:       1,
			probeZ:       1,
			orderXWeight: -1,
			orderZWeight: -1,
		}
	}
}

func GetVirtualSize(size *types.Pos, view string) *types.Pos {
	switch NormalizeIsoView(view) {
	case "nw", "se":
		return types.NewPos(size.Z(), size.Y(), size.X())
	default:
		return types.NewPos(size.X(), size.Y(), size.Z())
	}
}

func GetVirtualPos(rel_pos *types.Pos, size *types.Pos, view string) *types.Pos {
	switch NormalizeIsoView(view) {
	case "nw":
		return types.NewPos(size.Z()-1-rel_pos.Z(), rel_pos.Y(), rel_pos.X())
	case "sw":
		return types.NewPos(size.X()-1-rel_pos.X(), rel_pos.Y(), size.Z()-1-rel_pos.Z())
	case "se":
		return types.NewPos(rel_pos.Z(), rel_pos.Y(), size.X()-1-rel_pos.X())
	case "ne":
		fallthrough
	default:
		return types.NewPos(rel_pos.X(), rel_pos.Y(), rel_pos.Z())
	}
}

func GetIsometricImageSize(size *types.Pos, cube_len int) (int, int) {
	width := (size.X() * cube_len / 2) +
		(size.Z() * cube_len / 2)

	height := (size.X() * cube_len / 4) +
		(size.Y() * cube_len / 2) +
		(size.Z() * cube_len / 4)

	return width, height
}

func GetIsoCenterCubeOffset(size *types.Pos, cube_len int) (int, int) {
	x := (size.Z() * cube_len / 2) -
		(cube_len / 2)

	y := (size.X() * cube_len / 4) +
		(size.Y() * cube_len / 2) +
		(size.Z() * cube_len / 4) -
		cube_len

	return x, y
}

func GetIsoCubePosition(center_x, center_y, cube_len int, pos *types.Pos) (int, int) {
	x := center_x -
		(pos.Z() * cube_len / 2) +
		(pos.X() * cube_len / 2)

	y := center_y -
		(pos.X() * cube_len / 4) -
		(pos.Y() * cube_len / 2) -
		(pos.Z() * cube_len / 4)

	return x, y
}

func DrawIsoCube(img *image.RGBA, cube_len, x_offset, y_offset int, c1, c2, c3 color.Color) error {
	if cube_len%4 != 0 {
		return fmt.Errorf("cube_len must be divisible by 4")
	}
	if cube_len <= 4 {
		return fmt.Errorf("cube_len must be greater than 4")
	}

	half_len_zero_indexed := (cube_len / 2) - 1
	quarter_len := cube_len / 4

	// left/right part
	yo := 0
	for x := 0; x <= half_len_zero_indexed; x++ {
		for y := 0; y <= half_len_zero_indexed; y++ {
			img.Set(x_offset+x, y_offset+y+quarter_len+yo, c1)
			img.Set(x_offset+cube_len-1-x, y_offset+y+quarter_len+yo, c2)
		}
		if x%2 == 0 {
			yo = yo + 1
		}
	}

	// upper part
	yo = 0
	yl := 1
	for x := 0; x <= half_len_zero_indexed-1; x++ {
		for y := 0; y <= yl; y++ {
			img.Set(x_offset+1+x, y_offset+quarter_len-1-yo+y, c3)
			img.Set(x_offset+cube_len-2-x, y_offset+quarter_len-1-yo+y, c3)
		}
		if x%2 != 0 {
			yo = yo + 1
			yl = yl + 2
		}
	}

	return nil
}

func GetIsoNodeOrderForView(p *types.Pos, view string) int {
	cfg := getIsoViewConfig(view)
	return (64000 + p.Z()*cfg.orderZWeight) + p.Y() + (64000 + p.X()*cfg.orderXWeight)
}

func addAndClampUint8(a uint8, b int) uint8 {
	v := int(a) + b
	if v > 255 {
		return 255
	} else if v < 0 {
		return 0
	} else {
		return uint8(v)
	}
}

func ColorAdjust(c *color.RGBA, value int) *color.RGBA {
	return &color.RGBA{
		R: addAndClampUint8(c.R, value),
		G: addAndClampUint8(c.G, value),
		B: addAndClampUint8(c.B, value),
		A: c.A,
	}
}

func BlendColor(bg, fg *color.RGBA, bf float64) *color.RGBA {
	a := float64(fg.A) / 255 / bf
	ai := 1 - a

	return &color.RGBA{
		R: uint8((float64(fg.R) * a) + (float64(bg.R) * ai)),
		G: uint8((float64(fg.G) * a) + (float64(bg.G) * ai)),
		B: uint8((float64(fg.B) * a) + (float64(bg.B) * ai)),
		A: max(bg.A, fg.A),
	}
}

func SortNodesForView(view string) func(n1, n2 *NodeWithColor) int {
	return func(n1, n2 *NodeWithColor) int {
		o1 := GetIsoNodeOrderForView(n1.Pos, view)
		o2 := GetIsoNodeOrderForView(n2.Pos, view)
		return o1 - o2
	}
}
