package main

type state struct {
	block     block
	endCoords [2]int
	floor     map[[2]int]struct{}
	buttons   map[[2]int]button
}

type button struct {
	on            bool
	mustBeUpright bool
	tilesToToggle [][2]int
	state         *state
}

func (s *state) NewButton(coords [2]int, on, mustBeUpright bool, tiles [][2]int) {
	s.buttons[coords] = button{on, mustBeUpright, tiles, s}
}
func (b button) press() {
	b.on = !b.on
	toggle := func() {
		for _, coords := range b.tilesToToggle {
			b.state.floor[coords] = struct{}{}
		}
	}
	if !b.on {
		toggle = func() {
			for _, coords := range b.tilesToToggle {
				delete(b.state.floor, coords)
			}
		}
	}
	toggle()
}

type (
	orientation int
	direction   int
)

func (d direction) String() string {
	switch d {
	case UP:
		return "UP"
	case DOWN:
		return "DOWN"
	case LEFT:
		return "LEFT"
	case RIGHT:
		return "RIGHT"
	}
	return ""
}

const (
	UPRIGHT orientation = iota
	HORIZONTAL
	VERTICAL
)

const (
	UP direction = iota
	DOWN
	LEFT
	RIGHT
)

type block struct {
	orientation orientation
	coords      [][2]int
}

func (s *state) checkButtons() {
	for _, coords := range s.block.coords {
		if b, ok := s.buttons[coords]; ok {
			if (b.mustBeUpright && s.block.orientation == UPRIGHT) || (!b.mustBeUpright) {
				b.press()
			}
		}
	}
}

func (b block) Move(dir direction) block {
	newBlock := block{orientation: b.orientation}
	switch b.orientation {
	case UPRIGHT:
		switch dir {
		case UP:
			newBlock.orientation = VERTICAL
			newBlock.coords = [][2]int{{b.coords[0][0] - 1, b.coords[0][1]}, {b.coords[0][0] - 2, b.coords[0][1]}}
		case DOWN:
			newBlock.orientation = VERTICAL
			newBlock.coords = [][2]int{{b.coords[0][0] + 2, b.coords[0][1]}, {b.coords[0][0] + 1, b.coords[0][1]}}
		case LEFT:
			newBlock.orientation = HORIZONTAL
			newBlock.coords = [][2]int{{b.coords[0][0], b.coords[0][1] - 1}, {b.coords[0][0], b.coords[0][1] - 2}}
		case RIGHT:
			newBlock.orientation = HORIZONTAL
			newBlock.coords = [][2]int{{b.coords[0][0], b.coords[0][1] + 2}, {b.coords[0][0], b.coords[0][1] + 1}}
		}
	case HORIZONTAL:
		switch dir {
		case LEFT:
			newBlock.orientation = UPRIGHT
			newBlock.coords = [][2]int{{b.coords[1][0], b.coords[1][1] - 1}}
		case RIGHT:
			newBlock.orientation = UPRIGHT
			newBlock.coords = [][2]int{{b.coords[0][0], b.coords[0][1] + 1}}
		case UP:
			newBlock.coords = b.coords
			newBlock.coords[0][0]--
			newBlock.coords[1][0]--
		case DOWN:
			newBlock.coords = b.coords
			newBlock.coords[0][0]++
			newBlock.coords[1][0]++
		}
	case VERTICAL:
		switch dir {
		case UP:
			newBlock.orientation = UPRIGHT
			newBlock.coords = [][2]int{{b.coords[1][0] - 1, b.coords[1][1]}}
		case DOWN:
			newBlock.orientation = UPRIGHT
			newBlock.coords = [][2]int{{b.coords[0][0] + 1, b.coords[0][1]}}
		case LEFT:
			newBlock.coords = b.coords
			newBlock.coords[0][1]--
			newBlock.coords[1][1]--
		case RIGHT:
			newBlock.coords = b.coords
			newBlock.coords[0][1]++
			newBlock.coords[1][1]++
		}
	}
	return newBlock
}

type result int

const (
	WIN result = iota
	LOSE
	CONTINUE
)

func CheckState(s state) result {
	for _, c := range s.block.coords {
		_, ok := s.floor[c]
		if !ok {
			return LOSE
		}
	}
	if s.block.orientation == UPRIGHT && s.block.coords[0] == s.endCoords {
		return WIN
	}
	return CONTINUE
}
