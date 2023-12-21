package main

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
	"time"
)

type (
	lightGrid     [][]lightGridCell
	lightGridCell struct{ dirs []direction }
	direction     int
)

type ray struct {
	i, j int
	dir  direction
}

const (
	DIRECTION_U direction = iota
	DIRECTION_L
	DIRECTION_D
	DIRECTION_R
)

func newLightGrid(mg mirrorGrid, start ray, print bool) lightGrid {
	assert(len(mg) > 0)
	assert(len(mg[0]) > 0)
	assert(allEqualLen(mg))

	lg := make(lightGrid, len(mg))
	for i := range lg {
		lg[i] = make([]lightGridCell, len(mg[0]))
	}

	rays := mg.firstStep(start)
	for _, r := range rays {
		lg[r.i][r.j].dirs = append(lg[r.i][r.j].dirs, r.dir)
	}

	for len(rays) > 0 {
		lg, rays = lg.step(mg, rays)
		if print {
			lg.printGrid(mg, rays)
			time.Sleep(50 * time.Millisecond)
		}
	}
	return lg
}

func (lg lightGrid) step(mg mirrorGrid, rays []ray) (lightGrid, []ray) {
	nextRays := []ray{}
	for _, r := range rays {
		nextRays = append(nextRays, mg.stepRay(r)...)
	}

	// Deduplicate the next rays.
	slices.SortFunc(nextRays, func(fst, snd ray) int { return fst.cmp(snd) })
	nextRays = slices.CompactFunc(nextRays, func(fst, snd ray) bool { return fst.cmp(snd) == 0 })

	// Loop erasure: remove the rays that have already been seen.
	nextRays = slices.DeleteFunc(nextRays, func(r ray) bool { return slices.Contains(lg[r.i][r.j].dirs, r.dir) })
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

func (lg lightGrid) printGrid(mg mirrorGrid, rays []ray) {
	fmt.Printf("\033[0;0H\n")

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
			switch count {
			case 0:
				if len(lg[i][j].dirs) > 0 {
					sb.WriteString("#")
				} else {
					sb.WriteString(mg[i][j].String())
				}
			case 1:
				sb.WriteString(rays[rayIdx].dir.String())
			default:
				sb.WriteString("2")
			}
		}
		sb.WriteString("\n")
	}

	fmt.Println(sb.String())
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

func (lgc lightGridCell) String() string {
	if len(lgc.dirs) == 0 {
		return "."
	}
	return "#"
}

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
