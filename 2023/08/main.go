package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"golang.org/x/exp/maps"
)

type data struct {
	movements []movement
	mapping   map[string]pair
}

type (
	movement string
	location string
	pair     struct{ left, right string }
)

const (
	movementL movement = "L"
	movementR movement = "R"
)

func main() {
	d := partTwo("input")
	fmt.Printf("%d\n", d)
}

func partTwo(filepath string) int {
	d := parseFile(filepath)
	filter := func(s string) bool { return !strings.HasSuffix(s, "A") }
	starts := maps.Keys(d.mapping)
	starts = slices.DeleteFunc(starts, filter)

	// start -> Z-suffixed sites -> [visitations]
	periodicData := make(map[string]map[string][]int)

	for _, start := range starts {
		m := map[string][]int{}
		p := pair{}
		count := 0
		next := start
	NavigationLoop:
		for {
			p = d.mapping[next]
			switch movement := d.movements[count%len(d.movements)]; movement {
			case movementL:
				next = p.left
			case movementR:
				next = p.right
			default:
				panic("unexpected movement")
			}

			count++
			if !strings.HasSuffix(next, "Z") {
				continue
			}

			visitations, ok := m[next]
			if !ok {
				m[next] = []int{count}
				continue
			}

			m[next] = append(visitations, count)
			for _, visitation := range visitations {
				if (count-visitation)%len(d.movements) == 0 {
					periodicData[start] = m
					break NavigationLoop
				}
			}
		}
	}

	// The below allows us to specialise the calculation.
	xs := make([]int, 0, len(periodicData))
	for _, visitations := range periodicData {
		assert(len(visitations) == 1, "len != 1")
		for _, counts := range visitations {
			for i := 1; i < len(counts); i++ {
				assert(counts[0]*(i+1) == counts[i], "expected a_n = an")
			}
			xs = append(xs, counts[0])
		}
	}

	return lcm(xs)
}

func partOne(filepath string) int {
	d := parseFile(filepath)
	var p pair
	next := "AAA"
	count := 0
	for {
		p = d.mapping[next]
		switch movement := d.movements[count%len(d.movements)]; movement {
		case movementL:
			next = p.left
		case movementR:
			next = p.right
		default:
			panic("unexpected movement")
		}

		count++
		if next == "ZZZ" {
			break
		}
	}

	return count
}

func parseFile(filepath string) data {
	bs, err := os.ReadFile(filepath)
	assert(err == nil, "failed to read a file")
	blocks := strings.Split(strings.TrimSpace(string(bs)), "\n\n")
	assert(len(blocks) == 2, "expected 2 blocks")
	movements := parseMovements(blocks[0])
	mapping := parseMapping(blocks[1])
	return data{movements: movements, mapping: mapping}
}

func parseMovements(ln string) []movement {
	chars := strings.Split(strings.TrimSpace(ln), "")
	movements := make([]movement, 0, len(chars))
	for _, char := range chars {
		switch char {
		case "L":
			movements = append(movements, movementL)
		case "R":
			movements = append(movements, movementR)
		default:
			panic("invalid movement")
		}
	}
	return movements
}

func parseMapping(block string) map[string]pair {
	lines := strings.Split(strings.TrimSpace(block), "\n")
	out := make(map[string]pair)
	for _, line := range lines {
		before, after, ok := strings.Cut(line, " = ")
		assert(ok, "expected things to be ok")
		after = strings.TrimPrefix(after, "(")
		after = strings.TrimSuffix(after, ")")
		left, right, ok := strings.Cut(after, ", ")
		assert(ok, "expected things to be ok")
		out[before] = pair{left: left, right: right}
	}
	return out
}

func lcm(xs []int) int {
	assert(len(xs) > 0, "non-empty slice expected")
	if len(xs) == 1 {
		return xs[0]
	}
	if len(xs) == 2 {
		return (xs[0] * xs[1]) / gcd(xs[0], xs[1])
	}

	return lcm(append([]int{lcm(xs[:2])}, xs[2:]...))
}

func gcd(a, b int) int {
	for b != 0 { // Euclid
		t := b
		b = a % b
		a = t
	}
	return a
}

func assert(b bool, msg string, args ...any) {
	if !b {
		panic("assertion failed: " + fmt.Sprintf(msg, args...))
	}
}
