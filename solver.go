package main

import (
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
	sNode              *solverNode
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
	start := solverNode{visitedNode{nil, 1, 1, -1, -1}, nil, nil, nil, 0, make(map[[2]int]struct{})}
	start.visitedNode.sNode = &start
	queue = append(queue, start)
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
						continue
					}
					cur.onButtonTiles[coords] = struct{}{}
					on = true
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
			s.floor[c] = true
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
			newNode.visitedNode.sNode = &newNode
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
			if err != nil {
				panic("error creating node")
			}
			nodes[newNode.visitedNode.String()] = node
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
	// var buf bytes.Buffer

	img, err := g.RenderImage(ctx, graph)
	if err != nil {
		panic("error rendering")
	}
	fmt.Println("creating")
	file, err := os.Create("output.png")
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}

	defer file.Close()

	fmt.Println("encoding")
	err = png.Encode(file, img)
	if err != nil {
		log.Fatalf("failed to encode image: %v", err)
	}

	// fmt.Println("rendering")
	// if err := g.Render(ctx, graph, "dot", &buf); err != nil {
	// 	log.Fatal(err)
	// }
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
	str := fmt.Sprintf("%d,%d,%d,%d+B[", vn.bx1, vn.by1, vn.bx2, vn.by2)
	ary := make([][2]int, len(vn.sNode.onButtonTiles))
	i := 0
	for b := range vn.sNode.onButtonTiles {
		ary[i] = b
		i++
	}
	slices.SortFunc(ary, func(a, b [2]int) int {
		switch {
		case a[0] == b[0] && a[1] == b[1]:
			return 0
		case a[0] < b[0]:
			return -1
		case a[0] > b[0]:
			return 1
		case a[1] < b[1]:
			return -1
		}
		return 1
	})
	for _, ar := range ary {
		str += fmt.Sprintf("(%d,%d)", ar[0], ar[1])
	}
	return str + "]"
}
