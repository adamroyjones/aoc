package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

type sitetype int

const (
	ST_ASH   sitetype = 0
	ST_ROCKS sitetype = 1
)

func (st sitetype) String() string {
	switch st {
	case ST_ASH:
		return "."
	case ST_ROCKS:
		return "#"
	default:
		panic("unexpected site type")
	}
}

type block [][]sitetype

func (b block) String() string {
	var sb strings.Builder

	for _, row := range b {
		for _, cell := range row {
			sb.WriteString(cell.String())
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (b block) col(j int) []sitetype {
	out := make([]sitetype, len(b))
	for i := range b {
		out[i] = b[i][j]
	}
	return out
}

// We'll casually assume that there is exactly one index.
func (b block) horizontalIndex() (int, bool) {
Loop:
	// The -1 matters as the last row is, by definition, inadmissible as a line of
	// reflection.
	for i := 0; i < len(b)-1; i++ {
		for j := i; j >= 0; j-- {
			k := i + (i - j) + 1 // i+1, i+2, ...
			if k == len(b) {
				break
			}
			if !slices.Equal(b[j], b[k]) {
				continue Loop
			}
		}
		return i, true
	}
	return 0, false
}

// We'll casually assume that there is exactly one index.
func (b block) horizontalIndexNearMiss() (int, bool) {
Loop:
	// The -1 matters as the last row is, by definition, inadmissible as a line of
	// reflection.
	for i := 0; i < len(b)-1; i++ {
		nearMisses := 0
		for j := i; j >= 0; j-- {
			k := i + (i - j) + 1 // i+1, i+2, ...
			if k == len(b) {
				break
			}
			if slices.Equal(b[j], b[k]) {
				// No-op.
			} else if nearMiss(b[j], b[k]) {
				nearMisses++
			} else {
				continue Loop
			}
		}
		if nearMisses == 1 {
			return i, true
		}
	}
	return 0, false
}

// We'll casually assume that there is exactly one index.
func (b block) verticalIndex() (int, bool) {
	if len(b) == 0 {
		return 0, false
	}

Loop:
	// The -1 matters as the last column is, by definition, inadmissible as a line
	// of reflection.
	for i := 0; i < len(b[0])-1; i++ {
		for j := i; j >= 0; j-- {
			k := i + (i - j) + 1 // i+1, i+2, ...
			if k == len(b[0]) {
				break
			}
			if !slices.Equal(b.col(j), b.col(k)) {
				continue Loop
			}
		}
		return i, true
	}
	return 0, false
}

// We'll casually assume that there is exactly one index.
func (b block) verticalIndexNearMiss() (int, bool) {
	if len(b) == 0 {
		return 0, false
	}

Loop:
	// The -1 matters as the last column is, by definition, inadmissible as a line
	// of reflection.
	for i := 0; i < len(b[0])-1; i++ {
		nearMisses := 0
		for j := i; j >= 0; j-- {
			k := i + (i - j) + 1 // i+1, i+2, ...
			if k == len(b[0]) {
				break
			}
			if slices.Equal(b.col(j), b.col(k)) {
				// No-op.
			} else if nearMiss(b.col(j), b.col(k)) {
				nearMisses++
			} else {
				continue Loop
			}
		}
		if nearMisses == 1 {
			return i, true
		}
	}
	return 0, false
}

func nearMiss(xs, ys []sitetype) bool {
	assert(len(xs) == len(ys), "can only compare vectors of equal length")
	diffs := 0
	for i := range xs {
		if xs[i] != ys[i] {
			diffs++
		}
	}
	return diffs == 1
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

func parseFile(filepath string) []block {
	bs, err := os.ReadFile(filepath)
	assert(err == nil, "reading a file")

	blockStrs := strings.Split(strings.TrimSpace(string(bs)), "\n\n")
	blocks := make([]block, len(blockStrs))
	for i, blockStr := range blockStrs {
		blocks[i] = toBlock(blockStr)
	}
	return blocks
}

func toBlock(str string) block {
	lines := strings.Split(strings.TrimSpace(str), "\n")
	rows := make([][]sitetype, len(lines))
	for i, line := range lines {
		rows[i] = make([]sitetype, len(line))
		for j, ru := range line {
			switch ru {
			case '.':
				rows[i][j] = ST_ASH
			case '#':
				rows[i][j] = ST_ROCKS
			default:
				panic(fmt.Sprintf("unexpected site type rune: %q", string(ru)))
			}
		}
	}
	assert(equalLen(rows), "expected all rows to have equal length")
	return rows
}

func equalLen[T any](xss [][]T) bool {
	if len(xss) == 0 {
		return true
	}
	l := len(xss[0])
	for _, xs := range xss[1:] {
		if len(xs) != l {
			return false
		}
	}
	return true
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

func assert(b bool, msg string, args ...any) {
	if !b {
		panic("assertion failed: " + fmt.Sprintf(msg, args...))
	}
}
