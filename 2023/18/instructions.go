package main

import (
	"strconv"
	"strings"
)

type instructions []instruction

type instruction struct {
	dir  direction
	dist int
}

type direction int

const (
	DIRECTION_U direction = iota
	DIRECTION_L
	DIRECTION_D
	DIRECTION_R
)

type pipeType int

const (
	PIPE_EMPTY pipeType = iota
	PIPE_STRAIGHT_UD
	PIPE_STRAIGHT_LR
	PIPE_BEND_UL
	PIPE_BEND_DL
	PIPE_BEND_DR
	PIPE_BEND_UR
)

func newInstructions(s string, isPartTwo bool) instructions {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	insts := make(instructions, 0, len(lines))

	for _, line := range lines {
		fields := strings.Fields(line)
		assert(len(fields) == 3)

		var dir direction
		var dist int
		var err error

		if !isPartTwo {
			dirStr, distStr := fields[0], fields[1]
			switch dirStr {
			case "U":
				dir = DIRECTION_U
			case "L":
				dir = DIRECTION_L
			case "D":
				dir = DIRECTION_D
			case "R":
				dir = DIRECTION_R
			default:
				panic("unreachable")
			}

			dist, err = strconv.Atoi(distStr)
			assert(err == nil)
		} else {
			hex := fields[2]
			hex = strings.TrimPrefix(hex, "(#")
			hex = strings.TrimSuffix(hex, ")")
			assert(len(hex) == 6)

			distStr, dirStr := hex[:5], hex[5]
			dist64, err := strconv.ParseInt(distStr, 16, 64)
			assert(err == nil)

			dist = int(dist64)
			switch dirStr {
			case '0':
				dir = DIRECTION_R
			case '1':
				dir = DIRECTION_D
			case '2':
				dir = DIRECTION_L
			case '3':
				dir = DIRECTION_U
			default:
				panic("unexpected encoding of the direction")
			}
		}

		insts = append(insts, instruction{dir: dir, dist: dist})
	}

	return insts
}

func (dir direction) pipeType(nextDir direction) pipeType {
	switch dir {
	case DIRECTION_U:
		switch nextDir {
		case DIRECTION_L:
			return PIPE_BEND_UR
		case DIRECTION_R:
			return PIPE_BEND_UL
		default:
			panic("unexpected direction")
		}
	case DIRECTION_L:
		switch nextDir {
		case DIRECTION_U:
			return PIPE_BEND_DL
		case DIRECTION_D:
			return PIPE_BEND_UL
		default:
			panic("unexpected direction")
		}
	case DIRECTION_D:
		switch nextDir {
		case DIRECTION_L:
			return PIPE_BEND_DR
		case DIRECTION_R:
			return PIPE_BEND_DL
		default:
			panic("unexpected direction")
		}
	case DIRECTION_R:
		switch nextDir {
		case DIRECTION_U:
			return PIPE_BEND_DR
		case DIRECTION_D:
			return PIPE_BEND_UR
		default:
			panic("unexpected direction")
		}
	default:
		panic("unexpected direction")
	}
}

func directionsToPipeType(prev *direction, cur direction) pipeType {
	if prev == nil || *prev == cur {
		return directionToPipeType(cur)
	}
	return prev.pipeType(cur)
}

func directionToPipeType(dir direction) pipeType {
	switch dir {
	case DIRECTION_U, DIRECTION_D:
		return PIPE_STRAIGHT_UD
	case DIRECTION_L, DIRECTION_R:
		return PIPE_STRAIGHT_LR
	default:
		panic("unexpected direction")
	}
}

// dims returns the top-left and bottom-right corners of the smallest grid that
// can carry out the instructions. The bottom-right corner is exclusive.
func (insts instructions) dims() (pair, pair) {
	i, j := 0, 0
	imin, jmin := 0, 0
	imax, jmax := 0, 0
	for _, inst := range insts {
		switch inst.dir {
		case DIRECTION_U:
			i = i - inst.dist
		case DIRECTION_L:
			j = j - inst.dist
		case DIRECTION_D:
			i = i + inst.dist
		case DIRECTION_R:
			j = j + inst.dist
		default:
			panic("unreachable")
		}
		imin, jmin = min(i, imin), min(j, jmin)
		imax, jmax = max(i, imax), max(j, jmax)
	}
	return pair{i: imin, j: jmin}, pair{i: imax + 1, j: jmax + 1}
}
