package main

import (
	"fmt"
	"os"
	"slices"
)

func main() {
	filename := "testdata/integration-part-1"
	out := filenameToGrid(filename)
	fmt.Println(out.trenchString())
	fmt.Println()
	fmt.Println(out.pipeString())
}

func partOne(filename string) int {
	gryd := filenameToGrid(filename)
	return gryd.area()
}

func filenameToGrid(filename string) *grid {
	bs, err := os.ReadFile(filename)
	assert(err == nil)
	insts := newInstructions(string(bs))
	return newGrid(insts)
}

func equalLen[T any](xss [][]T) bool {
	if len(xss) == 0 {
		return true
	}
	return !slices.ContainsFunc(xss[1:], func(xs []T) bool { return len(xs) != len(xss[0]) })
}

func assertf(b bool, msg string, args ...any) {
	if !b {
		panic("assertion failed: " + fmt.Sprintf(msg, args...))
	}
}

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}

func ptr[T any](x T) *T { return &x }
