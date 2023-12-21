package main

import (
	"os"
	"slices"
	"strconv"
	"strings"
)

func partOne(filename string) int {
	d := parseFile(filename)
	out := 0
	for i := range d {
		out += extrapolateRight(d[i])
	}
	return out
}

func partTwo(filename string) int {
	d := parseFile(filename)
	out := 0
	for i := range d {
		out += extrapolateLeft(d[i])
	}
	return out
}

// parseFile returns rows of integers, one for each input line.
func parseFile(filename string) [][]int {
	bs, err := os.ReadFile(filename)
	assert(err == nil)

	lines := strings.Split(strings.TrimSpace(string(bs)), "\n")
	parsedLines := make([][]int, 0, len(lines))
	for _, line := range lines {
		strs := strings.Fields(line)
		parsedLine := make([]int, 0, len(strs))
		for _, str := range strs {
			x, err := strconv.Atoi(str)
			assert(err == nil)
			parsedLine = append(parsedLine, x)
		}
		parsedLines = append(parsedLines, parsedLine)
	}
	return parsedLines
}

func extrapolateRight(vals []int) int {
	diffs := [][]int{vals}
	for i := 0; ; i++ {
		next := step(diffs[i])
		diffs = append(diffs, next)
		if isConstant(next) {
			break
		}
	}

	for i := len(diffs) - 2; i >= 0; i-- {
		diffs[i] = append(diffs[i], last(diffs[i])+last(diffs[i+1]))
	}
	return last(diffs[0])
}

func extrapolateLeft(vals []int) int {
	diffs := [][]int{vals}
	for i := 0; ; i++ {
		next := step(diffs[i])
		diffs = append(diffs, next)
		if isConstant(next) {
			break
		}
	}

	for i := len(diffs) - 2; i >= 0; i-- {
		diffs[i] = append([]int{diffs[i][0] - diffs[i+1][0]}, diffs[i]...)
	}
	return diffs[0][0]
}

func isConstant(xs []int) bool {
	if len(xs) < 2 {
		return true
	}
	return !slices.ContainsFunc(xs[1:], func(x int) bool { return x != xs[0] })
}

func last(xs []int) int {
	assert(len(xs) > 0)
	return xs[len(xs)-1]
}

func step(xs []int) []int {
	out := make([]int, 0, len(xs)-1)
	for i := 1; i < len(xs); i++ {
		out = append(out, xs[i]-xs[i-1])
	}
	assert(len(out) == len(xs)-1)
	return out
}

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}
