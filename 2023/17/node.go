package main

type node struct {
	data  pairs
	edges map[string]*node
}

func (n *node) label() string {
	assert(n != nil, "the node must not be nil")
	assert(len(n.data) > 0, "expected a pair")
	return n.data.label()
}
