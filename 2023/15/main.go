package main

import (
	"os"
	"slices"
	"strconv"
	"strings"
)

type box struct {
	m     map[string]int
	order []string
}

func partOne(filename string) int {
	bs, err := os.ReadFile(filename)
	assert(err == nil)
	strs := strings.Split(strings.TrimSpace(string(bs)), ",")
	sum := 0
	for _, str := range strs {
		sum += hash(str)
	}
	return sum
}

func partTwo(filename string) int {
	bs, err := os.ReadFile(filename)
	assert(err == nil)

	boxes := make([]*box, 256)
	for i := range boxes {
		boxes[i] = newBox()
	}

	strs := strings.Split(strings.TrimSpace(string(bs)), ",")
	for _, str := range strs {
		label, valueStr, ok := strings.Cut(str, "=")
		if ok {
			value, err := strconv.Atoi(valueStr)
			assert(err == nil)
			boxes[hash(label)].set(label, value)
			continue
		}

		label, after, ok := strings.Cut(str, "-")
		assert(ok && after == "")
		boxes[hash(label)].del(label)
	}

	return power(boxes)
}

func newBox() *box {
	m := make(map[string]int)
	order := []string{}
	return &box{m: m, order: order}
}

func (b *box) set(label string, value int) {
	_, ok := b.m[label]
	b.m[label] = value
	if !ok {
		b.order = append(b.order, label)
	}
}

func (b *box) del(label string) {
	if _, ok := b.m[label]; !ok {
		return
	}
	delete(b.m, label)
	b.order = slices.DeleteFunc(b.order, func(s string) bool { return s == label })
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

func power(boxes []*box) int {
	result := 0
	for i, b := range boxes {
		for j, label := range b.order {
			result += (i + 1) * (j + 1) * b.m[label]
		}
	}
	return result
}

func assert(b bool) {
	if !b {
		panic("assertion failed")
	}
}
