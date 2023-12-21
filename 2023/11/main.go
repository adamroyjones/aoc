package main

import (
	"os"
	"slices"
	"strings"
)

type pair struct{ i, j int }

func partOne(filepath string) int {
	bs, err := os.ReadFile(filepath)
	assert(err == nil)
	lines := strings.Split(strings.TrimSpace(string(bs)), "\n")
	expanded := expand(lines)
	pairs := locations(expanded)
	sum := 0
	for i, fst := range pairs {
		for _, snd := range pairs[i+1:] {
			sum += lInfinity(fst, snd)
		}
	}
	return sum
}

func partTwo(filepath string, factor int) int {
	bs, err := os.ReadFile(filepath)
	assert(err == nil)
	lines := strings.Split(strings.TrimSpace(string(bs)), "\n")
	pairs := locations(lines)
	expandableRows := locateExpandableRows(lines)
	expandableColumns := locateExpandableColumns(lines)

	sum := 0
	for i, fst := range pairs {
		for _, snd := range pairs[i+1:] {
			sum += lInfinityExpandable(fst, snd, expandableRows, expandableColumns, factor)
		}
	}
	return sum
}

func expand(lines []string) []string { return expandColumns(expandRows(lines)) }

func expandRows(lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		out = append(out, line)
		if !strings.ContainsFunc(line, func(r rune) bool { return r != '.' }) {
			out = append(out, line)
		}
	}
	return out
}

func expandColumns(lines []string) []string {
	out := make([]string, len(lines))
	assert(len(lines) > 0)
	assert(len(lines[0]) > 0)

	for col := range lines[0] {
		shouldExpand := true
		for _, line := range lines {
			if line[col] != '.' {
				shouldExpand = false
				break
			}
		}

		for i := range out {
			if shouldExpand {
				out[i] = out[i] + ("..")
			} else {
				out[i] = out[i] + string(lines[i][col])
			}
		}
	}

	return out
}

func locations(lines []string) []pair {
	pairs := []pair{}
	for i := range lines {
		for j, r := range lines[i] {
			if r == '#' {
				pairs = append(pairs, pair{i: i, j: j})
			}
		}
	}
	return pairs
}

func locateExpandableRows(lines []string) []int {
	out := []int{}
	for i, line := range lines {
		if !strings.ContainsFunc(line, func(r rune) bool { return r != '.' }) {
			out = append(out, i)
		}
	}
	return out
}

func locateExpandableColumns(lines []string) []int {
	out := []int{}
ColumnLoop:
	for j := range lines[0] {
		for _, line := range lines {
			if line[j] == '#' {
				continue ColumnLoop
			}
		}
		out = append(out, j)
	}
	return out
}

func lInfinityExpandable(fst, snd pair, expandableRows, expandableColumns []int, factor int) int {
	xDist := abs(fst.j - snd.j)
	for j := min(fst.j, snd.j) + 1; j < max(fst.j, snd.j); j++ {
		if _, ok := slices.BinarySearch(expandableColumns, j); ok {
			xDist += factor - 1
		}
	}

	yDist := abs(fst.i - snd.i)
	for i := min(fst.i, snd.i) + 1; i < max(fst.i, snd.i); i++ {
		if _, ok := slices.BinarySearch(expandableRows, i); ok {
			yDist += factor - 1
		}
	}

	return xDist + yDist
}

func lInfinity(fst, snd pair) int { return abs(fst.i-snd.i) + abs(fst.j-snd.j) }

func abs(x int) int { return max(x, -x) }

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}
