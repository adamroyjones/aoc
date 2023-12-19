package main

import (
	"strconv"
	"strings"
)

type (
	pairs     []pair
	pair      struct{ i, j int }
	direction int
)

const (
	DIRECTION_U direction = iota
	DIRECTION_L
	DIRECTION_D
	DIRECTION_R
)

func (d direction) opposite() direction {
	switch d {
	case DIRECTION_U:
		return DIRECTION_D
	case DIRECTION_L:
		return DIRECTION_R
	case DIRECTION_D:
		return DIRECTION_U
	case DIRECTION_R:
		return DIRECTION_L
	default:
		panic("unreachable")
	}
}

func (ps pairs) label() string {
	var sb strings.Builder
	sb.WriteString("[")
	for _, p := range ps {
		sb.WriteString("(" + strconv.Itoa(p.i) + "," + strconv.Itoa(p.j) + "),") // Yes, a trailing comma.
	}
	sb.WriteString("]")
	return sb.String()
}

func (fst pair) eq(snd pair) bool { return fst.i == snd.i && fst.j == snd.j }

func (fst pair) direction(snd pair) direction {
	assert(fst.i == snd.i || fst.j == snd.j, "the pairs must equal one another in at least one coordinate (fst: %v, snd: %v)", fst, snd)
	assert(!fst.eq(snd), "the pairs must not be equal")

	if snd.i < fst.i {
		return DIRECTION_U
	}
	if snd.i > fst.i {
		return DIRECTION_D
	}
	if snd.j < fst.j {
		return DIRECTION_L
	}
	return DIRECTION_R
}
