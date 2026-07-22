package maprenderer

import (
	"fmt"
	"strings"

	"github.com/minetest-go/types"
)

func Probe(min, max, pos, ipos *types.Pos, na types.NodeAccessor, cr types.ColorResolver, skip_alpha bool) ([]*NodeWithColor, error) {
	nodes := []*NodeWithColor{}
	// Comment out the node names below that you want to keep.
	// By default, some nodes are skipped to produce a cleaner render.
	cpos := pos
	for cpos.IsWithin(min, max) {
		node, err := na(cpos)
		if err != nil {
			return nil, fmt.Errorf("getNode error @ %s: %v", cpos, err)
		}

		if node != nil &&
			node.Name != "air" &&
			node.Name != "ignore" &&

// Add additional node names here before compiling if you want to exclude them.
// Here are some of the nodes I found creating mess while taking a picture of the builds.

			//node.Name != "default:aspen_tree" &&
			//node.Name != "default:aspen_leaves" &&
			//node.Name != "default:acacia_tree" &&
			//node.Name != "default:acacia_leaves"&&
			//!strings.Contains(node.Name, "flora") &&
			//!strings.Contains(node.Name, "sign") &&
			//!strings.Contains(node.Name, "flower") &&


			!strings.Contains(node.Name, "protector") {
			c := cr(node.Name, node.Param2)
			if c != nil {
				nodes = append(nodes, &NodeWithColor{
					Node:  node,
					Color: c,
				})

				if c.A == 255 || skip_alpha {
					break
				}
			}
		}

		cpos = cpos.Add(ipos)
	}

	return nodes, nil
}
// Later, name the builds inside Python_GUI/bin accordingly
// After that go to the Python_GUI.py and add your renderer builds to it ...
