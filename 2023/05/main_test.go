package main

import "testing"

func TestPartOne(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "integration-part-1", exp: 35},
		{filename: "input", exp: 265018614},
	} {
		solution := partOne(tc.filename)
		if tc.exp != solution {
			t.Fatalf("%s: not correct (given: %d, expected: %d)", tc.filename, solution, tc.exp)
		}
	}
}

func TestPartTwo(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "integration-part-2", exp: 46},
		{filename: "input", exp: 63179500},
	} {
		solution := partTwo(tc.filename)
		if tc.exp != solution {
			t.Fatalf("%s: not correct (given: %d, expected: %d)", tc.filename, solution, tc.exp)
		}
	}
}
