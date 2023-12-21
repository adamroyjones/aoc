package main

import (
	"os"
	"slices"
	"strings"
)

func partOne(filepath string) int {
	blocks := parseFile(filepath)

	horizontalIndices := make([]int, 0, len(blocks))
	verticalIndices := make([]int, 0, len(blocks))
	for _, block := range blocks {
		verticalIndex, ok := block.verticalIndex()
		if ok {
			verticalIndices = append(verticalIndices, verticalIndex)
		}

		horizontalIndex, ok := block.horizontalIndex()
		if ok {
			horizontalIndices = append(horizontalIndices, horizontalIndex)
		}
	}

	return score(horizontalIndices, verticalIndices)
}

func partTwo(filepath string) int {
	blocks := parseFile(filepath)

	horizontalIndices := make([]int, 0, len(blocks))
	verticalIndices := make([]int, 0, len(blocks))
	for _, block := range blocks {
		verticalIndex, ok := block.verticalIndexNearMiss()
		if ok {
			verticalIndices = append(verticalIndices, verticalIndex)
		}

		horizontalIndex, ok := block.horizontalIndexNearMiss()
		if ok {
			horizontalIndices = append(horizontalIndices, horizontalIndex)
		}
	}

	return score(horizontalIndices, verticalIndices)
}

func parseFile(filepath string) []block {
	bs, err := os.ReadFile(filepath)
	assert(err == nil)

	blockStrs := strings.Split(strings.TrimSpace(string(bs)), "\n\n")
	blocks := make([]block, len(blockStrs))
	for i, blockStr := range blockStrs {
		blocks[i] = newBlock(blockStr)
	}
	return blocks
}

func score(horizontalIndices, verticalIndices []int) int {
	s := 0
	for _, verticalIndex := range verticalIndices {
		s += verticalIndex + 1
	}
	for _, horizontalIndex := range horizontalIndices {
		s += 100 * (horizontalIndex + 1)
	}
	return s
}

func equalLen[T any](xss [][]T) bool {
	if len(xss) == 0 {
		return true
	}
	return !slices.ContainsFunc(xss[1:], func(xs []T) bool { return len(xss[0]) != len(xs) })
}

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}
