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
		{filename: "testdata/integration_part_1", exp: 19114},
		{filename: "testdata/input", exp: 398527},
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
		{filename: "testdata/integration_part_2", exp: 167409079868000},
		{filename: "testdata/input", exp: 133973513090020},
	} {
		if out := partTwo(tc.filename); tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}

func TestProcessParts(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      []bool
	}{
		{filename: "testdata/integration_part_1", exp: []bool{true, false, true, false, true}},
	} {
		workflows, parts := parseFile(tc.filename)
		out := make([]bool, len(parts))
		for i := range parts {
			out[i] = processPart(workflows, parts[i]) > 0
		}
		if !slices.Equal(out, tc.exp) {
			t.Errorf("%s: expected %v, given %v", tc.filename, tc.exp, out)
		}
	}
}
