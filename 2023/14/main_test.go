package main

import (
	"slices"
	"testing"
)

func TestPartTwo(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "integration-part-2", exp: 64},
		{filename: "input", exp: 90928},
	} {
		out := partTwo(tc.filename)
		if out != tc.exp {
			t.Errorf("%s: expected: %v, given: %v", tc.filename, tc.exp, out)
		}
	}
}

func TestPartOne(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "integration-part-1", exp: 136},
		{filename: "input", exp: 109755},
	} {
		out := partOne(tc.filename)
		if out != tc.exp {
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
		out := slideRowEW(tc.in, DIRECTION_W)
		if !slices.Equal(out, tc.exp) {
			t.Errorf("in: %v, expected: %v, given: %v", tc.in, tc.exp, out)
		}
	}
}
