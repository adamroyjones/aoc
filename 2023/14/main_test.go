package main

import (
	"slices"
	"testing"
)

func TestPartOne(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "testdata/integration_part_1", exp: 136},
		{filename: "testdata/input", exp: 109755},
	} {
		if out := partOne(tc.filename); out != tc.exp {
			t.Errorf("%s: expected: %v, given: %v", tc.filename, tc.exp, out)
		}
	}
}

func TestPartTwo(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "testdata/integration_part_2", exp: 64},
		{filename: "testdata/input", exp: 90928},
	} {
		if out := partTwo(tc.filename); out != tc.exp {
			t.Errorf("%s: expected: %v, given: %v", tc.filename, tc.exp, out)
		}
	}
}

func TestSlideRowWest(t *testing.T) {
	for _, tc := range []struct {
		in  []tiletype
		exp []tiletype
	}{
		{in: []tiletype{}, exp: []tiletype{}},
		{in: []tiletype{TT_EMPTY}, exp: []tiletype{TT_EMPTY}},
		{in: []tiletype{TT_CUBIC_ROCK}, exp: []tiletype{TT_CUBIC_ROCK}},
		{in: []tiletype{TT_ROUND_ROCK}, exp: []tiletype{TT_ROUND_ROCK}},
		{in: []tiletype{TT_EMPTY, TT_ROUND_ROCK}, exp: []tiletype{TT_ROUND_ROCK, TT_EMPTY}},
		{in: []tiletype{TT_EMPTY, TT_EMPTY, TT_ROUND_ROCK}, exp: []tiletype{TT_ROUND_ROCK, TT_EMPTY, TT_EMPTY}},
		{in: []tiletype{TT_EMPTY, TT_CUBIC_ROCK, TT_ROUND_ROCK}, exp: []tiletype{TT_EMPTY, TT_CUBIC_ROCK, TT_ROUND_ROCK}},
		{in: []tiletype{TT_CUBIC_ROCK, TT_EMPTY, TT_ROUND_ROCK}, exp: []tiletype{TT_CUBIC_ROCK, TT_ROUND_ROCK, TT_EMPTY}},
	} {
		if out := slideRowEW(tc.in, DIRECTION_W); !slices.Equal(out, tc.exp) {
			t.Errorf("%v: expected: %v, given: %v", tc.in, tc.exp, out)
		}
	}
}
