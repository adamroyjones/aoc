package main

import (
	"os"
	"strconv"
	"strings"
)

func parseFile(path string) data {
	bs, err := os.ReadFile(path)
	assert(err == nil)

	blocks := strings.Split(string(bs), "\n\n")

	seedsLn := blocks[0]
	seedsLn = strings.TrimSpace(seedsLn)
	var ok bool
	_, seedsLn, ok = strings.Cut(seedsLn, "seeds: ")
	assert(ok)
	seeds := lineToInts(seedsLn)

	lookups := make([]lookup, 0, len(blocks)-1)
	for _, block := range blocks[1:] {
		lines := strings.Split(strings.TrimSpace(block), "\n")
		assert(len(lines) >= 2)

		lts := make([]lookupTriple, 0, len(lines)-1)
		for _, line := range lines[1:] {
			ints := lineToInts(line)
			assert(len(ints) == 3)
			lts = append(lts, lookupTriple{dst: ints[0], src: ints[1], rng: ints[2]})
		}

		lookups = append(lookups, lookup{lookupTriples: lts})
	}

	return data{seeds: seeds, lookups: lookups}
}

func lineToInts(line string) []int {
	strs := strings.Split(line, " ")
	xs := make([]int, 0, len(strs))
	for _, str := range strs {
		x, err := strconv.Atoi(str)
		assert(err == nil)
		xs = append(xs, x)
	}
	return xs
}
