package main

type node struct {
	// segment is the path segment of the walker. The current site of the walker
	// is the head of the path segment, which is the last element in this slice.
	segment pairs
	edges   map[string]*node
}

func (n *node) tail() pair {
	assert(len(n.segment) > 0, "there must be at least one datum")
	return n.segment[0]
}

func (n *node) head() pair {
	assert(len(n.segment) > 0, "there must be at least one datum")
	return n.segment[len(n.segment)-1]
}

func (n *node) direction() direction {
	assert(len(n.segment) > 1, "there must be at least two data")
	return n.tail().direction(n.head())
}

func (n *node) label() string {
	assert(n != nil, "the node must not be nil")
	assert(len(n.segment) > 0, "there must be at least one datum")
	return n.segment.label()
}
