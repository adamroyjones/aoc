package main

import (
	"fmt"
	"strings"
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
	POS_UD position = iota
	POS_LR
)

type (
	tiletype  int
	position  int
	matrix    [][]tiletype
	point     struct{ i, j int }
	distances [][]int
)

func (m matrix) String() string {
	mapper := func(tt tiletype) string {
		switch tt {
		case TT_START:
			return "S"
		case TT_GROUND:
			return "."
		case TT_STRAIGHT_UD:
			return "|"
		case TT_STRAIGHT_LR:
			return "-"
		case TT_BEND_UL:
			return "╭"
		case TT_BEND_UR:
			return "╮"
		case TT_BEND_DL:
			return "╰"
		case TT_BEND_DR:
			return "╯"
		default:
			panic("unexpected tile type")
		}
	}

	var sb strings.Builder
	for _, row := range m {
		for _, tt := range row {
			sb.WriteString(mapper(tt))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (fst point) eq(snd point) bool { return fst.i == snd.i && fst.j == snd.j }

func start(m matrix) point {
	for i, row := range m {
		for j, cell := range row {
			if cell == TT_START {
				return point{i: i, j: j}
			}
		}
	}
	panic("unable to find the start")
}

func neighbours(m matrix, s point) []point {
	out := []point{}

	if s.i > 0 && compatible(m[s.i-1][s.j], m[s.i][s.j], POS_UD) {
		out = append(out, point{i: s.i - 1, j: s.j})
	}
	if s.i < len(m)-1 && compatible(m[s.i][s.j], m[s.i+1][s.j], POS_UD) {
		out = append(out, point{i: s.i + 1, j: s.j})
	}
	if s.j > 0 && compatible(m[s.i][s.j-1], m[s.i][s.j], POS_LR) {
		out = append(out, point{i: s.i, j: s.j - 1})
	}
	if s.j < len(m[0])-1 && compatible(m[s.i][s.j], m[s.i][s.j+1], POS_LR) {
		out = append(out, point{i: s.i, j: s.j + 1})
	}

	assert(len(out) == 2, "expected exactly 2 neighbours; given %d (point: %v, out: %v)", len(out), s, out)
	return out
}

// If pos is POS_LR, we expect fst to be on the left of snd.
// If pos is POS_UD, we expect fst to be above snd.
func compatible(fst, snd tiletype, pos position) bool {
	switch pos {
	case POS_LR:
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
	case POS_UD:
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

func toMatrix(bs []byte) matrix {
	lines := strings.Split(strings.TrimSpace(string(bs)), "\n")
	m := make([][]tiletype, 0, len(lines))
	for _, line := range lines {
		m = append(m, row(line))
	}
	return m
}

func row(ln string) []tiletype {
	mapper := func(r rune) tiletype {
		switch r {
		case 'S':
			return TT_START
		case '.':
			return TT_GROUND
		case '|':
			return TT_STRAIGHT_UD
		case '-':
			return TT_STRAIGHT_LR
		case 'F':
			return TT_BEND_UL
		case '7':
			return TT_BEND_UR
		case 'L':
			return TT_BEND_DL
		case 'J':
			return TT_BEND_DR
		default:
			panic("unexpected rune")
		}
	}

	out := make([]tiletype, 0, len(ln))
	for _, r := range ln {
		out = append(out, mapper(r))
	}
	return out
}

func assert(b bool, msg string, args ...any) {
	if !b {
		panic("assertion failed: " + fmt.Sprintf(msg, args...))
	}
}
