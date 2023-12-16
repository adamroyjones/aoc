package main

import "testing"

func TestPartOne(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "integration-part-1.1", exp: 4},
		{filename: "integration-part-1.2", exp: 8},
		{filename: "input", exp: 6754},
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
		{filename: "integration-part-2.1", exp: 4},
		{filename: "integration-part-2.2", exp: 4},
		{filename: "integration-part-2.3", exp: 8},
		{filename: "integration-part-2.4", exp: 10},
		{filename: "input", exp: 567},
	} {
		out := partTwo(tc.filename)
		if tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}
