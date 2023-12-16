package main

import (
	"testing"
)

func TestPartTwo(t *testing.T) {
	for _, tc := range []struct {
		filename string
		exp      int
	}{
		{filename: "integration-part-2", exp: 525152},
		{filename: "input", exp: 10861030975833},
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
		{filename: "integration-part-1", exp: 21},
		{filename: "input", exp: 7771},
	} {
		out := partOne(tc.filename)
		if tc.exp != out {
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
		out := rowToArrangements(tc.r)
		if tc.exp != out {
			t.Errorf("%v: expected %v; given %d", tc.r, tc.exp, out)
		}
	}
}

func TestUnfoldRowToArrangements(t *testing.T) {
	for _, tc := range []struct {
		r           row
		unfoldCount int
		exp         int
	}{
		{r: row{rawStr: "?.?????#???#?", counts: []int{1, 1, 2, 2}}, unfoldCount: 1, exp: 22},
		{r: row{rawStr: "?.?????#???#?", counts: []int{1, 1, 2, 2}}, unfoldCount: 2, exp: 700},
		{r: row{rawStr: "?.?????#???#?", counts: []int{1, 1, 2, 2}}, unfoldCount: 3, exp: 22516},
		{r: row{rawStr: "?.?????#???#?", counts: []int{1, 1, 2, 2}}, unfoldCount: 4, exp: 727792},
		{r: row{rawStr: "?.?????#???#?", counts: []int{1, 1, 2, 2}}, unfoldCount: 5, exp: 23570904},
	} {
		ur := unfoldRow(tc.r, tc.unfoldCount)
		out := rowToArrangements(ur)
		if tc.exp != out {
			t.Errorf("%v: expected %v; given %d", tc.r, tc.exp, out)
		}
	}
}
