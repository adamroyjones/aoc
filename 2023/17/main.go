package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

// This uses graphviz to produce a visualisation.
func main() {
	filename := "integration-part-1-dot"
	bs, err := os.ReadFile(filename)
	assert(err == nil, "unable to read a file")
	g := newGrid(string(bs))
	graf := g.graph()
	renderDot(graf.String())
	fmt.Println(graf.shortestPath())
}

func partOne(filename string) int {
	bs, err := os.ReadFile(filename)
	assert(err == nil, "unable to read a file")
	g := newGrid(string(bs))
	return g.graph().shortestPath()
}

func renderDot(dot string) {
	f, err := os.Create("/tmp/dot.input")
	assert(err == nil, "unable to create file")
	if _, err := f.Write([]byte(dot)); err != nil {
		ferr := f.Close()
		assert(errors.Join(err, ferr) == nil, "writing to the file")
	}
	err = f.Close()
	assert(err == nil, "closing the file")

	cmd := exec.Command("dot", "-T", "svg", "-o", "/tmp/dot.svg", "/tmp/dot.input")
	bs, err := cmd.CombinedOutput()
	assert(err == nil, "running dot: %s", strings.TrimSpace(string(bs)))

	cmd = exec.Command("firefox", "/tmp/dot.svg")
	bs, err = cmd.CombinedOutput()
	assert(err == nil, "opening the dot graph: %s", strings.TrimSpace(string(bs)))
}

func equalLenStr(ss []string) bool {
	if len(ss) == 0 {
		return true
	}
	return !slices.ContainsFunc(ss[1:], func(s string) bool { return len(s) != len(ss[0]) })
}

func compact(ns ...*node) []*node {
	out := make([]*node, 0, len(ns))
	for _, n := range ns {
		if n != nil {
			out = append(out, n)
		}
	}
	return out
}

// safeappend was introduced as I had a horrific bug with append mutating a
// backing slice that was reused. The bug basically amounts to
//
//	xs := make([]int, 2, 16)
//	xs[0], xs[1] = 0, 1
//	ys := append(xs[1:], 2)
//	xs = append(xs, 3)
//
// where y is now [1, 3]. Yes, this is well-documented. Yes, I would prefer
// value semantics and immutability.
func safeappend[T any](xs []T, ys ...T) []T {
	zs := make([]T, 0, len(xs)+len(ys))
	for _, x := range xs {
		zs = append(zs, x)
	}
	for _, y := range ys {
		zs = append(zs, y)
	}
	return zs
}

func assert(b bool, msg string, args ...any) {
	if !b {
		panic("assertion failed: " + fmt.Sprintf(msg, args...))
	}
}
