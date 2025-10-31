package main

func LevelOne() state {
	s := state{}
	s.floor = make(map[[2]int]struct{})
	s.floor[[2]int{0, 0}] = struct{}{}
	s.floor[[2]int{0, 1}] = struct{}{}
	s.floor[[2]int{0, 2}] = struct{}{}
	s.floor[[2]int{1, 0}] = struct{}{}
	s.floor[[2]int{1, 1}] = struct{}{}
	s.floor[[2]int{1, 2}] = struct{}{}
	s.floor[[2]int{1, 3}] = struct{}{}
	s.floor[[2]int{1, 4}] = struct{}{}
	s.floor[[2]int{1, 5}] = struct{}{}
	s.floor[[2]int{2, 0}] = struct{}{}
	s.floor[[2]int{2, 1}] = struct{}{}
	s.floor[[2]int{2, 2}] = struct{}{}
	s.floor[[2]int{2, 3}] = struct{}{}
	s.floor[[2]int{2, 4}] = struct{}{}
	s.floor[[2]int{2, 5}] = struct{}{}
	s.floor[[2]int{2, 6}] = struct{}{}
	s.floor[[2]int{2, 7}] = struct{}{}
	s.floor[[2]int{2, 8}] = struct{}{}
	s.floor[[2]int{3, 1}] = struct{}{}
	s.floor[[2]int{3, 2}] = struct{}{}
	s.floor[[2]int{3, 3}] = struct{}{}
	s.floor[[2]int{3, 4}] = struct{}{}
	s.floor[[2]int{3, 5}] = struct{}{}
	s.floor[[2]int{3, 6}] = struct{}{}
	s.floor[[2]int{3, 7}] = struct{}{}
	s.floor[[2]int{3, 8}] = struct{}{}
	s.floor[[2]int{3, 9}] = struct{}{}
	s.floor[[2]int{4, 5}] = struct{}{}
	s.floor[[2]int{4, 6}] = struct{}{}
	s.floor[[2]int{4, 7}] = struct{}{}
	s.floor[[2]int{4, 8}] = struct{}{}
	s.floor[[2]int{4, 9}] = struct{}{}
	s.floor[[2]int{5, 6}] = struct{}{}
	s.floor[[2]int{5, 7}] = struct{}{}
	s.floor[[2]int{5, 8}] = struct{}{}
	s.endCoords = [2]int{4, 7}
	s.block = block{orientation: UPRIGHT, coords: [][2]int{{1, 1}}}
	return s
}
func LevelTwo() state {
	s := state{}
	s.floor = make(map[[2]int]struct{})
	s.buttons = make(map[[2]int]button)
	s.block = block{orientation: UPRIGHT, coords: [][2]int{{4, 1}}}
	s.addFloor(1, 0, 5, 4)
	s.NewButton([2]int{2, 2}, false, false, [][2]int{{4, 4}, {4, 5}})
	s.addFloor(0, 6, 5, 9)
	s.NewButton([2]int{1, 8}, false, true, [][2]int{{4, 10}, {4, 11}})
	s.addFloor(0, 12, 4, 14)
	s.endCoords = [2]int{1, 13}
	return s
}

func (s *state) addFloor(leftCornerX, leftCornerY, rightCornerX, rightCornerY int) {
	for i := leftCornerX; i <= rightCornerX; i++ {
		for j := leftCornerY; j <= rightCornerY; j++ {
			s.floor[[2]int{i, j}] = struct{}{}
		}
	}
}
