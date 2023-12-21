package main

import "testing"

func TestPartOne(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "testdata/integration_part_1", exp: 1320},
		{filename: "testdata/input", exp: 517965},
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
		{filename: "testdata/integration_part_2", exp: 145},
		{filename: "testdata/input", exp: 267372},
	} {
		if out := partTwo(tc.filename); tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}

func TestHash(t *testing.T) {
	for _, tc := range []struct {
		in  string
		exp int
	}{
		{in: "HASH", exp: 52},
	} {
		if out := hash(tc.in); tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.in, tc.exp, out)
		}
	}
}
