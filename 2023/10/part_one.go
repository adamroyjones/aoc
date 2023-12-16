package main

import (
	"os"
	"slices"
)

func partOne(filepath string) int {
	bs, err := os.ReadFile(filepath)
	assert(err == nil, "reading a file")
	m := toMatrix(bs)
	assert(len(m) > 0, "expected the matrix to be non-empty")
	s := start(m)
	nrs := neighbours(m, s)
	dss := []distances{}
	for _, nr := range nrs {
		ds := make(distances, len(m))
		for rowIdx := range ds {
			ds[rowIdx] = make([]int, len(m[0]))
		}

		prev, current := s, nr
		distance := 1
		ds[current.i][current.j] = distance

		for {
			nextNeighbours := neighbours(m, current)
			idx := slices.IndexFunc(nextNeighbours, func(p point) bool { return !p.eq(prev) })
			assert(idx != -1, "expected to find a new neighbour")
			nextNeighbour := nextNeighbours[idx]
			if nextNeighbour.eq(s) {
				break
			}

			prev, current = current, nextNeighbour
			distance++
			ds[current.i][current.j] = distance
		}

		dss = append(dss, ds)
	}

	assert(len(dss) == 2, "expected 2 distance matrices")
	return maximin(dss[0], dss[1])
}

func maximin(fst, snd distances) int {
	d := 0
	imax := len(fst)
	jmax := len(fst[0])
	for i := 0; i < imax; i++ {
		for j := 0; j < jmax; j++ {
			d = max(d, min(fst[i][j], snd[i][j]))
		}
	}
	return d
}
