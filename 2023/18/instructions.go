package main

import (
	"fmt"
	"strconv"
	"strings"
)

type instructions []instruction

type instruction struct {
	dir  direction
	dist int
	col  colour
}

type colour struct {
	r, g, b byte
}

type direction int

const (
	DIRECTION_U direction = iota
	DIRECTION_L
	DIRECTION_D
	DIRECTION_R
)

func (dir direction) opposite() direction {
	switch dir {
	case DIRECTION_U:
		return DIRECTION_D
	case DIRECTION_L:
		return DIRECTION_R
	case DIRECTION_D:
		return DIRECTION_U
	case DIRECTION_R:
		return DIRECTION_L
	default:
		panic("unexpected direction")
	}
}

type pair struct{ i, j int }

func directionsToPipeType(prev *direction, cur direction) pipeType {
	if prev == nil || *prev == cur {
		return directionToPipeType(cur)
	}

	assert(prev.opposite() != cur)

	switch *prev {
	case DIRECTION_U:
		switch cur {
		case DIRECTION_L:
			return PIPE_BEND_UR
		case DIRECTION_R:
			return PIPE_BEND_UL
		default:
			panic("unexpected direction")
		}
	case DIRECTION_L:
		switch cur {
		case DIRECTION_U:
			return PIPE_BEND_DL
		case DIRECTION_D:
			return PIPE_BEND_UL
		default:
			panic("unexpected direction")
		}
	case DIRECTION_D:
		switch cur {
		case DIRECTION_L:
			return PIPE_BEND_DR
		case DIRECTION_R:
			return PIPE_BEND_DL
		default:
			panic("unexpected direction")
		}
	case DIRECTION_R:
		switch cur {
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
		i, j = inst.step(i, j)
		imax, jmax = max(i, imax), max(j, jmax)
		imin, jmin = min(i, imin), min(j, jmin)
	}
	return pair{i: imin, j: jmin}, pair{i: imax + 1, j: jmax + 1}
}

func (inst instruction) step(i, j int) (int, int) {
	switch inst.dir {
	case DIRECTION_U:
		return i - inst.dist, j
	case DIRECTION_L:
		return i, j - inst.dist
	case DIRECTION_D:
		return i + inst.dist, j
	case DIRECTION_R:
		return i, j + inst.dist
	default:
		panic("unreachable")
	}
}

func newInstructions(s string) instructions {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	insts := make(instructions, 0, len(lines))
	for _, line := range lines {
		fields := strings.Fields(line)
		assert(len(fields) == 3)
		dirStr, distStr, colourStr := fields[0], fields[1], fields[2]
		dir := newDirection(dirStr)
		dist, err := strconv.Atoi(distStr)
		assert(err == nil)
		col := newColour(colourStr)
		insts = append(insts, instruction{dir: dir, dist: dist, col: col})
	}
	return insts
}

func newColour(s string) colour {
	s = strings.TrimPrefix(s, "(#")
	s = strings.TrimSuffix(s, ")")
	assertf(len(s) == 6, "string %q should have length 6", s)
	rStr, gStr, bStr := s[:2], s[2:4], s[4:]
	r, err := strconv.ParseUint(rStr, 16, 8)
	assert(err == nil)
	g, err := strconv.ParseUint(gStr, 16, 8)
	assert(err == nil)
	b, err := strconv.ParseUint(bStr, 16, 8)
	assert(err == nil)
	return colour{r: byte(r), g: byte(g), b: byte(b)}
}

func newDirection(s string) direction {
	switch s {
	case "U":
		return DIRECTION_U
	case "L":
		return DIRECTION_L
	case "D":
		return DIRECTION_D
	case "R":
		return DIRECTION_R
	default:
		assert(false)
	}
	panic("unreachable")
}

func (insts instructions) String() string {
	var sb strings.Builder
	for _, inst := range insts {
		sb.WriteString(inst.String())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (inst instruction) String() string {
	return inst.dir.String() + " " + strconv.Itoa(inst.dist) + " " + fmt.Sprintf("%02x/%02x/%02x", inst.col.r, inst.col.g, inst.col.b)
}

func (d direction) String() string {
	switch d {
	case DIRECTION_U:
		return "U"
	case DIRECTION_L:
		return "L"
	case DIRECTION_D:
		return "D"
	case DIRECTION_R:
		return "R"
	default:
		assert(false)
	}
	panic("unreachable")
}
