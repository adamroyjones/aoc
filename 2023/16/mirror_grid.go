package main

import "strings"

type (
	mirrorGrid     [][]mirrorCellType
	mirrorCellType int
)

const (
	MIRROR_EMPTY mirrorCellType = iota
	MIRROR_VERTICAL
	MIRROR_HORIZONTAL
	MIRROR_FORWARD_SLASH
	MIRROR_BACKSLASH
)

func newMirrorGrid(s string) mirrorGrid {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	assert(len(lines) > 0)
	assert(allEqualLenStrs(lines))

	cells := make([][]mirrorCellType, len(lines))
	for i := range cells {
		cells[i] = make([]mirrorCellType, len(lines[0]))
	}

	for i, line := range lines {
		for j, ru := range line {
			switch ru {
			case '.':
				cells[i][j] = MIRROR_EMPTY
			case '|':
				cells[i][j] = MIRROR_VERTICAL
			case '-':
				cells[i][j] = MIRROR_HORIZONTAL
			case '/':
				cells[i][j] = MIRROR_FORWARD_SLASH
			case '\\':
				cells[i][j] = MIRROR_BACKSLASH
			default:
				panic("unexpected rune")
			}
		}
	}

	return cells
}

func (mg mirrorGrid) startingRays() []ray {
	rays := []ray{}
	for i := range mg {
		rays = append(rays, ray{i: i, j: -1, dir: DIRECTION_R}, ray{i: i, j: len(mg[0]), dir: DIRECTION_L})
	}
	for j := range mg[0] {
		rays = append(rays, ray{i: -1, j: j, dir: DIRECTION_D}, ray{i: len(mg), j: j, dir: DIRECTION_U})
	}
	return rays
}

func (mg mirrorGrid) firstStep(start ray) []ray {
	// Firing down from the top.
	if start.i == -1 {
		switch mg[0][start.j] {
		case MIRROR_EMPTY:
			return []ray{{i: 0, j: start.j, dir: DIRECTION_D}}
		case MIRROR_VERTICAL:
			return []ray{{i: 0, j: start.j, dir: DIRECTION_D}}
		case MIRROR_HORIZONTAL:
			return []ray{{i: 0, j: start.j, dir: DIRECTION_L}, {i: 0, j: start.j, dir: DIRECTION_R}}
		case MIRROR_FORWARD_SLASH:
			return []ray{{i: 0, j: start.j, dir: DIRECTION_L}}
		case MIRROR_BACKSLASH:
			return []ray{{i: 0, j: start.j, dir: DIRECTION_R}}
		default:
			panic("unexpected mirror")
		}
	}

	// Firing up from the bottom.
	if start.i == len(mg) {
		switch mg[0][start.j] {
		case MIRROR_EMPTY:
			return []ray{{i: len(mg) - 1, j: start.j, dir: DIRECTION_U}}
		case MIRROR_VERTICAL:
			return []ray{{i: len(mg) - 1, j: start.j, dir: DIRECTION_U}}
		case MIRROR_HORIZONTAL:
			return []ray{{i: len(mg) - 1, j: start.j, dir: DIRECTION_L}, {i: len(mg) - 1, j: start.j, dir: DIRECTION_R}}
		case MIRROR_FORWARD_SLASH:
			return []ray{{i: len(mg) - 1, j: start.j, dir: DIRECTION_R}}
		case MIRROR_BACKSLASH:
			return []ray{{i: len(mg) - 1, j: start.j, dir: DIRECTION_L}}
		default:
			panic("unexpected mirror")
		}
	}

	// Firing right from the left.
	if start.j == -1 {
		switch mg[start.i][0] {
		case MIRROR_EMPTY:
			return []ray{{i: start.i, j: 0, dir: DIRECTION_R}}
		case MIRROR_VERTICAL:
			return []ray{{i: start.i, j: 0, dir: DIRECTION_U}, {i: start.i, j: 0, dir: DIRECTION_D}}
		case MIRROR_HORIZONTAL:
			return []ray{{i: start.i, j: 0, dir: DIRECTION_R}}
		case MIRROR_FORWARD_SLASH:
			return []ray{{i: start.i, j: 0, dir: DIRECTION_U}}
		case MIRROR_BACKSLASH:
			return []ray{{i: start.i, j: 0, dir: DIRECTION_D}}
		default:
			panic("unexpected mirror")
		}
	}

	// Firing left from the right.
	if start.j == len(mg[0]) {
		switch mg[start.i][len(mg[0])-1] {
		case MIRROR_EMPTY:
			return []ray{{i: start.i, j: len(mg[0]) - 1, dir: DIRECTION_L}}
		case MIRROR_VERTICAL:
			return []ray{{i: start.i, j: len(mg[0]) - 1, dir: DIRECTION_U}, {i: start.i, j: 0, dir: DIRECTION_D}}
		case MIRROR_HORIZONTAL:
			return []ray{{i: start.i, j: len(mg[0]) - 1, dir: DIRECTION_L}}
		case MIRROR_FORWARD_SLASH:
			return []ray{{i: start.i, j: len(mg[0]) - 1, dir: DIRECTION_D}}
		case MIRROR_BACKSLASH:
			return []ray{{i: start.i, j: len(mg[0]) - 1, dir: DIRECTION_U}}
		default:
			panic("unexpected mirror")
		}
	}

	panic("unexpected starting ray")
}

