package main

import (
	"fmt"
	"slices"
	"strings"
)

type (
	block    [][]sitetype
	sitetype int
)

const (
	ST_ASH   sitetype = 0
	ST_ROCKS sitetype = 1
)

func newBlock(str string) block {
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
	assert(equalLen(rows))
	return rows
}

func (b block) col(j int) []sitetype {
	out := make([]sitetype, len(b))
	for i := range b {
		out[i] = b[i][j]
	}
	return out
}

// horizontalIndex returns the index of the row of reflection. We'll casually
// assume that there is exactly zero or one index.
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

// horizontalIndexNearMiss returns the index of the near-miss variant of the row
// of reflection; that is, a row that is one modification (elsewhere) from being
// a genuine row of reflection. We'll casually assume that there is exactly zero
// or one index.
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
				continue
			}
			if !nearMiss(b[j], b[k]) {
				continue Loop
			}
			nearMisses++
		}
		if nearMisses == 1 {
			return i, true
		}
	}
	return 0, false
}

// verticalIndex returns the index of the column of reflection. We'll casually
// assume that there is exactly zero or one index.
func (b block) verticalIndex() (int, bool) {
	if len(b) == 0 {
		return 0, false
	}

Loop:
	// The -1 matters as the last column is, by definition, inadmissible as a line
	// of reflection.
	for i := 0; i < len(b[0])-1; i++ {
		for j := i; j >= 0; j-- {
			// And so k = i+1, i+2, ...
			k := i + (i - j) + 1
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

// verticalIndexNearMiss returns the index of the near-miss variant of the
// column of reflection; that is, a row that is one modification (elsewhere)
// from being a genuine column of reflection. We'll casually assume that there
// is exactly zero or one index.
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
			// And so k = i+1, i+2, ...
			k := i + (i - j) + 1
			if k == len(b[0]) {
				break
			}
			if slices.Equal(b.col(j), b.col(k)) {
				continue
			}
			if !nearMiss(b.col(j), b.col(k)) {
				continue Loop
			}
			nearMisses++
		}
		if nearMisses == 1 {
			return i, true
		}
	}
	return 0, false
}

func nearMiss(xs, ys []sitetype) bool {
	assert(len(xs) == len(ys))
	diffs := 0
	for i := range xs {
		if xs[i] != ys[i] {
			diffs++
		}
	}
	return diffs == 1
}
