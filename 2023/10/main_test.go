package main

import "testing"

func TestPartOne(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "testdata/integration_part_1.1", exp: 4},
		{filename: "testdata/integration_part_1.2", exp: 8},
		{filename: "testdata/input", exp: 6754},
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
		{filename: "testdata/integration_part_2.1", exp: 4},
		{filename: "testdata/integration_part_2.2", exp: 4},
		{filename: "testdata/integration_part_2.3", exp: 8},
		{filename: "testdata/integration_part_2.4", exp: 10},
		{filename: "testdata/input", exp: 567},
	} {
		if out := partTwo(tc.filename); tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}
