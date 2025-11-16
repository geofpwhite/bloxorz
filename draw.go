package main

import (
	"fmt"
	"strconv"
	"strings"

	"fortio.org/terminal/ansipixels"
	"fortio.org/terminal/ansipixels/tcolor"
)

func DrawToScale(ap *ansipixels.AnsiPixels, scale, x, y int) {
	for i := range scale {
		for j := range scale/2 + 1 {
			ap.WriteAt(x*scale+i, y*scale/2+j, " ")
		}
	}
}

func DrawToScaleColor(ap *ansipixels.AnsiPixels, scale, x, y int, color tcolor.Color) {
	ap.WriteBg(color)
	for i := range scale {
		for j := range scale {
			ap.WriteAt(x*scale+i, y*scale+j, " ")
		}
	}
}

func DrawGame(ap *ansipixels.AnsiPixels, s *state) {
	// ap.ClearScreen()
	size := GetScale(s.floor, ap.W, ap.H)
	ap.WriteBg(tcolor.Gray.Color())
	for c, e := range s.floor {
		if !e {
			continue
		}
		DrawToScale(ap, size, c[1], c[0])
	}
	ap.WriteBg(tcolor.Red.Color())
	DrawToScale(ap, size, s.endCoords[1], s.endCoords[0])
	for coords, b := range s.buttons {
		if b.mustBeUpright {
			ap.WriteBg(tcolor.Blue.Color())
		} else {
			ap.WriteBg(tcolor.Purple.Color())
		}
		DrawToScale(ap, size, coords[1], coords[0])
	}

	ap.WriteBg(tcolor.Green.Color())
	for _, c := range s.block.coords {
		DrawToScale(ap, size, c[1], c[0])
	}

	ap.WriteBg(ap.Background.Color())
	builder := strings.Builder{}
	builder.WriteString(strconv.Itoa(len(s.floor)))
	for c, b := range s.buttons {
		builder.WriteString(fmt.Sprintf("%d,%d=%v", c[0], c[1], b.on))
	}
	ap.WriteAt(ap.W-25, ap.H-5, "%s", builder.String())
}

func GetScale(floor map[[2]int]bool, w, h int) int {
	largestXcoords, largestYcoords := [2]int{}, [2]int{}
	for c, e := range floor {
		if !e {
			continue
		}
		if c[0] > largestXcoords[0] {
			largestXcoords = c
		}
		if c[1] > largestYcoords[1] {
			largestYcoords = c
		}
	}
	return min(h/(largestXcoords[0]), w/largestYcoords[1])
}