func (mg mirrorGrid) stepRay(r ray) []ray {
	switch r.dir {
	case DIRECTION_U:
		if r.i == 0 {
			return []ray{}
		}

		i, j := r.i-1, r.j
		switch mgc := mg[i][j]; mgc {
		case MIRROR_EMPTY:
			return []ray{{i: i, j: j, dir: DIRECTION_U}}
		case MIRROR_VERTICAL:
			return []ray{{i: i, j: j, dir: DIRECTION_U}}
		case MIRROR_HORIZONTAL:
			return []ray{{i: i, j: j, dir: DIRECTION_L}, {i: i, j: j, dir: DIRECTION_R}}
		case MIRROR_FORWARD_SLASH:
			return []ray{{i: i, j: j, dir: DIRECTION_R}}
		case MIRROR_BACKSLASH:
			return []ray{{i: i, j: j, dir: DIRECTION_L}}
		default:
			panic("unexpected mirror")
		}

	case DIRECTION_L:
		if r.j == 0 {
			return []ray{}
		}

		i, j := r.i, r.j-1
		switch mgc := mg[i][j]; mgc {
		case MIRROR_EMPTY:
			return []ray{{i: i, j: j, dir: DIRECTION_L}}
		case MIRROR_VERTICAL:
			return []ray{{i: i, j: j, dir: DIRECTION_U}, {i: i, j: j, dir: DIRECTION_D}}
		case MIRROR_HORIZONTAL:
			return []ray{{i: i, j: j, dir: DIRECTION_L}}
		case MIRROR_FORWARD_SLASH:
			return []ray{{i: i, j: j, dir: DIRECTION_D}}
		case MIRROR_BACKSLASH:
			return []ray{{i: i, j: j, dir: DIRECTION_U}}
		default:
			panic("unexpected mirror")
		}

	case DIRECTION_D:
		if r.i == len(mg)-1 {
			return []ray{}
		}

		i, j := r.i+1, r.j
		switch mgc := mg[i][j]; mgc {
		case MIRROR_EMPTY:
			return []ray{{i: i, j: j, dir: DIRECTION_D}}
		case MIRROR_VERTICAL:
			return []ray{{i: i, j: j, dir: DIRECTION_D}}
		case MIRROR_HORIZONTAL:
			return []ray{{i: i, j: j, dir: DIRECTION_L}, {i: i, j: j, dir: DIRECTION_R}}
		case MIRROR_FORWARD_SLASH:
			return []ray{{i: i, j: j, dir: DIRECTION_L}}
		case MIRROR_BACKSLASH:
			return []ray{{i: i, j: j, dir: DIRECTION_R}}
		default:
			panic("unexpected mirror")
		}

	case DIRECTION_R:
		if r.j == len(mg[0])-1 {
			return []ray{}
		}

		i, j := r.i, r.j+1
		switch mgc := mg[i][j]; mgc {
		case MIRROR_EMPTY:
			return []ray{{i: i, j: j, dir: DIRECTION_R}}
		case MIRROR_VERTICAL:
			return []ray{{i: i, j: j, dir: DIRECTION_U}, {i: i, j: j, dir: DIRECTION_D}}
		case MIRROR_HORIZONTAL:
			return []ray{{i: i, j: j, dir: DIRECTION_R}}
		case MIRROR_FORWARD_SLASH:
			return []ray{{i: i, j: j, dir: DIRECTION_U}}
		case MIRROR_BACKSLASH:
			return []ray{{i: i, j: j, dir: DIRECTION_D}}
		default:
			panic("unexpected mirror")
		}

	default:
		panic("unexpected direction")
	}
}

func (mct mirrorCellType) String() string {
	switch mct {
	case MIRROR_EMPTY:
		return "."
	case MIRROR_VERTICAL:
		return "|"
	case MIRROR_HORIZONTAL:
		return "-"
	case MIRROR_FORWARD_SLASH:
		return "/"
	case MIRROR_BACKSLASH:
		return `\`
	default:
		panic("unexpected mirror cell type")
	}
}
