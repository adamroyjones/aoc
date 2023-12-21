package main

import (
	"slices"
)

type data struct {
	seeds   []int
	lookups []lookup
}

type (
	lookup             struct{ lookupTriples []lookupTriple }
	lookupTriple       struct{ dst, src, rng int }
	leftClosedInterval struct{ start, end int }
	closedInterval     struct{ start, end int }
)

func partOne(filename string) int {
	data := parseFile(filename)
	locations := make([]int, len(data.seeds))
	for i, seed := range data.seeds {
		location := seed
		for _, lookup := range data.lookups {
			location = lookup.next(location)
		}
		locations[i] = location
	}
	return slices.Min(locations)
}

func partTwo(filename string) int {
	data := parseFile(filename)
	assert(len(data.seeds)%2 == 0)

	// Converts the seeds into left-closed intervals.
	intervals := []leftClosedInterval{}
	for i := 0; i < len(data.seeds); i += 2 {
		start := data.seeds[i]
		end := start + data.seeds[i+1]
		intervals = append(intervals, leftClosedInterval{start: start, end: end})
	}

	// For each lookup map:
	//
	// - cut the intervals at the knot points (that is, the transitionary points,
	//   where we switch from mapping one interval to mapping another), and
	// - take the image of the cut intervals under the lookup map.
	//
	// At the end, we can then flatten the intervals and take the minimum.
	for _, lookup := range data.lookups {
		nextIntervals := []leftClosedInterval{}
		for _, interval := range intervals {
			knots := lookup.intervalToKnots(interval)
			srcClosedIntervals := knotsToClosedIntervals(knots)
			dstIntervals := make([]leftClosedInterval, len(srcClosedIntervals))
			for i := range srcClosedIntervals {
				dstIntervals[i] = leftClosedInterval{
					start: lookup.next(srcClosedIntervals[i].start),
					end:   lookup.next(srcClosedIntervals[i].end) + 1,
				}
			}
			nextIntervals = append(nextIntervals, dstIntervals...)
		}
		intervals = nextIntervals
	}

	return slices.Min(flatten(intervals))
}

func (l lookup) next(x int) int {
	for _, lt := range l.lookupTriples {
		if lt.src <= x && x < lt.src+lt.rng {
			return lt.dst + (x - lt.src)
		}
	}
	return x
}

func (l lookup) intervalToKnots(i leftClosedInterval) []int {
	knots := []int{i.start}
	for _, lt := range l.lookupTriples {
		if i.start < lt.src && lt.src < i.end {
			knots = append(knots, lt.src)
		}
		if end := lt.src + lt.rng; i.start < end && end < i.end {
			knots = append(knots, end)
		}
	}
	knots = append(knots, i.end)
	slices.Sort(knots)
	return slices.Compact(knots)
}

func knotsToClosedIntervals(knots []int) []closedInterval {
	endpoints := []int{knots[0]}
	for i := 1; i < len(knots)-1; i++ {
		endpoints = append(endpoints, knots[i]-1, knots[i])
	}
	endpoints = append(endpoints, knots[len(knots)-1]-1)

	intervals := []closedInterval{}
	for i := 0; i < len(endpoints); i += 2 {
		intervals = append(intervals, closedInterval{start: endpoints[i], end: endpoints[i+1]})
	}
	return intervals
}

func flatten(intervals []leftClosedInterval) []int {
	xs := make([]int, 0, 2*len(intervals))
	for _, intv := range intervals {
		xs = append(xs, intv.start, intv.end)
	}
	return xs
}

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}
