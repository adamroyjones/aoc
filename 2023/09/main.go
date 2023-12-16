package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Printf("%v\n", parseFile("integration-part-1"))
}

func partTwo(filename string) int {
	d := parseFile(filename)
	out := 0
	for i := range d {
		out += extrapolateLeft(d[i])
	}
	return out
}

func partOne(filename string) int {
	d := parseFile(filename)
	out := 0
	for i := range d {
		out += extrapolateRight(d[i])
	}
	return out
}

func extrapolateLeft(vals []int) int {
	diffs := [][]int{vals}
	for i := 0; ; i++ {
		next := step(diffs[i])
		diffs = append(diffs, next)
		if constant(next) {
			break
		}
	}

	for i := len(diffs) - 2; i >= 0; i-- {
		diffs[i] = append([]int{diffs[i][0] - diffs[i+1][0]}, diffs[i]...)
	}
	return diffs[0][0]
}

func extrapolateRight(vals []int) int {
	diffs := [][]int{vals}
	for i := 0; ; i++ {
		next := step(diffs[i])
		diffs = append(diffs, next)
		if constant(next) {
			break
		}
	}

	for i := len(diffs) - 2; i >= 0; i-- {
		diffs[i] = append(diffs[i], last(diffs[i])+last(diffs[i+1]))
	}
	return last(diffs[0])
}

func constant(xs []int) bool {
	if len(xs) < 2 {
		return true
	}
	for i := 1; i < len(xs); i++ {
		if xs[i] != xs[i-1] {
			return false
		}
	}
	return true
}

func last(xs []int) int { return xs[len(xs)-1] }

func step(xs []int) []int {
	out := make([]int, 0, len(xs)-1)
	for i := 1; i < len(xs); i++ {
		out = append(out, xs[i]-xs[i-1])
	}
	assert(len(out) == len(xs)-1, "expected to decrement the list length")
	return out
}

func parseFile(filename string) [][]int {
	bs, err := os.ReadFile(filename)
	assert(err == nil, "reading a file")
	lines := strings.Split(strings.TrimSpace(string(bs)), "\n")
	data := make([][]int, 0, len(lines))
	for _, line := range lines {
		data = append(data, parseLine(line))
	}
	return data
}

func parseLine(ln string) []int {
	strs := strings.Fields(ln)
	xs := make([]int, 0, len(strs))
	for _, str := range strs {
		x, err := strconv.Atoi(str)
		assert(err == nil, "parsing an int")
		xs = append(xs, x)
	}
	return xs
}

func assert(b bool, msg string, args ...any) {
	if !b {
		panic("assertion failed: " + fmt.Sprintf(msg, args...))
	}
}
