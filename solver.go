package main

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"log"
	"maps"
	"os"
	"slices"

	"fortio.org/terminal/ansipixels/tcolor"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

type solverNode struct {
	visitedNode
	prevMove      *direction
	prevNode      *solverNode
	prevBlock     *block
	curPath       int
	onButtonTiles map[[2]int]struct{}
}

type visitedNode struct {
	bx1, by1, bx2, by2 int
}

func solve(s state) ([]string, [][][2]int) {
	ctx := context.Background()
	g, err := graphviz.New(ctx)
	if err != nil {
		panic(err)
	}

	graph, err := g.Graph()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			panic(err)
		}
		g.Close()
	}()
	queue := make([]solverNode, 0)
	queue = append(queue, solverNode{visitedNode{1, 1, -1, -1}, nil, nil, nil, 0, make(map[[2]int]struct{})})
	gQueue := make([]*cgraph.Node, 0)
	bQueue := make([]block, 0)
	visited := make(map[visitedNode]int)
	block := block{coords: [][2]int{{4, 1}}, orientation: UPRIGHT}
	bQueue = append(bQueue, block)
	nodes := make(map[string]*cgraph.Node)
	edges := make(map[string]*cgraph.Edge)
	var done *solverNode
	gN, _ := graph.CreateNodeByName(queue[0].visitedNode.String())
	gQueue = append(gQueue, gN)
	nodes[queue[0].visitedNode.String()] = gN
outer:
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		b := bQueue[0]
		bQueue = bQueue[1:]
		n := gQueue[0]
		gQueue = gQueue[1:]
		visited[cur.visitedNode] = cur.curPath
		if button, ok := s.buttons[[2]int{cur.bx1, cur.by1}]; ok {
			on := false
			if (button.mustBeUpright && b.orientation == UPRIGHT) || !button.mustBeUpright {
				for _, coords := range button.tilesToToggle {
					if _, ok := cur.onButtonTiles[coords]; ok {
						delete(cur.onButtonTiles, coords)
					} else {
						cur.onButtonTiles[coords] = struct{}{}
						on = true
					}
				}
				fmt.Println("button pressed", on)
			}
		}
		if button, ok := s.buttons[[2]int{cur.bx2, cur.by2}]; ok && cur.bx2 != 0 && cur.by2 != 0 {
			on := false
			if !button.mustBeUpright {
				for _, coords := range button.tilesToToggle {
					if _, ok := cur.onButtonTiles[coords]; ok {
						delete(cur.onButtonTiles, coords)
					} else {
						cur.onButtonTiles[coords] = struct{}{}
						on = true
					}
				}
				fmt.Println("button pressed", on)
			}
		}
		for _, b := range s.buttons {
			if b.on {
				b.press()
			}
		}
		for c := range cur.onButtonTiles {
			s.floor[c] = struct{}{}
		}
		// g := graph.
		for i := range 4 {
			d := direction(i)
			if cur.prevMove != nil && opposite(*cur.prevMove, d) {
				continue
			}
			newBlock := b.Move(d)
			m := make(map[[2]int]struct{})
			maps.Copy(m, cur.onButtonTiles)
			newNode := solverNode{prevMove: &d, visitedNode: visitedNode{bx1: newBlock.coords[0][0], by1: newBlock.coords[0][1], bx2: -1, by2: -1}, prevNode: &cur, curPath: cur.curPath + 1, prevBlock: &b, onButtonTiles: m}
			if len(newBlock.coords) > 1 {
				newNode.bx2 = newBlock.coords[1][0]
				newNode.by2 = newBlock.coords[1][1]
			}
			// if num, ok := visited[newNode.visitedNode]; ok && num < newNode.curPath {
			// 	continue
			// }
			checkVal := check(s, newBlock)
			fmt.Println(checkVal, newBlock.coords, d.String())
			node, err := graph.CreateNodeByName(newNode.visitedNode.String())
			nodes[newNode.visitedNode.String()] = node
			if err != nil {
				panic("f")
			}
			key := cur.visitedNode.String() + "-" + newNode.String()

			e, err := graph.CreateEdgeByName(key, n, node)
			if err != nil {
				panic(err)
			}
			edges[key] = e
			switch checkVal {
			case LOSE:
				continue
			case WIN:
				done = &newNode
				break outer
				// continue
			}
			gQueue = append(gQueue, node)
			queue = append(queue, newNode)
			bQueue = append(bQueue, newBlock)

		}
	}
	if done == nil {
		return nil, nil
	}
	cur := done
	finish := nodes[cur.String()]
	finish.SetColor(tcolor.Green.String())
	path := []string{}
	coordPath := [][][2]int{}
	for cur.prevMove != nil {
		path = append(path, cur.prevMove.String())
		coordPath = append(coordPath, cur.prevBlock.coords)
		e := edges[cur.prevNode.String()+"-"+cur.String()]
		e.SetColor(tcolor.Green.String())
		cur = cur.prevNode
	}
	slices.Reverse(path)
	slices.Reverse(coordPath)
	var buf bytes.Buffer
	img, err := g.RenderImage(ctx, graph)
	if err != nil {
	}
	file, err := os.Create("output.png")
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}

	defer file.Close() // Ensure the file is closed

	// Encode the image to the file in PNG format
	err = png.Encode(file, img)
	if err != nil {
		log.Fatalf("failed to encode image: %v", err)
	}

	if err := g.Render(ctx, graph, "dot", &buf); err != nil {
		log.Fatal(err)
	}
	return path, coordPath
}

func check(s state, b block) result {
	for _, c := range b.coords {
		if _, ok := s.floor[c]; !ok {
			return LOSE
		}
	}
	if b.orientation == UPRIGHT && b.coords[0] == s.endCoords {
		return WIN
	}
	return CONTINUE
}
func opposite(d1, d2 direction) bool {
	return (min(d1, d2) == UP && max(d1, d2) == DOWN) || (max(d1, d2) == RIGHT && min(d1, d2) == LEFT)
}

func (vn *visitedNode) String() string {
	return fmt.Sprintf("%d,%d,%d,%d", vn.bx1, vn.by1, vn.bx2, vn.by2)
}
