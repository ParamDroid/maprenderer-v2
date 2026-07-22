package maprenderer

import (
	"fmt"
	"image"
	"slices"

	"github.com/minetest-go/types"
)

type IsoRenderOpts struct {
	CubeLen            int
	EnableTransparency bool
	View               string
}

func NewDefaultIsoRenderOpts() *IsoRenderOpts {
	return &IsoRenderOpts{
		CubeLen:            8,
		EnableTransparency: false,
		View:               "ne",
	}
}

func RenderIsometric(na types.NodeAccessor, cr types.ColorResolver, from, to *types.Pos, opts *IsoRenderOpts) (image.Image, error) {
	if opts == nil {
		opts = NewDefaultIsoRenderOpts()
	}

	opts.View = NormalizeIsoView(opts.View)

	min, max := types.SortPos(from, to)
	size := to.Subtract(from).Add(types.NewPos(1, 1, 1))
	vsize := GetVirtualSize(size, opts.View)

	width, height := GetIsometricImageSize(vsize, opts.CubeLen)
	center_x, center_y := GetIsoCenterCubeOffset(vsize, opts.CubeLen)
	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: width, Y: height},
	})

	skip_alpha := !opts.EnableTransparency

	view_cfg := getIsoViewConfig(opts.View)

	// Opposite/front-facing probe direction
	ipos := types.NewPos(view_cfg.probeX, -1, view_cfg.probeZ)

	nodes := []*NodeWithColor{}

	// top layer
	for x := min.X(); x <= max.X(); x++ {
		for z := min.Z(); z <= max.Z(); z++ {
			pnodes, err := Probe(min, max, types.NewPos(x, max.Y(), z), ipos, na, cr, skip_alpha)
			if err != nil {
				return nil, fmt.Errorf("probe error, top layer: %v", err)
			}
			nodes = append(nodes, pnodes...)
		}
	}

	// z-layer (right or left depending on view)
	z_start := max.Z()
	if view_cfg.probeZ == 1 {
		z_start = min.Z()
	}
	for x := min.X(); x <= max.X(); x++ {
		for y := min.Y(); y <= max.Y()-1; y++ {
			pnodes, err := Probe(min, max, types.NewPos(x, y, z_start), ipos, na, cr, skip_alpha)
			if err != nil {
				return nil, fmt.Errorf("probe error, z layer: %v", err)
			}
			nodes = append(nodes, pnodes...)
		}
	}

	// x-layer
	x_start := max.X()
	if view_cfg.probeX == 1 {
		x_start = min.X()
	}
	// Avoid double-probing the corner pillar
	z_min := min.Z()
	z_max := max.Z()
	if view_cfg.probeZ == -1 {
		z_max--
	} else {
		z_min++
	}
	for z := z_min; z <= z_max; z++ {
		for y := min.Y(); y <= max.Y()-1; y++ {
			pnodes, err := Probe(min, max, types.NewPos(x_start, y, z), ipos, na, cr, skip_alpha)
			if err != nil {
				return nil, fmt.Errorf("probe error, x layer: %v", err)
			}
			nodes = append(nodes, pnodes...)
		}
	}

	slices.SortFunc(nodes, SortNodesForView(opts.View))

	for _, n := range nodes {
		rel_pos := n.Pos.Subtract(min)
		vpos := GetVirtualPos(rel_pos, size, opts.View)

		c1 := ColorAdjust(n.Color, 0)
		if !opts.EnableTransparency {
			c1.A = 255
		}

		c2 := ColorAdjust(c1, -10)
		c3 := ColorAdjust(c1, 10)

		x, y := GetIsoCubePosition(center_x, center_y, opts.CubeLen, vpos)
		err := DrawIsoCube(img, opts.CubeLen, x, y, c1, c2, c3)
		if err != nil {
			return nil, fmt.Errorf("DrawIsoCube error: %v", err)
		}
	}

	return img, nil
}
