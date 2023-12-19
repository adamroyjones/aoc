package main

import (
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestPartTwo(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "testdata/integration-part-2.1", exp: 94},
		{filename: "testdata/integration-part-2.2", exp: 71},
		{filename: "testdata/input", exp: 1082},
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
		{filename: "testdata/integration-part-1", exp: 102},
		{filename: "testdata/input", exp: 942},
	} {
		out := partOne(tc.filename)
		if tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}

func TestShortestPathPartTwo(t *testing.T) {
	bs, err := os.ReadFile("testdata/integration-part-2-manual")
	if err != nil {
		t.Fatalf("expected to be able to read a file")
	}

	tcs := strings.Split(strings.TrimSpace(string(bs)), "=")
	for _, tc := range tcs {
		input, output, ok := strings.Cut(tc, "-")
		if !ok {
			t.Fatalf("no '-' in the test case\ncase:\n%s", tc)
		}
		input, output = strings.TrimSpace(input), strings.TrimSpace(output)
		g := newGrid(input)
		graf := g.graph(partTwoMin, partTwoMax)

		exp, err := strconv.Atoi(output)
		if err != nil {
			t.Fatalf("expected to be able to parse an int (given: %s)", output)
		}
		if out := graf.shortestPath(); exp != out {
			t.Fatalf("expected %d, given %d\ncase:\n%s", exp, out, input)
		}
	}
}

func TestShortestPathPartOne(t *testing.T) {
	bs, err := os.ReadFile("testdata/integration-part-1-manual")
	if err != nil {
		t.Fatalf("expected to be able to read a file")
	}

	tcs := strings.Split(strings.TrimSpace(string(bs)), "=")
	for _, tc := range tcs {
		input, output, ok := strings.Cut(tc, "-")
		if !ok {
			t.Fatalf("no '-' in the test case\ncase:\n%s", tc)
		}
		input, output = strings.TrimSpace(input), strings.TrimSpace(output)
		g := newGrid(input)
		graf := g.graph(partOneMin, partOneMax)

		exp, err := strconv.Atoi(output)
		if err != nil {
			t.Fatalf("expected to be able to parse an int (given: %s)", output)
		}
		if out := graf.shortestPath(); exp != out {
			t.Fatalf("expected %d, given %d\ncase:\n%s", exp, out, input)
		}
	}
}
