package main

import (
	"fmt"

	"fortio.org/terminal/ansipixels"
)

func main() {
	ap := ansipixels.NewAnsiPixels(60)
	s := LevelTwo()
	path, coords := solve(s)
	curSteps := make([]direction, 0)
	fmt.Scanln()
	ap.Open()
	ap.HideCursor()
	ap.ClearScreen()
	defer func() {
		ap.ClearScreen()
		ap.ShowCursor()
		ap.Restore()
		fmt.Println(path)
		fmt.Println(coords)
		fmt.Println(curSteps)
	}()
	ap.FPSTicks(func() bool {
		if len(ap.Data) > 0 && ap.Data[0] == 'q' {
			return false
		}
		if len(ap.Data) > 2 {
			switch ap.Data[2] {
			case 'A':
				s.block = s.block.Move(UP)
				curSteps = append(curSteps, UP)
			case 'B':
				s.block = s.block.Move(DOWN)
				curSteps = append(curSteps, DOWN)
			case 'C':
				s.block = s.block.Move(RIGHT)
				curSteps = append(curSteps, RIGHT)
			case 'D':
				s.block = s.block.Move(LEFT)
				curSteps = append(curSteps, LEFT)
			}
		}
		ap.ClearScreen()
		result := CheckState(s)
		s.checkButtons()
		switch result {
		case LOSE:
			s = LevelTwo()
		case WIN:
			return false
		}
		DrawGame(ap, s)
		return true
	})
}
