package main

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type game struct{ time, distance int }

func partOne(filepath string) int {
	games := parseFile(filepath)
	results := make([]int, len(games))
	for i := range games {
		results[i] = winningWays(games[i])
	}
	return product(results...)
}

func partTwo(filepath string) int {
	games := parseFile(filepath)
	var timeStr, distanceStr string
	for _, g := range games {
		timeStr += strconv.Itoa(g.time)
		distanceStr += strconv.Itoa(g.distance)
	}
	time, err := strconv.Atoi(timeStr)
	assert(err == nil, "failed to parse int")
	distance, err := strconv.Atoi(distanceStr)
	assert(err == nil, "failed to parse int")

	f := func(i int) int { return i*(time-i) - distance }
	mid := time / 2
	rng := make([]int, mid+1)
	for i := 0; i <= mid; i++ {
		rng[i] = i
	}
	idx, ok := slices.BinarySearchFunc(rng, 0, func(x, y int) int { return cmp.Compare(f(x), y) })
	assert(!ok, "not ok!")
	return (time - idx) - (idx - 1)
}

func winningWays(g game) int {
	count := 0
	for i := 1; i < g.time; i++ {
		if i*(g.time-i) > g.distance {
			count++
		}
	}
	return count
}

func parseFile(filepath string) []game {
	bs, err := os.ReadFile(filepath)
	assert(err == nil, "could not read file")

	lines := strings.Split(strings.TrimSpace(string(bs)), "\n")
	assert(len(lines) == 2, "len(lines) != 2")

	timeStrs := strings.Fields(strings.TrimPrefix(lines[0], "Time:"))
	distanceStrs := strings.Fields(strings.TrimPrefix(lines[1], "Distance:"))
	assert(len(timeStrs) == len(distanceStrs), "length mismatch")

	games := make([]game, len(timeStrs))
	for i := range timeStrs {
		time, err := strconv.Atoi(timeStrs[i])
		assert(err == nil, "could not parse time %s (fields: %v)", timeStrs[i], timeStrs)
		distance, err := strconv.Atoi(distanceStrs[i])
		assert(err == nil, "could not parse distance")
		games[i] = game{time: time, distance: distance}
	}
	return games
}

func product(xs ...int) int {
	out := 1
	for _, x := range xs {
		out *= x
	}
	return out
}

func assert(b bool, msg string, args ...any) {
	if !b {
		panic("assertion failed: " + fmt.Sprintf(msg, args...))
	}
}
