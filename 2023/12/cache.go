package main

import (
	"strconv"
	"strings"
)

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
