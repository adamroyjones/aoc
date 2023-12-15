package main

import (
	"os"
	"strconv"
	"strings"
)

func parseFile(path string) data {
	bs, err := os.ReadFile(path)
	assert(err == nil, "error reading file: %v", err)

	blocks := strings.Split(string(bs), "\n\n")

	seedsLn := blocks[0]
	seedsLn = strings.TrimSpace(seedsLn)
	var ok bool
	_, seedsLn, ok = strings.Cut(seedsLn, "seeds: ")
	assert(ok, "'seeds: ' not found (line: %s)", seedsLn)
	seeds := lineToInts(seedsLn)

	lookups := make([]lookup, 0, len(blocks)-1)
	for _, block := range blocks[1:] {
		lookup := blockToLookup(block)
		lookups = append(lookups, lookup)
	}

	return data{seeds: seeds, lookups: lookups}
}

func blockToLookup(str string) lookup {
	lines := strings.Split(strings.TrimSpace(str), "\n")
	assert(len(lines) >= 2, "len(lines) = %d < 2", len(lines))

	lts := make([]lookupTriple, 0, len(lines)-1)
	for _, line := range lines[1:] {
		ints := lineToInts(line)
		assert(len(ints) == 3, "len(ints) = %d != 3", len(ints))
		lts = append(lts, lookupTriple{dst: ints[0], src: ints[1], rng: ints[2]})
	}

	return lookup{lookupTriples: lts}
}

func lineToInts(line string) []int {
	strs := strings.Split(line, " ")
	xs := make([]int, 0, len(strs))
	for _, str := range strs {
		x, err := strconv.Atoi(str)
		assert(err == nil, "failed to convert to %s to an int: %v", str, err)
		xs = append(xs, x)
	}
	return xs
}
