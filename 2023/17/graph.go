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

	minSegmentLength int
	maxSegmentLength int
}

func newGraph(g *grid, minSegmentLength, maxSegmentLength int) *graph {
	assert(g != nil, "the grid must not be nil")
	assert(maxSegmentLength > minSegmentLength && minSegmentLength > 0, "the segment lengths are inconsistent")
	root := node{segment: []pair{{i: 0, j: 0}}}
	return &graph{
		gryd:             g,
		labelToNode:      map[string]*node{root.label(): &root},
		root:             &root,
		minSegmentLength: minSegmentLength,
		maxSegmentLength: maxSegmentLength,
	}
}

func (graf graph) dims() (int, int) { i, j := graf.gryd.dims(); return i, j }

func (graf graph) weight(n *node) int {
	assert(n != nil, "the node must not be nil")
	assert(len(n.segment) > 0, "there must be at least one datum")
	assert(graf.gryd != nil, "the grid must not be nil")
	assert(len(*graf.gryd) > n.head().i, "the grid must accommodate the pair (pair: %v, node: %s)", n.head(), n.label())
	assert(len((*graf.gryd)[n.head().i]) > n.head().j, "the grid must accommodate the pair (pair: %v, node: %s)", n.head(), n.label())
	return int((*graf.gryd)[n.head().i][n.head().j])
}

func (graf *graph) addEdge(currentNode, nextNode *node) {
	assert(currentNode != nil && nextNode != nil, "neither node may be nil")
	assert(len(nextNode.segment) > 0, "there must be at least one datum")
	if currentNode.edges == nil {
		currentNode.edges = make(map[string]*node)
	}
	currentNode.edges[nextNode.label()] = nextNode
}

func (graf *graph) neighbours(n *node) []*node {
	assert(len(n.segment) > 0, "at least one datum must be present")
	assert(graf != nil, "the graph must not be nil")
	assert(graf.gryd != nil, "the grid must not be nil")

	if graf.isTerminal(n) {
		return []*node{}
	}

	return compact(
		graf.neighbour(n, DIRECTION_U),
		graf.neighbour(n, DIRECTION_L),
		graf.neighbour(n, DIRECTION_D),
		graf.neighbour(n, DIRECTION_R),
	)
}

func (graf *graph) neighbour(n *node, dir direction) *node {
	assert(n != nil, "the node must be non-nil")

	var loc pair
	switch dir {
	case DIRECTION_U:
		loc = pair{i: n.head().i - 1, j: n.head().j}
	case DIRECTION_L:
		loc = pair{i: n.head().i, j: n.head().j - 1}
	case DIRECTION_D:
		loc = pair{i: n.head().i + 1, j: n.head().j}
	case DIRECTION_R:
		loc = pair{i: n.head().i, j: n.head().j + 1}
	default:
		panic("unexpected direction")
	}
	candidate := &node{segment: safeappend(n.segment, loc)}

	// Are we in the grid?
	imax, jmax := graf.dims()
	if candidate.head().i < 0 || candidate.head().i >= imax {
		return nil
	}
	if candidate.head().j < 0 || candidate.head().j >= jmax {
		return nil
	}

	// Are we beginning a segment?
	if len(candidate.segment) == 2 {
		return candidate
	}

	// Given the above, we've a non-trivial segment. We need to examine the
	// direction of travel.
	//
	// Are we travelling in the opposite direction?
	neck := candidate.segment[len(candidate.segment)-2]
	if candidate.tail().direction(neck).opposite() == neck.direction(candidate.head()) {
		return nil
	}

	// Are we continuing in the same direction?
	if candidate.tail().direction(neck) == neck.direction(candidate.head()) {
		if len(candidate.segment) > graf.maxSegmentLength {
			return nil
		}

		// Are we at the end of the path...
		if candidate.head().i == imax-1 && candidate.head().j == jmax-1 {
			// ...and within the allowed limit?
			if len(candidate.segment) < graf.minSegmentLength {
				return nil
			}

			candidate.segment = []pair{candidate.head()}
			return candidate
		}

		return candidate
	}

	// We've turned—but are we allowed to? Note that we need -1 here to exclude
	// the head.
	if len(candidate.segment)-1 < graf.minSegmentLength {
		return nil
	}

	// Are we at the end of the path? If so, we can only include this if the
	// minimum segment length is 1, as we've just turned.
	if candidate.head().i == imax-1 && candidate.head().j == jmax-1 && graf.minSegmentLength == 1 {
		candidate.segment = []pair{candidate.head()}
		return candidate
	}

	// Finally: we've turned, we're allowed to, we're not at the end, but now we
	// need to now fix up the history to erase the section of the segment befoer
	// the turn.
	candidate.segment = candidate.segment[len(candidate.segment)-2:]
	return candidate
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

func (graf graph) isTerminal(n *node) bool {
	assert(n != nil, "the node must not be nil")
	imax, jmax := graf.dims()
	return len(n.segment) == 1 && n.head().i == imax-1 && n.head().j == jmax-1
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
