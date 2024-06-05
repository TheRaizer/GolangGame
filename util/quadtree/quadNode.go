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

// Returns the index of the quadrant in quadRect that contains all of el as well as the quadrants corresponding rect.
// Returns -1 and a nil pointer when no quadrant in quadRect contains all of el.
func QuadrantContaining(nodeBox Rect, el QuadElement) (quadrantIdx int32, quadBox *Rect) {
	if !nodeBox.Contains(&el.Rect) {
		panic("element is not contained in the given nodeBox")
	}

	width := nodeBox.W / 2
	height := nodeBox.H / 2

	// north west
	rect := Rect{X: nodeBox.X, Y: nodeBox.Y, W: width, H: height}
	if rect.Contains(&el.Rect) {
		return 0, &rect
	}

	// north east
	rect.X = width
	if rect.Contains(&el.Rect) {
		return 1, &rect
	}

	// south west
	rect.X = nodeBox.X
	rect.Y = height
	if rect.Contains(&el.Rect) {
		return 2, &rect
	}

	// south east
	rect.X = width
	if rect.Contains(&el.Rect) {
		return 3, &rect
	}

	return -1, nil
}
