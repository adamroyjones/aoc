package main

import (
	"container/heap"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/exp/maps"
)

type graph struct {
	gryd        *grid
	labelToNode map[string]*node
	root        *node
}

func newGraph(g *grid) *graph {
	assert(g != nil, "the grid must not be nil")
	root := node{data: []pair{{i: 0, j: 0}}}
	return &graph{
		gryd:        g,
		labelToNode: map[string]*node{root.label(): &root},
		root:        &root,
	}
}

func (graf graph) dims() (int, int) { i, j := graf.gryd.dims(); return i, j }

func (graf graph) weight(n *node) int {
	assert(n != nil, "the node must not be nil")
	assert(len(n.data) > 0, "there must be at least one datum")
	lastPair := n.data[len(n.data)-1]

	assert(graf.gryd != nil, "the grid must not be nil")
	assert(len(*graf.gryd) > lastPair.i, "the grid must accommodate the pair (pair: %v, node: %s)", lastPair, n.label())
	assert(len((*graf.gryd)[lastPair.i]) > lastPair.j, "the grid must accommodate the pair (pair: %v, node: %s)", lastPair, n.label())
	return int((*graf.gryd)[lastPair.i][lastPair.j])
}

func (graf *graph) addEdge(currentNode, nextNode *node) {
	assert(currentNode != nil && nextNode != nil, "neither node may be nil")
	assert(len(nextNode.data) > 0, "there must be at least one datum")
	if currentNode.edges == nil {
		currentNode.edges = make(map[string]*node)
	}
	currentNode.edges[nextNode.label()] = nextNode
}

func (graf *graph) getNeighbours(p pair) (up, left, down, right *node) {
	if p.i > 0 {
		up = &node{data: []pair{p, {i: p.i - 1, j: p.j}}}
	}
	if p.j > 0 {
		left = &node{data: []pair{p, {i: p.i, j: p.j - 1}}}
	}
	if imax, _ := graf.dims(); p.i < imax-1 {
		down = &node{data: []pair{p, {i: p.i + 1, j: p.j}}}
	}
	if _, jmax := graf.dims(); p.j < jmax-1 {
		right = &node{data: []pair{p, {i: p.i, j: p.j + 1}}}
	}
	return
}

