package main

import (
	"slices"
	"strconv"
	"strings"
)

type grid [][]cell

// newGrid takes in a string representation of a grid and produces a... grid.
func newGrid(s string) *grid {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	assert(equalLenStr(lines), "the input must be a rectangle")

	gryd := make(grid, len(lines))
	for i := range gryd {
		gryd[i] = make([]cell, len(lines[0]))
		for j := range gryd[i] {
			v, err := strconv.Atoi(string(lines[i][j]))
			assert(err == nil, "unable to parse an int")
			gryd[i][j] = cell(v)
		}
	}
	return &gryd
}

func (gryd *grid) graph() *graph {
	assert(gryd != nil, "the grid must be non-nil")
	assert(len(*gryd) > 0 && len((*gryd)[0]) > 0, "the grid must be be non-empty")

	graf := newGraph(gryd)
	currentNodes := []*node{graf.root}
	for len(currentNodes) > 0 {
		nextNodes := []*node{}
		for _, currentNode := range currentNodes {
			nextNextNodes := graf.nextNodes(currentNode)
			nextNextNodes = slices.DeleteFunc(nextNextNodes, func(n *node) bool {
				l := n.label()
				graf.addEdge(currentNode, n)
				if _, ok := graf.labelToNode[l]; !ok {
					graf.labelToNode[l] = n
					return false
				}
				return true
			})
			nextNodes = append(nextNodes, nextNextNodes...)
		}
		currentNodes = nextNodes
	}

	return graf
}

// dims assumes that we were given a rectangle. This is given by newGrid but
// it's not something that can be enforced in the type system.
func (gryd grid) dims() (int, int) {
	imax := len(gryd)
	if imax == 0 {
		return 0, 0
	}
	return imax, len(gryd[0])
}

func (gryd grid) String() string {
	assert(gryd != nil, "the grid must be non-nil")
	var sb strings.Builder
	for i := range gryd {
		for j := range gryd[i] {
			sb.WriteString(gryd[i][j].String())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
