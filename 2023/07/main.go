package main

import (
	"cmp"
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

func partOne(filename string) int {
	partOne := true
	return run(filename, partOne)
}

func partTwo(filename string) int {
	partOne := false
	return run(filename, partOne)
}

func run(filename string, partOne bool) int {
	records := parseFile(filename, partOne)
	slices.SortFunc(records, func(fst, snd record) int { return fst.cmp(snd, partOne) })
	result := 0
	for i, r := range records {
		result += (i + 1) * r.bid
	}
	return result
}

func parseFile(filepath string, partOne bool) []record {
	bs, err := os.ReadFile(filepath)
	assert(err == nil)

	rows := strings.Split(strings.TrimSpace(string(bs)), "\n")
	records := make([]record, len(rows))
	for i := range rows {
		records[i] = rowToRecord(rows[i], partOne)
	}
	return records
}

func rowToRecord(row string, partOne bool) record {
	fields := strings.Fields(row)
	assert(len(fields) == 2)

	cardsJoinedStr, bidStr := fields[0], fields[1]
	bid, err := strconv.Atoi(bidStr)
	assert(err == nil)

	cardsStrs := strings.Split(cardsJoinedStr, "")
	assert(len(cardsStrs) == 5)

	cards := make([]card, len(cardsStrs))
	for i := range cardsStrs {
		cards[i] = card(cardsStrs[i])
	}
	rank := cardsToRank(cards, partOne)
	return record{cards: cards, bid: bid, rank: rank}
}

func cardsToRank(cs []card, partOne bool) int {
	cardToCount := make(map[card]int)
	for _, c := range cs {
		cardToCount[c] = cardToCount[c] + 1
	}

	if partOne {
		return cardsToRankPartOne(cardToCount)
	}
	return cardsToRankPartTwo(cardToCount)
}

func cardsToRankPartOne(cardToCount map[card]int) int {
	cardCountToCount := make(map[int]int)
	for _, count := range cardToCount {
		cardCountToCount[count] = cardCountToCount[count] + 1
	}

	switch maxcnt := slices.Max(maps.Keys(cardCountToCount)); maxcnt {
	case 5:
		return 7
	case 4:
		return 6
	case 3:
		if _, ok := cardCountToCount[2]; ok {
			return 5 // Full house.
		}
		return 4
	case 2:
		if count := cardCountToCount[2]; count == 2 {
			return 3 // Two-pair.
		}
		return 2
	case 1:
		return 1
	default:
		panic("unexpected value")
	}
}

func cardsToRankPartTwo(cardToCount map[card]int) int {
	jokerCount := cardToCount[card("J")]
	if jokerCount == 0 {
		return cardsToRankPartOne(cardToCount)
	}
	// This ensures we avoid calling slices.Max on an empty slice.
	if jokerCount == 5 {
		return 7
	}

	delete(cardToCount, card("J"))
	cardCountToCounts := make(map[int]int)
	for _, count := range cardToCount {
		cardCountToCounts[count] = cardCountToCounts[count] + 1
	}

	maxcnt := slices.Max(maps.Keys(cardCountToCounts))
	// jokerCount is never zero given the special handling at the top.
	switch maxcnt + jokerCount {
	case 5:
		return 7
	case 4:
		return 6
	case 3:
		switch maxcnt {
		case 2:
			if count := cardCountToCounts[2]; count == 2 {
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
}

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}