func (graf *graph) nextNodes(n *node) []*node {
	assert(len(n.data) > 0, "at least one datum must be present")
	assert(graf != nil, "the graph must not be nil")
	assert(graf.gryd != nil, "the grid must not be nil")

	head := n.data[len(n.data)-1]
	up, left, down, right := graf.getNeighbours(head)

	switch len(n.data) {
	case 1:
		if n.label() == graf.terminalNodeLabel() {
			return []*node{}
		}

		assert(head.i == 0 && head.j == 0, "this must be the originating node")
		_, _, down, right := graf.getNeighbours(head)
		return compact(down, right)

	case 2:
		// It's always shorter to jump to the end when you're a neighbour of it.
		if (head.i == len(*graf.gryd)-1 && head.j == len((*graf.gryd)[0])-2) || (head.i == len(*graf.gryd)-2 && head.j == len((*graf.gryd)[0])-1) {
			i, j := len(*graf.gryd)-1, len((*graf.gryd)[0])-1
			return []*node{{data: []pair{{i: i, j: j}}}}
		}

		switch d := n.data[0].direction(head); d {
		case DIRECTION_U:
			if up != nil {
				up.data = safeappend(n.data[:1], up.data...)
			}
			return compact(up, left, right)
		case DIRECTION_L:
			if left != nil {
				left.data = safeappend(n.data[:1], left.data...)
			}
			return compact(left, down, up)
		case DIRECTION_D:
			if down != nil {
				down.data = safeappend(n.data[:1], down.data...)
			}
			return compact(down, left, right)
		case DIRECTION_R:
			if right != nil {
				right.data = safeappend(n.data[:1], right.data...)
			}
			return compact(right, up, down)
		default:
			panic("unexpected direction")
		}

	case 3:
		// It's always shorter to jump to the end when you're a neighbour of it.
		if (head.i == len(*graf.gryd)-1 && head.j == len((*graf.gryd)[0])-2) || (head.i == len(*graf.gryd)-2 && head.j == len((*graf.gryd)[0])-1) {
			i, j := len(*graf.gryd)-1, len((*graf.gryd)[0])-1
			return []*node{{data: []pair{{i: i, j: j}}}}
		}

		switch d := n.data[1].direction(head); d {
		case DIRECTION_U:
			if up != nil {
				up.data = safeappend(n.data[:2], up.data...)
			}
			return compact(up, left, right)
		case DIRECTION_L:
			if left != nil {
				left.data = safeappend(n.data[:2], left.data...)
			}
			return compact(up, left, down)
		case DIRECTION_D:
			if down != nil {
				down.data = safeappend(n.data[:2], down.data...)
			}
			return compact(left, down, right)
		case DIRECTION_R:
			if right != nil {
				right.data = safeappend(n.data[:2], right.data...)
			}
			return compact(up, down, right)
		default:
			panic("unexpected direction")
		}

	case 4:
		// We have to change direction.
		switch d := n.data[1].direction(head); d {
		case DIRECTION_U:
			return compact(left, right)
		case DIRECTION_L:
			return compact(up, down)
		case DIRECTION_D:
			// It's always shorter to jump to the end when you're a neighbour of it.
			if head.i == len(*graf.gryd)-1 && head.j == len((*graf.gryd)[0])-2 {
				return []*node{{data: []pair{{i: len(*graf.gryd) - 1, j: len((*graf.gryd)[0]) - 1}}}}
			}
			return compact(left, right)
		case DIRECTION_R:
			// It's always shorter to jump to the end when you're a neighbour of it.
			if head.i == len(*graf.gryd)-2 && head.j == len((*graf.gryd)[0])-1 {
				return []*node{{data: []pair{{i: len(*graf.gryd) - 1, j: len((*graf.gryd)[0]) - 1}}}}
			}
			return compact(up, down)

		default:
			panic("unexpected direction")
		}

	default:
		assert(false, "n.data must have fewer than 3 elements (n.data = %v)", n.data)
	}

	panic("unreachable")
}

// shortestPath uses Dijkstra's algorithm with a priority queue to find the
// shortest path.
func (graf *graph) shortestPath() int {
	labelToPQI := make(map[string]*priorityQueueItem, len(graf.labelToNode))
	pq := make(priorityQueue, len(graf.labelToNode))

	// Initialise the priority queue.
	index := 0
	for label, node := range graf.labelToNode {
		pqi := priorityQueueItem{node: node, distance: math.MaxInt, index: index}
		if label == graf.root.label() {
			pqi.distance = 0
		}
		labelToPQI[label], pq[index] = &pqi, &pqi
		index++
	}
	heap.Init(&pq)

	// Dijkstra.
	for pq.Len() > 0 {
		pqi := heap.Pop(&pq).(*priorityQueueItem)
		dist := pqi.distance

		for neighbourLabel, neighbourNode := range pqi.node.edges {
			neighbourWeight := graf.weight(neighbourNode)
			neighbourPQI := labelToPQI[neighbourLabel]
			oldDistance := neighbourPQI.distance
			if newDistance := dist + neighbourWeight; newDistance < oldDistance {
				pq.update(neighbourPQI, newDistance)
			}
		}
	}

	return labelToPQI[graf.terminalNodeLabel()].distance
}

func (graf graph) terminalNodeLabel() string {
	imax, jmax := graf.dims()
	return "[(" + strconv.Itoa(imax-1) + "," + strconv.Itoa(jmax-1) + "),]"
}

// This can be passed into graphviz.
func (graf graph) String() string {
	var sb strings.Builder
	sb.WriteString("digraph {\n")
	srcLabels := maps.Keys(graf.labelToNode)
	slices.Sort(srcLabels)
	for _, srcLabel := range srcLabels {
		n := graf.labelToNode[srcLabel]
		dstLabels := maps.Keys(n.edges)
		slices.Sort(dstLabels)
		for _, dstLabel := range dstLabels {
			sb.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\"\n", srcLabel, dstLabel))
		}
	}
	sb.WriteString("}\n")
	return sb.String()
}
