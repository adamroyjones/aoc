package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

type (
	grid      [][]tiletype
	tiletype  int
	direction int
)

const (
	TT_EMPTY tiletype = iota
	TT_ROUND_ROCK
	TT_CUBIC_ROCK
)

const (
	DIRECTION_N direction = iota
	DIRECTION_W
	DIRECTION_S
	DIRECTION_E
)

func main() {
	// This illustrates the periodic behaviour. This requires an ANSI terminal.
	filename := "testdata/integration_part_2"
	grid := parseFile(filename)
	for i := 0; i <= 1000; i++ {
		fmt.Printf("\033[0;0H\n")
		fmt.Println(grid)
		grid = cycle(grid)
		time.Sleep(100 * time.Millisecond)
	}
}

func partOne(filename string) int {
	grid := parseFile(filename)
	finalGrid := slide(grid, DIRECTION_N)
	return score(finalGrid)
}

func partTwo(filename string) int {
	g := parseFile(filename)
	history := []grid{g}

	cycleLimit := 1_000_000_000
	periodStart := 0
	periodEnd := 0
Loop:
	for i := 0; i < cycleLimit; i++ {
		g = cycle(g)
		history = append(history, g)
		for j, gg := range history[:len(history)-1] {
			if g.eq(gg) {
				periodStart, periodEnd = j, i+1
				break Loop
			}
		}
	}

	// Express cycleLimit = periodStart + K(periodEnd-periodStart) + r and then
	// extract the (periodStart + r)th element from the history.
	r := (cycleLimit - periodStart) % (periodEnd - periodStart)
	g = history[periodStart+r]
	return score(g)
}

func parseFile(filename string) grid {
	bs, err := os.ReadFile(filename)
	assert(err == nil, "unable to read a file")
	lines := strings.Split(strings.TrimSpace(string(bs)), "\n")
	grid := make([][]tiletype, len(lines))
	for i, line := range lines {
		grid[i] = newGridRow(line)
	}
	return grid
}

func score(g grid) int {
	s := 0
	for i := len(g) - 1; i >= 0; i-- {
		count := 0
		for _, tt := range g[i] {
			if tt == TT_ROUND_ROCK {
				count++
			}
		}
		s += count * (len(g) - i) // The second factor is 1, 2, ...
	}
	return s
}

func cycle(g grid) grid {
	return slide(slide(slide(slide(g, DIRECTION_N), DIRECTION_W), DIRECTION_S), DIRECTION_E)
}

func slide(g grid, dir direction) grid {
	switch dir {
	case DIRECTION_N, DIRECTION_S:
		return slideGridNS(g, dir)
	case DIRECTION_E, DIRECTION_W:
		return slideGridEW(g, dir)
	default:
		panic("unexpected direction")
	}
}

func slideGridNS(g grid, dir direction) grid {
	nextGrid := transpose(g)
	switch dir {
	case DIRECTION_N:
		return transpose(slideGridEW(nextGrid, DIRECTION_W))
	case DIRECTION_S:
		return transpose(slideGridEW(nextGrid, DIRECTION_E))
	default:
		panic("unexpected direction")
	}
}

func slideGridEW(g grid, dir direction) grid {
	nextGrid := make([][]tiletype, len(g))
	for i := range g {
		nextGrid[i] = slideRowEW(g[i], dir)
	}
	return nextGrid
}

func slideRowEW(row []tiletype, dir direction) []tiletype {
	if len(row) == 0 {
		return row
	}

	if row[0] == TT_CUBIC_ROCK {
		return append([]tiletype{TT_CUBIC_ROCK}, slideRowEW(row[1:], dir)...)
	}

	if idx := slices.Index(row, TT_CUBIC_ROCK); idx != -1 {
		return append(slideRowEW(row[:idx], dir), slideRowEW(row[idx:], dir)...)
	}

	roundRockCount := 0
	for _, tt := range row {
		if tt == TT_ROUND_ROCK {
			roundRockCount++
		}
	}
	switch dir {
	case DIRECTION_W:
		return append(repeat(TT_ROUND_ROCK, roundRockCount), repeat(TT_EMPTY, len(row)-roundRockCount)...)
	case DIRECTION_E:
		return append(repeat(TT_EMPTY, len(row)-roundRockCount), repeat(TT_ROUND_ROCK, roundRockCount)...)
	default:
		panic("unexpected direction")
	}
}

func transpose(g grid) grid {
	imax, jmax := len(g), len(g[0])
	transposed := make([][]tiletype, jmax)
	for i := range transposed {
		transposed[i] = make([]tiletype, imax)
	}

	for i := range g {
		for j := range g[i] {
			transposed[j][i] = g[i][j]
		}
	}

	return transposed
}

func newGridRow(line string) []tiletype {
	out := make([]tiletype, len(line))
	for i, ru := range line {
		switch ru {
		case '.':
			out[i] = TT_EMPTY
		case 'O':
			out[i] = TT_ROUND_ROCK
		case '#':
			out[i] = TT_CUBIC_ROCK
		default:
			panic("unexpected tile type")
		}
	}
	return out
}

func (xss grid) eq(yss grid) bool {
	fn := func(xs, ys []tiletype) bool { return slices.Equal(xs, ys) }
	return slices.EqualFunc(xss, yss, fn)
}

func repeat[T any](x T, count int) []T {
	out := make([]T, count)
	for i := range out {
		out[i] = x
	}
	return out
}

func assert(b bool, msg string, args ...any) {
	if !b {
		panic("assertion failed: " + fmt.Sprintf(msg, args...))
	}
}

func (g grid) String() string {
	var sb strings.Builder
	for _, row := range g {
		for _, tt := range row {
			sb.WriteString(tt.String())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (tt tiletype) String() string {
	switch tt {
	case TT_EMPTY:
		return "."
	case TT_ROUND_ROCK:
		return "O"
	case TT_CUBIC_ROCK:
		return "#"
	default:
		panic("unexpected tile type")
	}
}
