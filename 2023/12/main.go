package main

import (
	"os"
	"slices"
	"strconv"
	"strings"
)

// row is the parsed form of the lines. A line is a line of input, for example
// "???.### 1,1,3".
type row struct {
	// rawStr corresponds to "???.###" in the example above.
	rawStr string
	// counts corresponds to [1, 2, 3] in the example above.
	counts []int
	// chunks is used internally. It's the result of cutting rawStr at . and
	// removing empty elements and so, with the example above, ["???", "###"] is
	// what's produced. This is done because dots are irrelevant in the analysis
	// except insofar as they serve as barriers.
	chunks []string
}

func partOne(filename string) int {
	rows := filenameToRows(filename)
	return rowsToArrangements(rows)
}

func partTwo(filename string) int {
	rows := filenameToRows(filename)
	return rowsToArrangements(unfold(rows, 5))
}

func filenameToRows(filename string) []row {
	bs, err := os.ReadFile(filename)
	assert(err == nil)
	lines := strings.Split(strings.TrimSpace(string(bs)), "\n")
	rows := make([]row, len(lines))
	for i, line := range lines {
		rawStr, countsStr, ok := strings.Cut(line, " ")
		assert(ok)

		ss := strings.Split(countsStr, ",")
		counts := make([]int, len(ss))
		var err error
		for i := range ss {
			counts[i], err = strconv.Atoi(ss[i])
			assert(err == nil)
		}
		rows[i] = row{rawStr: rawStr, counts: counts}
	}
	return rows
}

func rowsToArrangements(rows []row) int {
	arrs := 0
	for _, row := range rows {
		fn := func(x, y byte) bool { return x == '.' && x == y }
		row.rawStr = string(slices.CompactFunc([]byte(row.rawStr), fn))
		arrs += rowToArrangements(row)
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
		fn := func(s string) bool { return strings.ContainsRune(s, '#') }
		if slices.ContainsFunc(chunks, fn) {
			return put(fstKey, sndKey, 0)
		}
		return put(fstKey, sndKey, 1)
	}

	// The chunks slice must have more available sites than the sum of
	// counts. Further, we must have at least one site to fill in. We now need to
	// do some work. We'll act recursively on the first chunk.
	//
	// Must we ignore the first chunk entirely?
	if chunkLengths[0] < r.counts[0] {
		if strings.ContainsRune(chunks[0], '#') {
			return put(fstKey, sndKey, 0)
		}
		val := rowToArrangements(row{chunks: chunks[1:], counts: r.counts})
		return put(fstKey, sndKey, val)
	}

	// Could we ignore the first chunk entirely?
	if chunkLengths[0] == r.counts[0] {
		useFirstChunk := rowToArrangements(row{chunks: chunks[1:], counts: r.counts[1:]})
		omitFirstChunk := 0
		if !strings.ContainsRune(chunks[0], '#') {
			omitFirstChunk = rowToArrangements(row{chunks: chunks[1:], counts: r.counts})
		}

		return put(fstKey, sndKey, useFirstChunk+omitFirstChunk)
	}

	// The first chunk is longer than r.counts[0]. We'll have to consider the
	// possibility of repeated use of the first chunk. We'll now work recursively
	// from the first element of the first chunk.
	//
	// dotCase tries replacing the first character with a ".".
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

	// hashCase tries replacing the first character with a "#".
	hashCase := func() int {
		nextChunks := make([]string, len(chunks))
		copy(nextChunks, chunks)
		assert(len(nextChunks[0]) >= r.counts[0]+1)

		// If we start with # and need r.counts[0] # in the first filled-in chunk,
		// then this condition makes this impossible.
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
		r := rows[i]
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

		out[i] = row{rawStr: rawStr, counts: counts}
	}
	return out
}

func sum(xs []int) int {
	s := 0
	for _, x := range xs {
		s += x
	}
	return s
}

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}
