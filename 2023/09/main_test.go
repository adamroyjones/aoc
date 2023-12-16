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
		{filename: "integration-part-2", exp: 2},
		{filename: "input", exp: 803},
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
		{filename: "integration-part-1", exp: 114},
		{filename: "input", exp: 1681758908},
	} {
		out := partOne(tc.filename)
		if tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}

func TestExtrapolateRight(t *testing.T) {
	for _, tc := range []struct {
		input []int
		exp   int
	}{
		{input: []int{0, 3, 6, 9, 12, 15}, exp: 18},
		{input: []int{1, 3, 6, 10, 15, 21}, exp: 28},
		{input: []int{10, 13, 16, 21, 30, 45}, exp: 68},
	} {
		out := extrapolateRight(tc.input)
		if tc.exp != out {
			t.Errorf("%v: expected %d, given %d", tc.input, tc.exp, out)
		}
	}
}

func TestStep(t *testing.T) {
	for _, tc := range []struct {
		input []int
		exp   []int
	}{
		{input: []int{0, 3, 6, 9, 12, 15}, exp: []int{3, 3, 3, 3, 3}},
		{input: []int{1, 3, 6, 10, 15, 21}, exp: []int{2, 3, 4, 5, 6}},
		{input: []int{10, 13, 16, 21, 30, 45}, exp: []int{3, 3, 5, 9, 15}},
	} {
		out := step(tc.input)
		if !slices.Equal(tc.exp, out) {
			t.Errorf("%v: expected %v, given %v", tc.input, tc.exp, out)
		}
	}
}
