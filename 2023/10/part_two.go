package main

import (
	"os"
	"slices"
)

func partTwo(filepath string) int {
	bs, err := os.ReadFile(filepath)
	assert(err == nil, "reading a file")
	m := toMatrix(bs)
	assert(len(m) > 0, "expected the matrix to be non-empty")
	loop := tidyLoop(m)
	return area(loop)
}

func tidyLoop(m matrix) matrix {
	loop := make(matrix, len(m))
	for i := range loop {
		loop[i] = make([]tiletype, len(m[0]))
		for j := range loop[i] {
			loop[i][j] = TT_GROUND
		}
	}

	s := start(m)
	nrs := neighbours(m, s)
	prev, current := s, nrs[0]
	loop[prev.i][prev.j] = m[prev.i][prev.j]
	loop[current.i][current.j] = m[current.i][current.j]
	for {
		nextNeighbours := neighbours(m, current)
		idx := slices.IndexFunc(nextNeighbours, func(p point) bool { return !p.eq(prev) })
		assert(idx != -1, "expected to find a new neighbour")

		nextNeighbour := nextNeighbours[idx]
		if nextNeighbour.eq(s) {
			break
		}

		prev, current = current, nextNeighbour
		loop[current.i][current.j] = m[current.i][current.j]
	}

	// Mutates loop.
	fixupStart(loop, s)
	loop = padLoop(loop)
	return loop
}

func fixupStart(loop matrix, s point) {
	nbrs := neighbours(loop, s)
	assert(len(nbrs) == 2, "expected 2 neighbours")

	startPtr := &loop[s.i][s.j]

	fst, snd := nbrs[0], nbrs[1]
	if fst.i == snd.i {
		*startPtr = TT_STRAIGHT_LR
		return
	}
	if fst.j == snd.j {
		*startPtr = TT_STRAIGHT_UD
		return
	}

	if fst.j < s.j || snd.j < s.j { // One to the left.
		if fst.i < s.i || snd.i < s.i { // One above.
			*startPtr = TT_BEND_DR
		} else { // One below.
			*startPtr = TT_BEND_UR
		}
	} else { // One to the right.
		if fst.i < s.i || snd.i < s.i { // One above.
			*startPtr = TT_BEND_DL
		} else { // One below.
			*startPtr = TT_BEND_UL
		}
	}
}

func padLoop(loop matrix) matrix {
	nextLoop := make(matrix, len(loop)+2)
	for i := range nextLoop {
		row := make([]tiletype, len(loop[0])+2)
		if i == 0 || i == len(nextLoop)-1 {
			for j := range row {
				row[j] = TT_GROUND
			}
		} else {
			for j := range row {
				if j == 0 || j == len(row)-1 {
					row[j] = TT_GROUND
				} else {
					row[j] = loop[i-1][j-1]
				}
			}
		}

		nextLoop[i] = row
	}
	return nextLoop
}

func area(loop matrix) int {
	a := 0
	for _, row := range loop {
		a += areaOnLine(row)
	}
	return a
}

func areaOnLine(line []tiletype) int {
	outside := true
	area := 0
	var prevBend *tiletype
	for _, cell := range line {
		switch cell {
		case TT_GROUND:
			if !outside {
				area++
			}
		case TT_STRAIGHT_UD:
			outside = !outside
		case TT_STRAIGHT_LR:
			continue
		case TT_BEND_UL:
			prevBend = ptr(TT_BEND_UL)
		case TT_BEND_UR:
			if prevBend == nil {
				panic("unexpected situation")
			}
			if *prevBend == TT_BEND_DL {
				outside = !outside
			}
			prevBend = nil
		case TT_BEND_DL:
			prevBend = ptr(TT_BEND_DL)
		case TT_BEND_DR:
			if prevBend == nil {
				panic("unexpected situation")
			}
			if *prevBend == TT_BEND_UL {
				outside = !outside
			}
			prevBend = nil
		default:
			panic("unexpected tile type")
		}
	}

	return area
}

func ptr[T any](t T) *T { return &t }
