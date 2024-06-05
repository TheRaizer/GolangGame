package quadtree

type QuadNode struct {
	// NOTE: according to https://stackoverflow.com/questions/41946007/efficient-and-well-explained-implementation-of-a-quadtree-for-2d-collision-det
	// it would be best not to store small lists in each node, but imma ignore that for the initial implementation

	children [4]*QuadNode
	els      []QuadElement
}

func (quadNode *QuadNode) isLeaf() bool {
	return quadNode.children[0] == nil
}

// Turns a leaf node into a branch node by splitting into 4 different children leaf nodes.
// Each element of the branch node will be distributed amongst the children if the corresponding child
// quadrant can contain the elements rect. Otherwise the element will return back to the branch node
func (quadNode *QuadNode) split(quadRect Rect) {
	if !quadNode.isLeaf() {
		panic("only a leaf node can be split")
	}

	var newEls []QuadElement

	for i := 0; i < 4; i++ {
		quadNode.children[i] = &QuadNode{}
	}

	for _, el := range quadNode.els {
		quadrant, _ := QuadrantContaining(quadRect, el)
		if quadrant == -1 {
			newEls = append(newEls, el)
		} else {
			quadNode.children[quadrant].els = append(quadNode.children[quadrant].els, el)
		}
	}

}

// computes the rect of a quadrant given its parent quad along with the specific quadrant idx
// 0 = NW
// 1 = NE
// 2 = SW
// 3 = SE
// otherwise nil
func ComputeQuadRect(parentRect Rect, quadrantIdx int) *Rect {
	width := parentRect.W / 2
	height := parentRect.H / 2

	switch quadrantIdx {
	case 0:
		return &Rect{X: parentRect.X, Y: parentRect.Y, W: width, H: height}
	case 1:
		return &Rect{X: parentRect.X + width, Y: parentRect.Y, W: width, H: height}
	case 2:
		return &Rect{X: parentRect.X, Y: parentRect.Y + height, W: width, H: height}
	case 3:
		return &Rect{X: parentRect.X + width, Y: parentRect.Y + height, W: width, H: height}
	default:
		return nil

	}
}

// Returns the index of the quadrant in quadRect that contains all of el as well as the quadrants corresponding rect.
// Returns -1 and a nil pointer when no quadrant in quadRect contains all of el.
func QuadrantContaining(nodeRect Rect, el QuadElement) (quadrantIdx int32, quadRect *Rect) {
	if !nodeRect.Contains(el.Rect) {
		panic("element is not contained in the given nodeBox")
	}

	rect := ComputeQuadRect(nodeRect, 0)
	if rect.Contains(el.Rect) {
		return 0, rect
	}

	rect = ComputeQuadRect(nodeRect, 1)
	if rect.Contains(el.Rect) {
		return 1, rect
	}

	rect = ComputeQuadRect(nodeRect, 2)
	if rect.Contains(el.Rect) {
		return 2, rect
	}

	rect = ComputeQuadRect(nodeRect, 3)
	if rect.Contains(el.Rect) {
		return 3, rect
	}

	return -1, nil
}
