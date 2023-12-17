package main

import (
	"testing"
)

func TestPartTwo(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "integration-part-2", exp: 400},
		{filename: "input", exp: 30842},
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
		{filename: "integration-part-1", exp: 405},
		{filename: "input", exp: 41859},
	} {
		out := partOne(tc.filename)
		if tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}

func TestHorizontalIndex(t *testing.T) {
	patternOne := `
#.##..##.
..#.##.#.
##......#
##......#
..#.##.#.
..##..##.
#.#.##.#.
`

	blockOne := toBlock(patternOne)
	d, ok := blockOne.horizontalIndex()
	if ok {
		t.Errorf("first block: expected !ok")
	}
	if d != 0 {
		t.Errorf("first block: expected 0, given %d", d)
	}

	patternTwo := `
#...##..#
#....#..#
..##..###
#####.##.
#####.##.
..##..###
#....#..#
`

	blockTwo := toBlock(patternTwo)
	d, ok = blockTwo.horizontalIndex()
	if !ok {
		t.Errorf("second block: expected ok")
	}
	if d != 3 {
		t.Errorf("second block: expected 0, given %d", d)
	}
}

func TestVerticalIndex(t *testing.T) {
	patternOne := `
#.##..##.
..#.##.#.
##......#
##......#
..#.##.#.
..##..##.
#.#.##.#.
`

	blockOne := toBlock(patternOne)
	d, ok := blockOne.verticalIndex()
	if !ok {
		t.Errorf("first block: expected ok")
	}
	if d != 4 {
		t.Errorf("first block: expected 4, given %d", d)
	}

	patternTwo := `
#...##..#
#....#..#
..##..###
#####.##.
#####.##.
..##..###
#....#..#
`

	blockTwo := toBlock(patternTwo)
	d, ok = blockTwo.verticalIndex()
	if ok {
		t.Errorf("second block: expected !ok")
	}
	if d != 0 {
		t.Errorf("second block: expected 0, given %d", d)
	}
}
