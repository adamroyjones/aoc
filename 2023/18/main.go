package main

import (
	"fmt"
	"os"
)

type pair struct{ i, j int }

func main() {
	filename := "testdata/integration-part-1"
	isPartTwo := false
	bs, err := os.ReadFile(filename)
	assert(err == nil)
	insts := newInstructions(string(bs), isPartTwo)
	gryd := newGrid(insts)
	fmt.Println(gryd.trenchString())
	fmt.Println()
	fmt.Println(gryd.pipeString())
}

func partTwo(filename string) int {
	isPartTwo := true
	bs, err := os.ReadFile(filename)
	assert(err == nil)
	insts := newInstructions(string(bs), isPartTwo)
	sparseGryd := newSparseGrid(insts)
	quads := newQuadrangulation(sparseGryd)
	return sparseGryd.area(quads)
}

func partOne(filename string) int {
	isPartTwo := false
	bs, err := os.ReadFile(filename)
	assert(err == nil)
	insts := newInstructions(string(bs), isPartTwo)
	gryd := newGrid(insts)
	return gryd.area()
}

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}

func repeat[T any](x T, count int) []T {
	assert(count >= 0)
	xs := make([]T, count)
	for i := range xs {
		xs[i] = x
	}
	return xs
}

func ptr[T any](x T) *T { return &x }
