package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type box struct {
	m     map[string]int
	order []string
}

func (b *box) set(label string, value int) {
	_, ok := b.m[label]
	b.m[label] = value
	if !ok {
		b.order = append(b.order, label)
	}
}

func (b *box) rm(label string) {
	_, ok := b.m[label]
	if !ok {
		return
	}
	delete(b.m, label)
	b.order = slices.DeleteFunc(b.order, func(s string) bool { return s == label })
}

func newBox() *box {
	m := make(map[string]int)
	order := []string{}
	return &box{m: m, order: order}
}

func partTwo(filename string) int {
	bs, err := os.ReadFile(filename)
	assert(err == nil, "failed to read the file")

	boxes := make([]*box, 256)
	for i := range boxes {
		boxes[i] = newBox()
	}

	strs := strings.Split(strings.TrimSpace(string(bs)), ",")
	for _, str := range strs {
		label, valueStr, ok := strings.Cut(str, "=")
		if ok {
			value, err := strconv.Atoi(valueStr)
			assert(err == nil, "failed to convert a string to an int")
			idx := hash(label)
			boxes[idx].set(label, value)
			continue
		}

		label, after, ok := strings.Cut(str, "-")
		assert(ok, "expected = or -")
		assert(after == "", "expected no value")
		idx := hash(label)
		boxes[idx].rm(label)
	}

	return power(boxes)
}

func power(boxes []*box) int {
	result := 0
	for i, b := range boxes {
		for j, label := range b.order {
			result += (i + 1) * (j + 1) * b.m[label]
		}
	}
	return result
}

func partOne(filename string) int {
	bs, err := os.ReadFile(filename)
	assert(err == nil, "failed to read the file")
	strs := strings.Split(strings.TrimSpace(string(bs)), ",")
	sum := 0
	for _, str := range strs {
		sum += hash(str)
	}
	return sum
}

func hash(s string) int {
	val := 0
	for _, ru := range s {
		val += int(ru)
		val *= 17
		val %= 256
	}
	return val
}

func assert(b bool, msg string, args ...any) {
	if !b {
		panic("assertion failed: " + fmt.Sprintf(msg, args...))
	}
}
