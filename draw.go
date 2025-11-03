package main

import (
	"fmt"
	"strconv"

	"fortio.org/terminal/ansipixels"
	"fortio.org/terminal/ansipixels/tcolor"
)

func DrawGame(ap *ansipixels.AnsiPixels, s *state) {
	// ap.ClearScreen()
	ap.WriteBg(tcolor.Gray.Color())
	for c, e := range s.floor {
		if !e {
			continue
		}
		ap.WriteAt(c[1], c[0], " ")
	}
	ap.WriteBg(tcolor.Red.Color())
	ap.WriteAt(s.endCoords[1], s.endCoords[0], " ")
	for coords, b := range s.buttons {
		if b.mustBeUpright {
			ap.WriteBg(tcolor.Blue.Color())
		} else {
			ap.WriteBg(tcolor.Purple.Color())
		}
		ap.WriteAt(coords[1], coords[0], " ")
	}

	ap.WriteBg(tcolor.Green.Color())
	for _, c := range s.block.coords {
		ap.WriteAt(c[1], c[0], " ")
	}

	ap.WriteBg(ap.Background.Color())
	str := ""
	for c, b := range s.buttons {
		str += fmt.Sprintf("%d,%d=%v", c[0], c[1], b.on)
	}
	ap.WriteAt(ap.W-25, ap.H-5, strconv.Itoa(len(s.floor))+str)

}
