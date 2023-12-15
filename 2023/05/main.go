package main

import (
	"fmt"
	"slices"
)

func partOne(filename string) int {
	data := parseFile(filename)
	locations := make([]int, len(data.seeds))
	for i, seed := range data.seeds {
		locations[i] = data.lookups.toLocation(seed)
	}
	return slices.Min(locations)
}

func partTwo(filename string) int {
	data := parseFile(filename)
	intervals := []leftClosedInterval{}
	assert(len(data.seeds)%2 == 0, "expected an even number of seeds")

	// Converts the seeds into left-closed intervals.
	for i := 0; i < len(data.seeds); i += 2 {
		start := data.seeds[i]
		end := start + data.seeds[i+1]
		intervals = append(intervals, leftClosedInterval{start: start, end: end})
	}

	for _, lookup := range data.lookups {
		nextIntervals := []leftClosedInterval{}
		for _, intv := range intervals {
			knots := lookup.intervalToKnots(intv)
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

type data struct {
	seeds   []int
	lookups lookups
}

type (
	lookups            []lookup
	lookup             struct{ lookupTriples []lookupTriple }
	lookupTriple       struct{ dst, src, rng int }
	leftClosedInterval struct{ start, end int }
	closedInterval     struct{ start, end int }
)

func (l lookups) toLocation(x int) int {
	for _, lookup := range l {
		x = lookup.next(x)
	}
	return x
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

func assert(b bool, msg string, args ...any) {
	if !b {
		panic("assertion failed: " + fmt.Sprintf(msg, args...))
	}
}
