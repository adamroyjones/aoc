package main

import (
	"slices"
	"strings"
	"testing"
)

func TestPartOne(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "testdata/integration_part_1", exp: 374},
		{filename: "testdata/input", exp: 9734203},
	} {
		if out := partOne(tc.filename); tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}

func TestPartTwo(t *testing.T) {
	for _, tc := range []struct {
		filename string
		factor   int
		exp      int
	}{
		{filename: "testdata/integration_part_1", factor: 2, exp: 374},
		{filename: "testdata/input", factor: 2, exp: 9734203},
		{filename: "testdata/integration_part_2", factor: 100, exp: 8410},
		{filename: "testdata/input", factor: 1_000_000, exp: 568914596391},
	} {
		if out := partTwo(tc.filename, tc.factor); tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}

func TestExpand(t *testing.T) {
	in := `
...#......
.......#..
#.........
..........
......#...
.#........
.........#
..........
.......#..
#...#.....
`

	exp := `
....#........
.........#...
#............
.............
.............
........#....
.#...........
............#
.............
.............
.........#...
#....#.......
`

	transform := func(s string) []string { return strings.Split(strings.TrimSpace(s), "\n") }
	if out := expand(transform(in)); !slices.Equal(transform(exp), out) {
		t.Errorf("expected %v, given %v", exp, out)
	}
}
