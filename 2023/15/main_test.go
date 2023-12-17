package main

import "testing"

func TestPartTwo(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "integration-part-2", exp: 145},
		{filename: "input", exp: 267372},
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
		{filename: "integration-part-1", exp: 1320},
		{filename: "input", exp: 517965},
	} {
		out := partOne(tc.filename)
		if tc.exp != out {
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
		out := hash(tc.in)
		if tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.in, tc.exp, out)
		}
	}
}
