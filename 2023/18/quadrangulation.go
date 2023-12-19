package main

import (
	"slices"

	"golang.org/x/exp/maps"
)

type (
	quad           struct{ irange, jrange closedInterval }
	closedInterval struct{ min, max int }
)

func newQuadrangulation(sg *sparseGrid) *[]quad {
	assert(sg != nil)

	iset, jset := make(map[int]struct{}, len(sg.edges)), make(map[int]struct{}, len(sg.edges))
	for _, sge := range sg.edges {
		iset[sge.start.i] = struct{}{}
		jset[sge.start.j] = struct{}{}
	}

	islice, jslice := maps.Keys(iset), maps.Keys(jset)
	slices.Sort(islice)
	slices.Sort(jslice)

	// Given, say, {0, 1, 4}, this will return the intervals [0, 0], [1, 1], [2, 3], [4, 4].
	sliceToIntervals := func(xs []int) []closedInterval {
		assert(len(xs) >= 2)
		ranges := make([]closedInterval, 0, len(xs))
		for i, x := range xs[:len(xs)-1] {
			ranges = append(ranges, closedInterval{min: x, max: x})
			if xs[i+1] > x+1 {
				ranges = append(ranges, closedInterval{min: x + 1, max: xs[i+1] - 1})
			}
		}
		ranges = append(ranges, closedInterval{min: xs[len(xs)-1], max: xs[len(xs)-1]})
		return ranges
	}

	iranges, jranges := sliceToIntervals(islice), sliceToIntervals(jslice)

	quads := make([]quad, 0, (len(iranges) * len(jranges)))
	for _, irange := range iranges {
		for _, jrange := range jranges {
			quads = append(quads, quad{irange: irange, jrange: jrange})
		}
	}
	return &quads
}

func (q quad) area() int {
	return (q.jrange.max - q.jrange.min + 1) * (q.irange.max - q.irange.min + 1)
}
