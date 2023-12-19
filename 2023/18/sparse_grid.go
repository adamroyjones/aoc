package main

type sparseGrid struct{ edges []sparseGridEdge }

type sparseGridEdge struct {
	start, end       pair
	startTyp, endTyp pipeType
}

func newSparseGrid(insts instructions) *sparseGrid {
	assert(len(insts) > 0)
	sg := &sparseGrid{}
	for _, inst := range insts {
		sg.step(inst)
	}
	sg.fixOrigin()
	return sg
}

func (sg *sparseGrid) step(inst instruction) {
	nextEdge := sparseGridEdge{start: pair{i: 0, j: 0}}

	var prevEdge *sparseGridEdge
	if len(sg.edges) > 0 {
		prevEdge = &sg.edges[len(sg.edges)-1]
		nextEdge.start = prevEdge.end
	}

	var end pair
	assert(inst.dist > 0)
	switch inst.dir {
	case DIRECTION_U:
		end = pair{i: nextEdge.start.i - inst.dist, j: nextEdge.start.j}
	case DIRECTION_L:
		end = pair{i: nextEdge.start.i, j: nextEdge.start.j - inst.dist}
	case DIRECTION_D:
		end = pair{i: nextEdge.start.i + inst.dist, j: nextEdge.start.j}
	case DIRECTION_R:
		end = pair{i: nextEdge.start.i, j: nextEdge.start.j + inst.dist}
	default:
		panic("unexpected direction")
	}
	nextEdge.end = end

	if prevEdge != nil {
		pt := prevEdge.direction().pipeType(nextEdge.direction())
		prevEdge.endTyp, nextEdge.startTyp = pt, pt
	}

	sg.edges = append(sg.edges, nextEdge)
}

func (sg *sparseGrid) fixOrigin() {
	assert(sg != nil && len(sg.edges) > 0)
	first, last := &sg.edges[0], &sg.edges[len(sg.edges)-1]
	pt := last.direction().pipeType(first.direction())
	first.startTyp, last.endTyp = pt, pt
}

func (sg *sparseGrid) area(quads *[]quad) int {
	assert(sg != nil && quads != nil)
	area := 0
	for _, quad := range *quads {
		// By construction, if any part of q is in the region, then all of q is.
		if sg.contains(pair{i: quad.irange.min, j: quad.jrange.min}) {
			area += quad.area()
		}
	}
	return area
}

// contains says whether the point v lies within the bounded region. This uses
// the same kind of crossing argument used in part 1.
func (sg *sparseGrid) contains(p pair) bool {
	outside := true
	for _, sge := range sg.edges {
		// If v is above, below, or to the left of the edge, then the edge is irrelevant.
		if imin := min(sge.start.i, sge.end.i); p.i < imin {
			continue
		}
		if imax := max(sge.start.i, sge.end.i); p.i > imax {
			continue
		}
		if jmin := min(sge.start.j, sge.end.j); p.j < jmin {
			continue
		}

		// We may be on the edge itself, in which case we can bail out early.
		if sge.contains(p) {
			return true
		}

		if sge.cross(p) {
			outside = !outside
		}
	}

	return !outside
}

func (sge sparseGridEdge) direction() direction {
	assert(sge.start.i != sge.end.i || sge.start.j != sge.end.j)
	assert(sge.start.i == sge.end.i || sge.start.j == sge.end.j)
	if sge.start.i < sge.end.i {
		return DIRECTION_D
	}
	if sge.start.i > sge.end.i {
		return DIRECTION_U
	}
	if sge.start.j < sge.end.j {
		return DIRECTION_R
	}
	return DIRECTION_L
}

func (sge sparseGridEdge) cross(p pair) bool {
	// The edge is horizontal.
	if sge.start.i == sge.end.i {
		assert(p.i == sge.start.i)

		// That is, if the start point is the first point encountered when moving right.
		if sge.start.j < sge.end.j {
			switch sge.startTyp {
			case PIPE_BEND_DL:
				return sge.endTyp == PIPE_BEND_UR
			case PIPE_BEND_UL:
				return sge.endTyp == PIPE_BEND_DR
			default:
				panic("invalid start type")
			}
		}

		switch sge.endTyp {
		case PIPE_BEND_DL:
			return sge.startTyp == PIPE_BEND_UR
		case PIPE_BEND_UL:
			return sge.startTyp == PIPE_BEND_DR
		default:
			panic("invalid end type")
		}
	}

	// The edge is vertical...
	assert(sge.start.j == sge.end.j)
	assert(p.j > sge.start.j)
	// ...and we may have crossed it by passing through the relative interior of the edge.
	if min(sge.start.i, sge.end.i) < p.i && p.i < max(sge.start.i, sge.end.i) {
		return true
	}

	// The only remaining cases are the endpoints of a vertical edge. However,
	// every endpoint of a vertical edge is the endpoint of a corresponding
	// horizontal edge, which we've considered above. We return false to avoid
	// counting this edge.
	return false
}

func (sge sparseGridEdge) contains(p pair) bool {
	// sge is a vertical edge.
	if sge.start.i == sge.end.i {
		return p.i == sge.start.i && min(sge.start.j, sge.end.j) <= p.j && p.j <= max(sge.start.j, sge.end.j)
	}

	// As sge is not a vertical edge, it must be a horizontal edge.
	assert(sge.start.j == sge.end.j)
	return p.j == sge.start.j && min(sge.start.i, sge.end.i) <= p.i && p.i <= max(sge.start.i, sge.end.i)
}
