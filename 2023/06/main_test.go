package main

import "testing"

func TestPartOne(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "integration-part-1", exp: 288},
		{filename: "input", exp: 4568778},
	} {
		out := partOne(tc.filename)
		if tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}

func TestPartTwo(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "integration-part-2", exp: 71503},
		{filename: "input", exp: 28973936},
	} {
		out := partTwo(tc.filename)
		if tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}
