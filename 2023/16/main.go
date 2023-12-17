package main

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

type direction int

const (
	DIRECTION_U direction = iota
	DIRECTION_L
	DIRECTION_D
	DIRECTION_R
)

func (d direction) String() string {
	switch d {
	case DIRECTION_U:
		return "^"
	case DIRECTION_L:
		return "<"
	case DIRECTION_D:
		return "v"
	case DIRECTION_R:
		return ">"
	default:
		panic("unexpected direction")
	}
}

type (
	mirrorCellType int
	mirrorGrid     [][]mirrorCellType
)

const (
	MIRROR_EMPTY mirrorCellType = iota
	MIRROR_VERTICAL
	MIRROR_HORIZONTAL
	MIRROR_FORWARD_SLASH
	MIRROR_BACKSLASH
)

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

func (mg mirrorGrid) String() string {
	var sb strings.Builder
	for _, row := range mg {
		for _, cell := range row {
			sb.WriteString(cell.String())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

type (
	lightGridCell struct{ dirs []direction }
	lightGrid     [][]lightGridCell
)

func (lgc lightGridCell) String() string {
	if len(lgc.dirs) == 0 {
		return "."
	}
	return "#"
}

func (lg lightGrid) String() string {
	var sb strings.Builder
	for _, row := range lg {
		for _, lgc := range row {
			sb.WriteString(lgc.String())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (lg lightGrid) count() int {
	count := 0
	for i := range lg {
		for j := range lg[0] {
			if len(lg[i][j].dirs) > 0 {
				count++
			}
		}
	}
	return count
}

type ray struct {
	i, j int
	dir  direction
}

func main() {
	// This illustrates the evolution. This requires an ANSI terminal.
	filename := "integration-part-1"
	filename = "input"
	bs, err := os.ReadFile(filename)
	assert(err == nil, "unable to read a file")
	mg := toMirrorGrid(string(bs))
	start := ray{i: 0, j: -1, dir: DIRECTION_R}
	_ = toLightGrid(mg, start, true)
}

func partTwo(filename string) int {
	bs, err := os.ReadFile(filename)
	assert(err == nil, "unable to read a file")
	mg := toMirrorGrid(string(bs))
	count := 0
	for _, ray := range startingRays(mg) {
		lg := toLightGrid(mg, ray, false)
		count = max(count, lg.count())
	}
	return count
}

func partOne(filename string) int {
	bs, err := os.ReadFile(filename)
	assert(err == nil, "unable to read a file")
	mg := toMirrorGrid(string(bs))
	start := ray{i: 0, j: -1, dir: DIRECTION_R}
	lg := toLightGrid(mg, start, false)
	return lg.count()
}

func toMirrorGrid(s string) mirrorGrid {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	assert(len(lines) > 0, "expected at least one line")
	assert(allEqualLenStrs(lines), "not a rectangle")

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

func startingRays(mg mirrorGrid) []ray {
	rays := []ray{}
	for i := range mg {
		rays = append(rays, ray{i: i, j: -1, dir: DIRECTION_R}, ray{i: i, j: len(mg[0]), dir: DIRECTION_L})
	}
	for j := range mg[0] {
		rays = append(rays, ray{i: -1, j: j, dir: DIRECTION_D}, ray{i: len(mg), j: j, dir: DIRECTION_U})
	}
	return rays
}

func toLightGrid(mg mirrorGrid, start ray, print bool) lightGrid {
	rays := firstStep(mg, start)

	lg := newLightGrid(mg, rays)
	for len(rays) > 0 {
		lg, rays = step(mg, lg, rays)
		if print {
			printGrid(mg, lg, rays)
			time.Sleep(16 * time.Millisecond)
		}
	}
	return lg
}

func firstStep(mg mirrorGrid, start ray) []ray {
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

func printGrid(mg mirrorGrid, lg lightGrid, rays []ray) {
	fmt.Printf("\033[0;0H")

	var sb strings.Builder
	for i := range lg {
		for j := range lg[0] {
			count := 0
			rayIdx := 0
			for k, r := range rays {
				if r.i == i && r.j == j {
					rayIdx = k
					count++
				}
			}
			if count >= 2 {
				sb.WriteString("2")
			}
			if count == 1 {
				sb.WriteString(rays[rayIdx].dir.String())
			}
			if count == 0 {
				if len(lg[i][j].dirs) > 0 {
					sb.WriteString("#")
				} else {
					sb.WriteString(mg[i][j].String())
				}
			}
		}
		sb.WriteString("\n")
	}

	fmt.Println(sb.String())
}

func newLightGrid(mg mirrorGrid, rays []ray) lightGrid {
	assert(len(mg) > 0, "expected at least one row in the mirror grid")
	assert(len(mg[0]) > 0, "expected at least one column in the mirror grid")
	assert(allEqualLen(mg), "expected a rectangular mirror grid")
	lg := make([][]lightGridCell, len(mg))
	for i := range lg {
		lg[i] = make([]lightGridCell, len(mg[0]))
	}

	for _, r := range rays {
		lg[r.i][r.j].dirs = append(lg[r.i][r.j].dirs, r.dir)
	}

	return lg
}

func step(mg mirrorGrid, lg lightGrid, rays []ray) (lightGrid, []ray) {
	nextRays := []ray{}
	for _, r := range rays {
		nextRays = append(nextRays, stepRay(mg, r)...)
	}

	slices.SortFunc(nextRays, func(fst, snd ray) int { return fst.cmp(snd) })
	nextRays = slices.CompactFunc(nextRays, func(fst, snd ray) bool { return fst.cmp(snd) == 0 })
	nextRays = slices.DeleteFunc(nextRays, func(r ray) bool { return lg.containsRay(r) })
	for _, r := range nextRays {
		lg[r.i][r.j].dirs = append(lg[r.i][r.j].dirs, r.dir)
	}

	return lg, nextRays
}

func (fst ray) cmp(snd ray) int {
	if fst.i != snd.i {
		return cmp.Compare(fst.i, snd.i)
	}
	if fst.j != snd.j {
		return cmp.Compare(fst.j, snd.j)
	}
	return cmp.Compare(fst.dir, snd.dir)
}

func (lg lightGrid) containsRay(r ray) bool {
	lgc := lg[r.i][r.j]
	return slices.Contains(lgc.dirs, r.dir)
}

func stepRay(mg mirrorGrid, r ray) []ray {
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

func allEqualLen[T any](xss [][]T) bool {
	if len(xss) == 0 {
		return true
	}
	l := len(xss[0])
	for _, xs := range xss[1:] {
		if l != len(xs) {
			return false
		}
	}
	return true
}

func allEqualLenStrs(ss []string) bool {
	if len(ss) == 0 {
		return true
	}
	l := len(ss[0])
	for _, s := range ss[1:] {
		if l != len(s) {
			return false
		}
	}
	return true
}

func assert(b bool, msg string, args ...any) {
	if !b {
		panic("assertion failed: " + fmt.Sprintf(msg, args...))
	}
}
