package main

import "testing"

func TestPartOne(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "testdata/integration_part_1.1", exp: 2},
		{filename: "testdata/integration_part_1.2", exp: 6},
		{filename: "testdata/input", exp: 12361},
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
		{filename: "testdata/integration_part_2", exp: 6},
		{filename: "testdata/input", exp: 18215611419223},
	} {
		if out := partTwo(tc.filename); tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}

func TestGCD(t *testing.T) {
	for _, tc := range []struct {
		x, y int
		exp  int
	}{
		{x: 1, y: 2, exp: 1},
		{x: 2, y: 2, exp: 2},
		{x: 2, y: 4, exp: 2},
		{x: 2, y: 3, exp: 1},
		{x: 4, y: 6, exp: 2},
		{x: 4, y: 5, exp: 1},
		{x: 12, y: 30, exp: 6},
	} {
		if out := gcd(tc.x, tc.y); tc.exp != out {
			t.Errorf("%d, %d: expected %d, given %d", tc.x, tc.y, tc.exp, out)
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
		if out := lcm(tc.xs); tc.exp != out {
			t.Errorf("%v: expected %d, given %d", tc.xs, tc.exp, out)
		}
	}
}
