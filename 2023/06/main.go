package main

import (
	"cmp"
	"os"
	"slices"
	"strconv"
	"strings"
)

type game struct{ time, distance int }

func partOne(filepath string) int {
	bs, err := os.ReadFile(filepath)
	assert(err == nil)

	lines := strings.Split(strings.TrimSpace(string(bs)), "\n")
	assert(len(lines) == 2)

	timeStrs := strings.Fields(strings.TrimPrefix(lines[0], "Time:"))
	distanceStrs := strings.Fields(strings.TrimPrefix(lines[1], "Distance:"))
	assert(len(timeStrs) == len(distanceStrs))

	product := 1
	for i := range timeStrs {
		time, err := strconv.Atoi(timeStrs[i])
		assert(err == nil)
		distance, err := strconv.Atoi(distanceStrs[i])
		assert(err == nil)

		winningWays := 0
		for i := 1; i < time; i++ {
			if i*(time-i) > distance {
				winningWays++
			}
		}
		product *= winningWays
	}
	return product
}

func partTwo(filepath string) int {
	bs, err := os.ReadFile(filepath)
	assert(err == nil)

	lines := strings.Split(strings.TrimSpace(string(bs)), "\n")
	assert(len(lines) == 2)

	timeStr := lines[0]
	timeStr = strings.ReplaceAll(strings.TrimPrefix(timeStr, "Time:"), " ", "")
	time, err := strconv.Atoi(timeStr)
	assert(err == nil)

	distanceStr := lines[1]
	distanceStr = strings.ReplaceAll(strings.TrimPrefix(distanceStr, "Distance:"), " ", "")
	distance, err := strconv.Atoi(distanceStr)
	assert(err == nil)

	fn := func(i int) int { return i*(time-i) - distance }
	mid := time / 2
	rng := make([]int, mid+1)
	for i := 0; i <= mid; i++ {
		rng[i] = i
	}
	idx, ok := slices.BinarySearchFunc(rng, 0, func(x, y int) int { return cmp.Compare(fn(x), y) })
	assert(!ok) // We could also handle the ok case, but I don't want to.
	return (time - idx) - (idx - 1)
}

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}
