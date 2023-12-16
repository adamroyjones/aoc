package main

import (
	"cmp"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

var (
	orderPartOne = []card{"A", "K", "Q", "J", "T", "9", "8", "7", "6", "5", "4", "3", "2"}
	orderPartTwo = []card{"A", "K", "Q", "T", "9", "8", "7", "6", "5", "4", "3", "2", "J"}
)

type record struct {
	cards cards
	bid   int
	rank  int
}

type (
	cards []card
	card  string
)

func (fst record) cmp(snd record, partOne bool) int {
	if fst.rank != snd.rank {
		return cmp.Compare(fst.rank, snd.rank)
	}
	return fst.cards.cmp(snd.cards, partOne)
}

func (fst cards) cmp(snd cards, partOne bool) int {
	for i := range fst {
		if fst[i] != snd[i] {
			return fst[i].cmp(snd[i], partOne)
		}
	}
	return 0
}

func (fst card) cmp(snd card, partOne bool) int {
	if partOne {
		return cmp.Compare(slices.Index(orderPartOne, snd), slices.Index(orderPartOne, fst))
	}
	return cmp.Compare(slices.Index(orderPartTwo, snd), slices.Index(orderPartTwo, fst))
}

func main() {
	partTwo("input")
}

func partTwo(filename string) int {
	partOne := false
	return run(filename, partOne)
}

func partOne(filename string) int {
	partOne := true
	return run(filename, partOne)
}

func run(filename string, partOne bool) int {
	rs := parseFile(filename, partOne)
	slices.SortFunc(rs, func(fst, snd record) int { return fst.cmp(snd, partOne) })
	result := 0
	for i, r := range rs {
		result += (i + 1) * r.bid
	}
	return result
}

func parseFile(filepath string, partOne bool) []record {
	bs, err := os.ReadFile(filepath)
	assert(err == nil, "unable to read file")

	rows := strings.Split(strings.TrimSpace(string(bs)), "\n")
	records := make([]record, len(rows))
	for i := range rows {
		records[i] = rowToRecord(rows[i], partOne)
	}
	return records
}

func rowToRecord(row string, partOne bool) record {
	fields := strings.Fields(row)
	assert(len(fields) == 2, "expected 2 fields (row: %s)", row)
	cardsJoinedStr, bidStr := fields[0], fields[1]
	bid, err := strconv.Atoi(bidStr)
	assert(err == nil, "unable to parse int")
	cardsStrs := strings.Split(cardsJoinedStr, "")
	assert(len(cardsStrs) == 5, "expected 5 cards in a hand")
	cards := make([]card, len(cardsStrs))
	for i := range cardsStrs {
		cards[i] = card(cardsStrs[i])
	}
	rank := cardsToRank(cards, partOne)
	return record{cards: cards, bid: bid, rank: rank}
}

func cardsToRank(cs []card, partOne bool) int {
	counts := make(map[card]int)
	for _, c := range cs {
		counts[c] = counts[c] + 1
	}

	if partOne {
		return cardsToRankPartOne(counts)
	}
	return cardsToRankPartTwo(counts)
}

func cardsToRankPartTwo(counts map[card]int) int {
	jokerCount := counts[card("J")]
	if jokerCount == 0 {
		return cardsToRankPartOne(counts)
	}
	if jokerCount == 5 { // To avoid slices.Max of an empty slice.
		return 7
	}

	delete(counts, card("J"))
	countOfCounts := make(map[int]int)
	for _, count := range counts {
		countOfCounts[count] = countOfCounts[count] + 1
	}

	maxcnt := slices.Max(maps.Keys(countOfCounts))
	switch maxcnt + jokerCount {
	case 5:
		return 7
	case 4:
		return 6
	case 3:
		switch maxcnt {
		case 3: // jokerCount must be zero.
			if _, ok := countOfCounts[2]; ok {
				return 5 // Full house.
			}
			return 4
		case 2:
			if count := countOfCounts[2]; count == 2 {
				return 5 // Full house.
			}
			return 4
		case 1:
			return 4
		default:
			panic("unreachable")
		}
	case 2:
		switch maxcnt {
		case 2:
			if count := countOfCounts[2]; count == 2 {
				return 3
			}
			return 2
		case 1:
			return 2
		default:
			panic("unreachable")
		}
	case 1:
		return 1
	default:
		panic("unreachable")
	}

	panic("unreachable")
}

func cardsToRankPartOne(counts map[card]int) int {
	countOfCounts := make(map[int]int)
	for _, count := range counts {
		countOfCounts[count] = countOfCounts[count] + 1
	}

	switch maxcnt := slices.Max(maps.Keys(countOfCounts)); maxcnt {
	case 5, 4:
		return maxcnt + 2
	case 3:
		if _, ok := countOfCounts[2]; ok {
			return 5 // full house
		}
		return 4
	case 2:
		if count := countOfCounts[2]; count == 2 {
			return 3 // two pair
		}
		return 2
	case 1:
		return 1
	default:
		panic("unexpected value")
	}
}

func assert(b bool, msg string, args ...any) {
	if !b {
		panic("assertion failed: " + fmt.Sprintf(msg, args...))
	}
}
