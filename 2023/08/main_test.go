package main

import "testing"

func TestPartOne(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "integration-part-1.1", exp: 2},
		{filename: "integration-part-1.2", exp: 6},
		{filename: "input", exp: 12361},
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
		{filename: "integration-part-2", exp: 6},
		{filename: "input", exp: 18215611419223},
	} {
		out := partTwo(tc.filename)
		if tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}

func TestGCD(t *testing.T) {
	for _, tc := range []struct {
		xs  []int
		exp int
	}{
		{xs: []int{1}, exp: 1},
		{xs: []int{1, 2}, exp: 1},
		{xs: []int{2, 2}, exp: 2},
		{xs: []int{2, 4}, exp: 2},
		{xs: []int{2, 3}, exp: 1},
		{xs: []int{4, 6}, exp: 2},
		{xs: []int{4, 6, 5}, exp: 1},
		{xs: []int{4, 6, 10}, exp: 2},
		{xs: []int{12, 6, 30}, exp: 6},
	} {
		out := gcd(tc.xs)
		if tc.exp != out {
			t.Errorf("%v: expected %d, given %d", tc.xs, tc.exp, out)
		}
	}
}

func TestLCM(t *testing.T) {
	for _, tc := range []struct {
		xs  []int
		exp int
	}{
		{xs: []int{1}, exp: 1},
		{xs: []int{1, 2}, exp: 2},
		{xs: []int{2, 2}, exp: 2},
		{xs: []int{2, 4}, exp: 4},
		{xs: []int{2, 3}, exp: 6},
		{xs: []int{4, 6}, exp: 12},
		{xs: []int{4, 6, 5}, exp: 60},
	} {
		out := lcm(tc.xs)
		if tc.exp != out {
			t.Errorf("%v: expected %d, given %d", tc.xs, tc.exp, out)
		}
	}
}
