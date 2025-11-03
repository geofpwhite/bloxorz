package main

import (
	"flag"
	"fmt"

	"fortio.org/terminal/ansipixels"
)

func main() {
	outputGraphFlag := flag.Bool("graph", false, "Graph state space of current level")
	levelFlag := flag.Int("level", 1, "choose your level")
	flag.Parse()
	ap := ansipixels.NewAnsiPixels(60)
	level := LevelOne
	switch *levelFlag {
	case 2:
		level = LevelTwo
	}
	s := level()
	if *outputGraphFlag {
		path, coords := solve(s)
		defer func() {
			fmt.Println(path)
			fmt.Println(coords)
		}()
	}
	curSteps := make([]direction, 0)
	ap.Open()
	ap.HideCursor()
	ap.ClearScreen()
	defer func() {
		ap.ClearScreen()
		ap.ShowCursor()
		ap.Restore()
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
				s.checkButtons()
			case 'B':
				s.block = s.block.Move(DOWN)
				curSteps = append(curSteps, DOWN)
				s.checkButtons()
			case 'C':
				s.block = s.block.Move(RIGHT)
				curSteps = append(curSteps, RIGHT)
				s.checkButtons()
			case 'D':
				s.block = s.block.Move(LEFT)
				curSteps = append(curSteps, LEFT)
				s.checkButtons()
			}
		}
		ap.ClearScreen()
		result := CheckState(s)
		switch result {
		case LOSE:
			s = level()
		case WIN:
			return false
		}
		DrawGame(ap, &s)
		return true
	})
}
