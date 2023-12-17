package main

import (
	"strings"
	"testing"
)

func TestPartTwo(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "integration-part-2", exp: 51},
		{filename: "input", exp: 7665},
	} {
		out := partTwo(tc.filename)
		if tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}

func TestPartOne(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "integration-part-1", exp: 46},
		{filename: "input", exp: 7392},
	} {
		out := partOne(tc.filename)
		if tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}

func TestLightGrid(t *testing.T) {
	input := `
.|...\....
|.-.\.....
.....|-...
........|.
..........
.........\
..../.\\..
.-.-/..|..
.|....-|.\
..//.|....
`

	exp := `
######....
.#...#....
.#...#####
.#...##...
.#...##...
.#...##...
.#..####..
########..
.#######..
.#...#.#..
`

	mg := toMirrorGrid(input)
	start := ray{i: 0, j: -1, dir: DIRECTION_R}
	lg := toLightGrid(mg, start, false)
	if strings.TrimSpace(lg.String()) != strings.TrimSpace(exp) {
		t.Errorf("expected:\n%s\n\ngiven:\n%s\n", strings.TrimSpace(exp), strings.TrimSpace(lg.String()))
	}
}
