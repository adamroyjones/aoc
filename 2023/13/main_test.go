package main

import "testing"

func TestPartOne(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "testdata/integration_part_1", exp: 405},
		{filename: "testdata/input", exp: 41859},
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
		{filename: "testdata/integration_part_2", exp: 400},
		{filename: "testdata/input", exp: 30842},
	} {
		if out := partTwo(tc.filename); tc.exp != out {
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

	blockOne := newBlock(patternOne)
	d, ok := blockOne.horizontalIndex()
	if ok {
		t.Fatalf("first block: expected !ok")
	}
	if d != 0 {
		t.Fatalf("first block: expected 0, given %d", d)
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

	blockTwo := newBlock(patternTwo)
	d, ok = blockTwo.horizontalIndex()
	if !ok {
		t.Fatalf("second block: expected ok")
	}
	if d != 3 {
		t.Fatalf("second block: expected 0, given %d", d)
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

	blockOne := newBlock(patternOne)
	d, ok := blockOne.verticalIndex()
	if !ok {
		t.Fatalf("first block: expected ok")
	}
	if d != 4 {
		t.Fatalf("first block: expected 4, given %d", d)
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

	blockTwo := newBlock(patternTwo)
	d, ok = blockTwo.verticalIndex()
	if ok {
		t.Fatalf("second block: expected !ok")
	}
	if d != 0 {
		t.Fatalf("second block: expected 0, given %d", d)
	}
}
