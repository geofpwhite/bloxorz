package main

func LevelOne() state {
	s := state{}
	s.floor = make(map[[2]int]bool)
	s.floor[[2]int{0, 0}] = true
	s.floor[[2]int{0, 1}] = true
	s.floor[[2]int{0, 2}] = true
	s.floor[[2]int{1, 0}] = true
	s.floor[[2]int{1, 1}] = true
	s.floor[[2]int{1, 2}] = true
	s.floor[[2]int{1, 3}] = true
	s.floor[[2]int{1, 4}] = true
	s.floor[[2]int{1, 5}] = true
	s.floor[[2]int{2, 0}] = true
	s.floor[[2]int{2, 1}] = true
	s.floor[[2]int{2, 2}] = true
	s.floor[[2]int{2, 3}] = true
	s.floor[[2]int{2, 4}] = true
	s.floor[[2]int{2, 5}] = true
	s.floor[[2]int{2, 6}] = true
	s.floor[[2]int{2, 7}] = true
	s.floor[[2]int{2, 8}] = true
	s.floor[[2]int{3, 1}] = true
	s.floor[[2]int{3, 2}] = true
	s.floor[[2]int{3, 3}] = true
	s.floor[[2]int{3, 4}] = true
	s.floor[[2]int{3, 5}] = true
	s.floor[[2]int{3, 6}] = true
	s.floor[[2]int{3, 7}] = true
	s.floor[[2]int{3, 8}] = true
	s.floor[[2]int{3, 9}] = true
	s.floor[[2]int{4, 5}] = true
	s.floor[[2]int{4, 6}] = true
	s.floor[[2]int{4, 7}] = true
	s.floor[[2]int{4, 8}] = true
	s.floor[[2]int{4, 9}] = true
	s.floor[[2]int{5, 6}] = true
	s.floor[[2]int{5, 7}] = true
	s.floor[[2]int{5, 8}] = true
	s.endCoords = [2]int{4, 7}
	s.block = block{orientation: UPRIGHT, coords: [][2]int{{1, 1}}}
	return s
}

func LevelTwo() state {
	s := state{}
	s.floor = make(map[[2]int]bool)
	s.buttons = make(map[[2]int]*button)
	s.block = block{orientation: UPRIGHT, coords: [][2]int{{4, 1}}}
	s.addFloor(1, 0, 5, 4)
	s.NewButton([2]int{2, 2}, false, false, [][2]int{{4, 5}})
	s.addFloor(0, 6, 5, 9)
	s.NewButton([2]int{1, 8}, false, true, [][2]int{{4, 10}, {4, 11}})
	s.addFloor(0, 12, 4, 14)
	s.endCoords = [2]int{1, 13}
	return s
}

func (s *state) addFloor(leftCornerX, leftCornerY, rightCornerX, rightCornerY int) {
	for i := leftCornerX; i <= rightCornerX; i++ {
		for j := leftCornerY; j <= rightCornerY; j++ {
			s.floor[[2]int{i, j}] = true
		}
	}
}
