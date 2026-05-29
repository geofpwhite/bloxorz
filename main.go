package main

import (
	"flag"
	"fmt"
	"net/http"

	"fortio.org/terminal/ansipixels"
)

func main() {
	outputGraphFlag := flag.Bool("graph", false, "Graph state space of current level")
	serveFlag := flag.Bool("serve", false, "Solve level and serve graph at localhost:8080")
	levelFlag := flag.Int("level", 1, "choose your level")
	flag.Parse()
	ap := ansipixels.NewAnsiPixels(60)
	level := LevelOne
	switch *levelFlag {
	case 1:
		level = LevelOne
	case 2:
		level = LevelTwo
	}
	s := level()
	if *serveFlag {
		path, coords, svg := solve(s)
		fmt.Println(path)
		fmt.Println(coords)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "image/svg+xml")
			w.Write(svg) //nolint:errcheck
		})
		fmt.Println("serving graph at http://localhost:8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println(err)
		}
		return
	}
	if *outputGraphFlag {
		path, coords, _ := solve(s)
		defer func() {
			fmt.Println(path)
			fmt.Println(coords)
		}()
	}
	curSteps := make([]direction, 0)
	_ = ap.Open()
	ap.HideCursor()
	ap.ClearScreen()
	defer func() {
		ap.ClearScreen()
		ap.ShowCursor()
		ap.Restore()
	}()
	err := ap.FPSTicks(func() bool {
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
		case CONTINUE:
		}
		DrawGame(ap, &s)
		return true
	})
	if err != nil {
		fmt.Println(err)
	}
}
