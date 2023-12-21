package main

import (
	"os"
	"slices"
)

func partOne(filepath string) int {
	bs, err := os.ReadFile(filepath)
	assert(err == nil)

	m := newMatrix(bs)
	assert(len(m) > 0)

	start := m.start()
	dss := make([]distances, 0, 2)
	for _, neighbour := range m.neighbours(start) {
		// Initialise the distance matrix to 0.
		ds := make(distances, len(m))
		for rowIdx := range ds {
			ds[rowIdx] = make([]int, len(m[0]))
		}

		// The neighbour is at a distance of 1 from the start.
		prev, current := start, neighbour
		distance := 1
		ds[current.i][current.j] = distance

		// We continue the walk start -> neighbour until we loop around to start once again.
		for {
			nextNeighbours := m.neighbours(current)
			idx := slices.IndexFunc(nextNeighbours, func(p point) bool { return !p.eq(prev) })
			assert(idx != -1)
			nextNeighbour := nextNeighbours[idx]
			if nextNeighbour.eq(start) {
				break
			}

			prev, current = current, nextNeighbour
			distance++
			ds[current.i][current.j] = distance
		}

		dss = append(dss, ds)
	}

	// We now find the maximin.
	assert(len(dss) == 2)
	d := 0
	fst, snd := dss[0], dss[1]
	assert(len(fst) == len(snd))
	for i := 0; i < len(fst); i++ {
		assert(len(fst[i]) == len(snd[i]))
		for j := 0; j < len(fst[i]); j++ {
			d = max(d, min(fst[i][j], snd[i][j]))
		}
	}
	return d
}

func partTwo(filepath string) int {
	bs, err := os.ReadFile(filepath)
	assert(err == nil)
	m := newMatrix(bs)
	assert(len(m) > 0)
	m.tidy()
	return m.area()
}

func ptr[T any](t T) *T { return &t }

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}
