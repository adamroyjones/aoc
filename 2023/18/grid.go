package main

import (
	"strings"
)

type grid struct {
	cells      [][]cell
	imin, jmin int // Inclusive.
	imax, jmax int // Exclusive.
}

type cell struct {
	cellTyp cellType
	pipeTyp pipeType
}

type cellType int

const (
	CELL_EMPTY cellType = iota
	CELL_TRENCH
)

func newGrid(insts instructions) *grid {
	assert(len(insts) > 0)
	ul, br := insts.dims()
	grydCells := make([][]cell, br.i-ul.i)
	for i := range grydCells {
		grydCells[i] = repeat(cell{cellTyp: CELL_EMPTY}, br.j-ul.j)
	}
	gryd := &grid{cells: grydCells, imin: ul.i, jmin: ul.j, imax: br.i, jmax: br.j}

	// Dig out the trenches and add pipe types.
	position := pair{i: 0, j: 0}
	var prevDirection *direction
	for _, inst := range insts {
		prevDirection, position = gryd.step(prevDirection, position, inst)
	}

	// We need to fix the origin; the origin's pipe type can only be known at the end of the traversal.
	assert(prevDirection != nil)
	gryd.cells[-gryd.imin][-gryd.jmin].pipeTyp = prevDirection.pipeType(insts[0].dir)

	return gryd
}

func (gryd grid) area() int { return gryd.boundaryArea() + gryd.boundedArea() }

func (gryd grid) boundaryArea() int {
	area := 0
	for i := range gryd.cells {
		for j := range gryd.cells[i] {
			if gryd.cells[i][j].cellTyp == CELL_TRENCH {
				area++
			}
		}
	}
	return area
}

func (gryd grid) boundedArea() int {
	area := 0
	for _, grydLine := range gryd.cells {
		area += gryd.areaOnLine(grydLine)
	}
	return area
}

func (gryd grid) areaOnLine(grydLine []cell) int {
	outside := true
	area := 0

	var prevBend *pipeType
	for _, cell := range grydLine {
		switch cell.pipeTyp {
		case PIPE_EMPTY:
			if !outside {
				area++
			}
		case PIPE_STRAIGHT_UD:
			outside = !outside
		case PIPE_STRAIGHT_LR:
			continue
		case PIPE_BEND_UL:
			prevBend = ptr(PIPE_BEND_UL)
		case PIPE_BEND_UR:
			if prevBend == nil {
				panic("unexpected situation")
			}
			if *prevBend == PIPE_BEND_DL {
				outside = !outside
			}
			prevBend = nil
		case PIPE_BEND_DL:
			prevBend = ptr(PIPE_BEND_DL)
		case PIPE_BEND_DR:
			if prevBend == nil {
				panic("unexpected situation")
			}
			if *prevBend == PIPE_BEND_UL {
				outside = !outside
			}
			prevBend = nil
		default:
			panic("unexpected pipe type")
		}
	}

	return area
}

// step takes a point and an instruction and executes the instruction, returning
// a new Cartesian point. The point MUST be remapped appropriately by this
// method to produce appropriate slice indices.
func (gryd *grid) step(prevDirection *direction, pos pair, inst instruction) (*direction, pair) {
	arrayI, arrayJ := pos.i-gryd.imin, pos.j-gryd.jmin
	nextPos := pos
	gryd.cells[arrayI][arrayJ].pipeTyp = directionsToPipeType(prevDirection, inst.dir)

	switch inst.dir {
	case DIRECTION_U:
		for i := 1; i <= inst.dist; i++ {
			gryd.cells[arrayI-i][arrayJ].cellTyp = CELL_TRENCH
			gryd.cells[arrayI-i][arrayJ].pipeTyp = directionToPipeType(inst.dir)
		}
		nextPos.i -= inst.dist
	case DIRECTION_L:
		for j := 1; j <= inst.dist; j++ {
			gryd.cells[arrayI][arrayJ-j].cellTyp = CELL_TRENCH
			gryd.cells[arrayI][arrayJ-j].pipeTyp = directionToPipeType(inst.dir)
		}
		nextPos.j -= inst.dist
	case DIRECTION_D:
		for i := 1; i <= inst.dist; i++ {
			gryd.cells[arrayI+i][arrayJ].cellTyp = CELL_TRENCH
			gryd.cells[arrayI+i][arrayJ].pipeTyp = directionToPipeType(inst.dir)
		}
		nextPos.i += inst.dist
	case DIRECTION_R:
		for j := 1; j <= inst.dist; j++ {
			gryd.cells[arrayI][arrayJ+j].cellTyp = CELL_TRENCH
			gryd.cells[arrayI][arrayJ+j].pipeTyp = directionToPipeType(inst.dir)
		}
		nextPos.j += inst.dist
	default:
		panic("unexpected direction")
	}

	assert(nextPos.i >= gryd.imin && nextPos.i < gryd.imax)
	assert(nextPos.j >= gryd.jmin && nextPos.j < gryd.jmax)
	return &inst.dir, nextPos
}

func (gryd grid) trenchString() string {
	var sb strings.Builder
	for _, row := range gryd.cells {
		for _, c := range row {
			sb.WriteString(c.cellTyp.String())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (gryd grid) pipeString() string {
	var sb strings.Builder
	for _, row := range gryd.cells {
		for _, c := range row {
			sb.WriteString(c.pipeTyp.String())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (c cellType) String() string {
	switch c {
	case CELL_EMPTY:
		return "."
	case CELL_TRENCH:
		return "#"
	default:
		panic("unexpected cell")
	}
}

func (pt pipeType) String() string {
	switch pt {
	case PIPE_EMPTY:
		return "."
	case PIPE_STRAIGHT_UD:
		return "|"
	case PIPE_STRAIGHT_LR:
		return "-"
	case PIPE_BEND_UL:
		return "╭"
	case PIPE_BEND_DL:
		return "╰"
	case PIPE_BEND_DR:
		return "╯"
	case PIPE_BEND_UR:
		return "╮"
	default:
		panic("unexpected pipe type")
	}
}
