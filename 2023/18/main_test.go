package main

import (
	"testing"
)

func TestPartTwo(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "testdata/integration-part-2", exp: 952408144115},
		{filename: "testdata/input", exp: 122109860712709},
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
		{filename: "testdata/integration-part-1", exp: 62},
		{filename: "testdata/input", exp: 62500},
	} {
		if out := partOne(tc.filename); tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}

func TestSparseGridArea(t *testing.T) {
	for _, tc := range []struct {
		name  string
		insts []instruction
		exp   int
	}{
		{
			name: "square",
			insts: []instruction{
				{dir: DIRECTION_R, dist: 1},
				{dir: DIRECTION_D, dist: 1},
				{dir: DIRECTION_L, dist: 1},
				{dir: DIRECTION_U, dist: 1},
			},
			exp: 4,
		},
		{
			name: "bigger_square",
			insts: []instruction{
				{dir: DIRECTION_R, dist: 2},
				{dir: DIRECTION_D, dist: 2},
				{dir: DIRECTION_L, dist: 2},
				{dir: DIRECTION_U, dist: 2},
			},
			exp: 9,
		},
		{
			name: "tetronimo/l",
			insts: []instruction{
				{dir: DIRECTION_R, dist: 1},
				{dir: DIRECTION_D, dist: 4},
				{dir: DIRECTION_R, dist: 2},
				{dir: DIRECTION_D, dist: 1},
				{dir: DIRECTION_L, dist: 3},
				{dir: DIRECTION_U, dist: 5},
			},
			exp: 16,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sparseGryd := newSparseGrid(tc.insts)
			quads := newQuadrangulation(sparseGryd)
			if out := sparseGryd.area(quads); out != tc.exp {
				t.Errorf("%s: expected %d, given %d", tc.name, tc.exp, out)
			}
		})
	}
}
