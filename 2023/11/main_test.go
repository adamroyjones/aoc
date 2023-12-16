package main

import (
	"slices"
	"strings"
	"testing"
)

func TestPartTwo(t *testing.T) {
	for _, tc := range []struct {
		filename string
		factor   int
		exp      int
	}{
		{filename: "integration-part-1", factor: 2, exp: 374},
		{filename: "input", factor: 2, exp: 9734203},
		{filename: "integration-part-2", factor: 100, exp: 8410},
		{filename: "input", factor: 1_000_000, exp: 568914596391},
	} {
		out := partTwo(tc.filename, tc.factor)
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
		{filename: "integration-part-1", exp: 374},
		{filename: "input", exp: 9734203},
	} {
		out := partOne(tc.filename)
		if tc.exp != out {
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
	out := expand(transform(in))
	if !slices.Equal(transform(exp), out) {
		t.Errorf("expected %v; given %v", exp, out)
	}
}
