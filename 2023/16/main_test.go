package main

import (
	"strings"
	"testing"
)

func TestPartOne(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "testdata/integration_part_1", exp: 46},
		{filename: "testdata/input", exp: 7392},
	} {
		if out := partOne(tc.filename); tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}

func TestPartTwo(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "testdata/integration_part_2", exp: 51},
		{filename: "testdata/input", exp: 7665},
	} {
		out := partTwo(tc.filename)
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

	mg := newMirrorGrid(input)
	start := ray{i: 0, j: -1, dir: DIRECTION_R}
	lg := newLightGrid(mg, start, false)
	if strings.TrimSpace(lg.String()) != strings.TrimSpace(exp) {
		t.Errorf("expected:\n%s\n\ngiven:\n%s\n", strings.TrimSpace(exp), strings.TrimSpace(lg.String()))
	}
}
