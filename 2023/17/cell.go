package main

import "fmt"

type cell int

// String prints out the surface from above: the redder, the higher.
//
// This didn't prove to be useful, but it was useful to see how vacuous the
// terminal colour libraries are at root, in spite of their sizes. A perfectly
// general implementation comprises 2 lines (that are gofumpted to 7).
func (c cell) String() string {
	if c < 1 || c > 9 {
		panic("unexpected cell")
	}
	r := byte(255 * float64(c-1) / 8.0)
	return fg(bg(" ", r, 0, 0), r, 0, 0)
}

func fg(s string, r, g, b byte) string {
	return fmt.Sprintf("\u001B[38;2;%d;%d;%dm%s\u001B[39m", r, g, b, s)
}

func bg(s string, r, g, b byte) string {
	return fmt.Sprintf("\u001B[48;2;%d;%d;%dm%s\u001B[49m", r, g, b, s)
}
