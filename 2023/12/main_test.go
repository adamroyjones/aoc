package main

import (
	"testing"
)

func TestPartOne(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "testdata/integration_part_1", exp: 21},
		{filename: "testdata/input", exp: 7771},
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
		{filename: "testdata/integration_part_2", exp: 525152},
		{filename: "testdata/input", exp: 10861030975833},
	} {
		if out := partTwo(tc.filename); tc.exp != out {
			t.Errorf("%s: expected %d, given %d", tc.filename, tc.exp, out)
		}
	}
}

func TestRowToArrangements(t *testing.T) {
	for _, tc := range []struct {
		r   row
		exp int
	}{
		{r: row{rawStr: "", counts: []int{}}, exp: 1},
		{r: row{rawStr: "", counts: []int{1}}, exp: 0},
		{r: row{rawStr: ".", counts: []int{1}}, exp: 0},
		{r: row{rawStr: "#", counts: []int{1}}, exp: 1},
		{r: row{rawStr: "#.", counts: []int{1}}, exp: 1},
		{r: row{rawStr: "#?", counts: []int{1}}, exp: 1},
		{r: row{rawStr: ".#", counts: []int{1}}, exp: 1},
		{r: row{rawStr: ".#.", counts: []int{1}}, exp: 1},
		{r: row{rawStr: "?", counts: []int{1}}, exp: 1},
		{r: row{rawStr: "??", counts: []int{1}}, exp: 2},
		{r: row{rawStr: "??", counts: []int{2}}, exp: 1},
		{r: row{rawStr: "??.", counts: []int{2}}, exp: 1},
		{r: row{rawStr: "???", counts: []int{1, 1}}, exp: 1},
		{r: row{rawStr: "?.#", counts: []int{1, 1}}, exp: 1},
		{r: row{rawStr: "?.#", counts: []int{1}}, exp: 1},
		{r: row{rawStr: ".??.", counts: []int{1}}, exp: 2},
		{r: row{rawStr: "??.??", counts: []int{1, 1}}, exp: 4},
		{r: row{rawStr: "??.??.??", counts: []int{1, 1, 1}}, exp: 8},
		{r: row{rawStr: "???.###", counts: []int{1, 1, 3}}, exp: 1},
		{r: row{rawStr: "?.#?", counts: []int{1, 1}}, exp: 1},
		{r: row{rawStr: "??.#?", counts: []int{1, 1}}, exp: 2},
		{r: row{rawStr: ".??..??...?##.", counts: []int{1, 1, 3}}, exp: 4},
		{r: row{rawStr: "?#", counts: []int{1}}, exp: 1},
		{r: row{rawStr: "?#?", counts: []int{1}}, exp: 1},
		{r: row{rawStr: "?#?#?#", counts: []int{1, 3}}, exp: 1},
		{r: row{rawStr: "?#?#?#?#?#?#?#?", counts: []int{1, 3, 1, 6}}, exp: 1},
		{r: row{rawStr: "????.#...#...", counts: []int{4, 1, 1}}, exp: 1},
		{r: row{rawStr: "????.######..#####.", counts: []int{1, 6, 5}}, exp: 4},
		{r: row{rawStr: "?###????????", counts: []int{3, 2, 1}}, exp: 10},
	} {
		if out := rowToArrangements(tc.r); tc.exp != out {
			t.Errorf("%v: expected %v; given %d", tc.r, tc.exp, out)
		}
	}
}
