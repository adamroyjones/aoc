package main

import (
	"slices"
	"strings"
)

type (
	tiletype  int
	direction int
	matrix    [][]tiletype

	point struct{ i, j int }

	// distances is a map i -> j -> distance. The value is the distance from (i, j)
	// to the creature. If we're outside of the loop, the distance is set to be 0.
	// Ultimately we'll be maximising the distance, so this is fine.
	distances [][]int
)

const (
	TT_START tiletype = iota
	TT_GROUND
	TT_STRAIGHT_UD
	TT_STRAIGHT_LR
	TT_BEND_UL
	TT_BEND_UR
	TT_BEND_DL
	TT_BEND_DR
)

const (
	DIR_D direction = iota
	DIR_R
)

func (fst point) eq(snd point) bool { return fst.i == snd.i && fst.j == snd.j }

func newMatrix(bs []byte) matrix {
	lines := strings.Split(strings.TrimSpace(string(bs)), "\n")
	m := make([][]tiletype, 0, len(lines))
	for _, line := range lines {
		row := make([]tiletype, 0, len(line))
		for _, ru := range line {
			var tt tiletype
			switch ru {
			case 'S':
				tt = TT_START
			case '.':
				tt = TT_GROUND
			case '|':
				tt = TT_STRAIGHT_UD
			case '-':
				tt = TT_STRAIGHT_LR
			case 'F':
				tt = TT_BEND_UL
			case '7':
				tt = TT_BEND_UR
			case 'L':
				tt = TT_BEND_DL
			case 'J':
				tt = TT_BEND_DR
			default:
				panic("unexpected rune")
			}
			row = append(row, tt)
		}

		m = append(m, row)
	}

	return m
}

// tidy mutates m, even though it's not a pointer receiver.
func (m matrix) tidy() {
	witnessed := make(map[point]struct{})

	s := m.start()
	neighbours := m.neighbours(s)
	assert(len(neighbours) == 2)
	prev, current := s, neighbours[0]
	witnessed[prev], witnessed[current] = struct{}{}, struct{}{}
	for {
		nextNeighbours := m.neighbours(current)
		idx := slices.IndexFunc(nextNeighbours, func(p point) bool { return !p.eq(prev) })
		assert(idx != -1)

		nextNeighbour := nextNeighbours[idx]
		if nextNeighbour.eq(s) {
			break
		}

		prev, current = current, nextNeighbour
		witnessed[current] = struct{}{}
	}

	for i := range m {
		row := m[i]
		for j := range row {
			if _, ok := witnessed[point{i: i, j: j}]; !ok {
				cell := &row[j]
				*cell = TT_GROUND
			}
		}
	}

	// Fix up the start.
	startPtr := &m[s.i][s.j]
	fst, snd := neighbours[0], neighbours[1]
	if fst.i == snd.i {
		*startPtr = TT_STRAIGHT_LR
		return
	}
	if fst.j == snd.j {
		*startPtr = TT_STRAIGHT_UD
		return
	}

	if fst.j < s.j || snd.j < s.j {
		if fst.i < s.i || snd.i < s.i {
			*startPtr = TT_BEND_DR
		} else {
			*startPtr = TT_BEND_UR
		}
		return
	}

	if fst.i < s.i || snd.i < s.i {
		*startPtr = TT_BEND_DL
	} else {
		*startPtr = TT_BEND_UL
	}
}

func (m matrix) start() point {
	for i, row := range m {
		for j, cell := range row {
			if cell == TT_START {
				return point{i: i, j: j}
			}
		}
	}
	panic("unable to find the start")
}

func (m matrix) neighbours(s point) []point {
	out := []point{}

	if s.i > 0 && compatible(m[s.i-1][s.j], m[s.i][s.j], DIR_D) {
		out = append(out, point{i: s.i - 1, j: s.j})
	}
	if s.i < len(m)-1 && compatible(m[s.i][s.j], m[s.i+1][s.j], DIR_D) {
		out = append(out, point{i: s.i + 1, j: s.j})
	}
	if s.j > 0 && compatible(m[s.i][s.j-1], m[s.i][s.j], DIR_R) {
		out = append(out, point{i: s.i, j: s.j - 1})
	}
	if s.j < len(m[0])-1 && compatible(m[s.i][s.j], m[s.i][s.j+1], DIR_R) {
		out = append(out, point{i: s.i, j: s.j + 1})
	}

	assert(len(out) == 2)
	return out
}

// area presupposes that the matrix represents a loop. The area is calculated as
// though it were a trivial kind of Fubini's theorem.
func (m matrix) area() int {
	a := 0
	for _, row := range m {
		a += areaOnLine(row)
	}
	return a
}

// areaOnLine calculates the interior area lying on the line. It uses a crossing
// argument. We start on the left, on the outside of the loop, and move to the
// right. When we cross the loop, we move into the inside of the loop and start
// counting tiles. If we cross it again, we're on the outside, and so we stop
// counting, etc.
//
// The only wrinkle is around what constitutes a crossing.
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

// compatible takes in two tile types and an orientation and says whether they
// could be neighbours in a loop.
//
// If dir is DIR_R, we expect fst to be on the left of snd.
// If dir is DIR_D, we expect fst to be above snd.
func compatible(fst, snd tiletype, dir direction) bool {
	switch dir {
	case DIR_R:
		switch fst {
		case TT_START:
			return snd == TT_STRAIGHT_LR || snd == TT_BEND_DR || snd == TT_BEND_UR
		case TT_STRAIGHT_LR, TT_BEND_DL, TT_BEND_UL:
			return snd == TT_STRAIGHT_LR || snd == TT_BEND_DR || snd == TT_BEND_UR || snd == TT_START
		case TT_GROUND, TT_STRAIGHT_UD, TT_BEND_DR, TT_BEND_UR:
			return false
		default:
			panic("unexpected tile type")
		}
	case DIR_D:
		switch fst {
		case TT_START:
			return snd == TT_STRAIGHT_UD || snd == TT_BEND_DL || snd == TT_BEND_DR
		case TT_STRAIGHT_UD, TT_BEND_UL, TT_BEND_UR:
			return snd == TT_STRAIGHT_UD || snd == TT_BEND_DL || snd == TT_BEND_DR || snd == TT_START
		case TT_GROUND, TT_STRAIGHT_LR, TT_BEND_DL, TT_BEND_DR:
			return false
		default:
			panic("unexpected tile type")
		}
	default:
		panic("unexpected position")
	}
}
