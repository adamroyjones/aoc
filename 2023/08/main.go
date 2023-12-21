package main

import (
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
	pair     struct{ left, right string }
)

const (
	movementL movement = "L"
	movementR movement = "R"
)

func partOne(filepath string) int {
	d := parseFile(filepath)
	next, count := "AAA", 0
	for {
		p := d.mapping[next]
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
			return count
		}
	}
}

func partTwo(filepath string) int {
	d := parseFile(filepath)
	starts := maps.Keys(d.mapping)
	starts = slices.DeleteFunc(starts, func(s string) bool { return !strings.HasSuffix(s, "A") })

	// start -> Z-suffixed sites -> [visitations]
	periodicData := make(map[string]map[string][]int)

	for _, start := range starts {
		// Z-suffixed sites -> [visitations]
		periodicDatum := map[string][]int{}
		count := 0
		next := start
	NavigationLoop:
		for {
			p := d.mapping[next]
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

			visitations, ok := periodicDatum[next]
			if !ok {
				periodicDatum[next] = []int{count}
				continue
			}

			periodicDatum[next] = append(visitations, count)
			// As soon as we have periodicity, we break.
			for _, visitation := range visitations {
				if (count-visitation)%len(d.movements) == 0 {
					periodicData[start] = periodicDatum
					break NavigationLoop
				}
			}
		}
	}

	// The below allows us to specialise the calculation. It asserts that we only
	// ever visit one Z-suffixed site per start position. It also asserts that the
	// periodic behaviour is about as simple as it gets: as soon as we visit the
	// site, we've entered periodic behaviour.
	xs := make([]int, 0, len(periodicData))
	for _, periodicDatum := range periodicData {
		assert(len(periodicDatum) == 1)
		for _, visitations := range periodicDatum {
			for i := 1; i < len(visitations); i++ {
				assert(visitations[0]*(i+1) == visitations[i])
			}
			xs = append(xs, visitations[0])
		}
	}

	return lcm(xs)
}

func parseFile(filepath string) data {
	bs, err := os.ReadFile(filepath)
	assert(err == nil)
	blocks := strings.Split(strings.TrimSpace(string(bs)), "\n\n")
	assert(len(blocks) == 2)
	return data{movements: parseMovements(blocks[0]), mapping: parseMapping(blocks[1])}
}

func parseMovements(s string) []movement {
	s = strings.TrimSpace(s)
	movements := make([]movement, 0, len(s))
	for _, ru := range s {
		switch ru {
		case 'L':
			movements = append(movements, movementL)
		case 'R':
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
		assert(ok)
		after = strings.TrimPrefix(after, "(")
		after = strings.TrimSuffix(after, ")")
		left, right, ok := strings.Cut(after, ", ")
		assert(ok)
		out[before] = pair{left: left, right: right}
	}
	return out
}

func lcm(xs []int) int {
	assert(len(xs) > 0)
	if len(xs) == 1 {
		return xs[0]
	}
	if len(xs) == 2 {
		return (xs[0] * xs[1]) / gcd(xs[0], xs[1])
	}
	return lcm(append([]int{lcm(xs[:2])}, xs[2:]...))
}

func gcd(a, b int) int {
	// Euclid's algorithm.
	for b != 0 {
		tmp := b
		b = a % b
		a = tmp
	}
	return a
}

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}
