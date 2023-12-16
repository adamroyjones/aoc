package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type row struct {
	rawStr string
	chunks []string
	counts []int
}

var cache = map[string]map[string]int{}

func cacheKeys(r row) (string, string) {
	fstKey := r.rawStr
	if r.rawStr == "" {
		fstKey = strings.Join(r.chunks, ".")
	}

	var sb strings.Builder
	for i, c := range r.counts {
		sb.WriteString(strconv.Itoa(c))
		if i < len(r.counts)-1 {
			sb.WriteString(",")
		}
	}

	return fstKey, sb.String()
}

func get(fstKey, sndKey string) (int, bool) {
	fstMap, ok := cache[fstKey]
	if !ok {
		return 0, false
	}
	val, ok := fstMap[sndKey]
	return val, ok
}

func put(fstKey, sndKey string, val int) int {
	fstMap, ok := cache[fstKey]
	if !ok {
		cache[fstKey] = make(map[string]int)
		fstMap = cache[fstKey]
	}
	fstMap[sndKey] = val
	return val
}

func partTwo(filename string) int {
	rows := filenameToRows(filename)
	rows = unfold(rows, 5)
	arrs := rowsToArrangements(rows)
	return sum(arrs)
}

func partOne(filename string) int {
	rows := filenameToRows(filename)
	arrs := rowsToArrangements(rows)
	return sum(arrs)
}

func filenameToRows(filename string) []row {
	bs, err := os.ReadFile(filename)
	assert(err == nil, "reading a file")
	lines := strings.Split(strings.TrimSpace(string(bs)), "\n")
	rows := make([]row, len(lines))
	for i, line := range lines {
		rawStr, countsStr, ok := strings.Cut(line, " ")
		assert(ok, "expected the line to contain a space")
		counts := strToInts(countsStr)
		rows[i] = row{rawStr: rawStr, counts: counts}
	}
	return rows
}

func rowsToArrangements(rows []row) []int {
	arrs := make([]int, len(rows))
	for i, row := range rows {
		row.rawStr = string(slices.CompactFunc([]byte(row.rawStr), func(fst, snd byte) bool { return fst == '.' && fst == snd }))
		arrs[i] = rowToArrangements(row)
	}
	return arrs
}

// rowToArrangements is the meat of the logic. The use of caching is critical
// for this to work quickly enough to be viable.
func rowToArrangements(r row) int {
	fstKey, sndKey := cacheKeys(r)
	val, ok := get(fstKey, sndKey)
	if ok {
		return val
	}

	chunks := r.chunks
	if chunks == nil {
		chunks = strings.Split(r.rawStr, ".")
		chunks = slices.DeleteFunc(chunks, func(s string) bool { return s == "" })
	}
	chunkLengths := make([]int, len(chunks))
	for i, chunk := range chunks {
		chunkLengths[i] = len(chunk)
	}

	// We begin with some basic validations.
	if sum(chunkLengths) < sum(r.counts) {
		return put(fstKey, sndKey, 0)
	}
	if sum(chunkLengths) == sum(r.counts) {
		if slices.Equal(chunkLengths, r.counts) {
			return put(fstKey, sndKey, 1)
		}
		return put(fstKey, sndKey, 0)
	}
	if sum(r.counts) == 0 {
		if slices.ContainsFunc(chunks, func(s string) bool { return strings.ContainsRune(s, '#') }) {
			return put(fstKey, sndKey, 0)
		}
		return put(fstKey, sndKey, 1)
	}

	// The chunks slice must have more available sites than the sum of counts. We
	// now need to do some work.
	//
	// Must we ignore the first chunk entirely?
	if chunkLengths[0] < r.counts[0] {
		if strings.ContainsRune(chunks[0], '#') {
			return put(fstKey, sndKey, 0)
		}
		val := rowToArrangements(row{chunks: chunks[1:], counts: r.counts})
		return put(fstKey, sndKey, val)
	}

	// Can we ignore the first chunk entirely?
	if chunkLengths[0] == r.counts[0] {
		useFirstChunk := rowToArrangements(row{chunks: chunks[1:], counts: r.counts[1:]})
		omitFirstChunk := 0
		if !strings.ContainsRune(chunks[0], '#') {
			omitFirstChunk = rowToArrangements(row{chunks: chunks[1:], counts: r.counts})
		}

		return put(fstKey, sndKey, useFirstChunk+omitFirstChunk)
	}

	// The first chunk is longer than r.counts[0]. We'll have to consider making
	// repeated use of the first chunk. We'll work recursively from the first
	// element of the first chunk.
	//
	// dotCase tries replacing the first character with a .
	dotCase := func() int {
		if chunks[0][0] != '?' {
			return put(fstKey, sndKey, 0)
		}

		nextChunks := make([]string, len(chunks))
		copy(nextChunks, chunks)
		if len(nextChunks[0]) >= 2 {
			nextChunks[0] = nextChunks[0][1:]
		} else {
			nextChunks = nextChunks[1:]
		}

		return rowToArrangements(row{chunks: nextChunks, counts: r.counts})
	}

	// dotCase tries replacing the first character with a #.
	hashCase := func() int {
		nextChunks := make([]string, len(chunks))
		copy(nextChunks, chunks)
		assert(len(nextChunks[0]) >= r.counts[0]+1, "the first chunk is too short")

		if nextChunks[0][r.counts[0]] == '#' {
			return put(fstKey, sndKey, 0)
		}

		if len(nextChunks[0]) > r.counts[0]+1 {
			nextChunks[0] = nextChunks[0][r.counts[0]+1:]
		} else {
			nextChunks = nextChunks[1:]
		}
		return rowToArrangements(row{chunks: nextChunks, counts: r.counts[1:]})
	}

	// We now put the results together.
	return put(fstKey, sndKey, dotCase()+hashCase())
}

func unfold(rows []row, count int) []row {
	out := make([]row, len(rows))
	for i := range rows {
		out[i] = unfoldRow(rows[i], count)
	}
	return out
}

func unfoldRow(r row, count int) row {
	var sb strings.Builder
	for i := 0; i < count; i++ {
		sb.WriteString(r.rawStr)
		if i < count-1 {
			sb.WriteString("?")
		}
	}
	rawStr := sb.String()

	counts := make([]int, 0, count*len(r.counts))
	for i := 0; i < count; i++ {
		counts = append(counts, r.counts...)
	}

	return row{rawStr: rawStr, counts: counts}
}

func strToInts(s string) []int {
	ss := strings.Split(s, ",")
	xs := make([]int, len(ss))
	var err error
	for i := range ss {
		xs[i], err = strconv.Atoi(ss[i])
		assert(err == nil, "unable to parse int")
	}
	return xs
}

func sum(xs []int) int {
	s := 0
	for _, x := range xs {
		s += x
	}
	return s
}

func assert(b bool, msg string, args ...any) {
	if !b {
		panic("assertion failed: " + fmt.Sprintf(msg, args...))
	}
}
