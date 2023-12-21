package main

import (
	"os"
	"slices"
)

func main() {
	// This illustrates the evolution. This requires an ANSI terminal.
	filename := "testdata/integration_part_1"
	bs, err := os.ReadFile(filename)
	assert(err == nil)
	mg := newMirrorGrid(string(bs))
	start := ray{i: 0, j: -1, dir: DIRECTION_R}
	print := true
	_ = newLightGrid(mg, start, print)
}

func partOne(filename string) int {
	bs, err := os.ReadFile(filename)
	assert(err == nil)
	mg := newMirrorGrid(string(bs))
	start := ray{i: 0, j: -1, dir: DIRECTION_R}
	print := false
	return newLightGrid(mg, start, print).count()
}

func partTwo(filename string) int {
	bs, err := os.ReadFile(filename)
	assert(err == nil)
	mg := newMirrorGrid(string(bs))
	count := 0
	print := false
	for _, ray := range mg.startingRays() {
		count = max(count, newLightGrid(mg, ray, print).count())
	}
	return count
}

func allEqualLen[T any](xss [][]T) bool {
	if len(xss) == 0 {
		return true
	}
	return !slices.ContainsFunc(xss[1:], func(xs []T) bool { return len(xs) != len(xss[0]) })
}

func allEqualLenStrs(ss []string) bool {
	if len(ss) == 0 {
		return true
	}
	return !slices.ContainsFunc(ss[1:], func(s string) bool { return len(s) != len(ss[0]) })
}

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}
