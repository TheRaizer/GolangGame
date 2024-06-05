package quadtree

// https://pvigier.github.io/2019/08/04/quadtree-collision-detection.html

// NOTE: in this implementation every node will contain only the elements that fit within their bounds.
// This is compared with an implementation where only leaf nodes contain their contents
// This is because hit boxes can vary greatly in size so a large hitbox that is contained in multiple
// leaf nodes would require it to be stored multiple times.

type QuadTree interface {
	Insert(el QuadElement)
	Remove(el QuadElement)
}

type BaseQuadTree struct {
	nodes []QuadNode
}

type CollisionQuadTree struct {
	*BaseQuadTree
	threshold int // max number of elements before we split the quad
	maxDepth  int // max number of times we will allow quads to be split

	globalRect Rect
	root       *QuadNode
}

// Inserts an element into the quadtree
func (quadtree *CollisionQuadTree) Insert(el QuadElement) {
	quadtree.insert(quadtree.root, quadtree.globalRect, 0, el)
}

// Queries for elements that lie inside the given rect
func (quadtree *CollisionQuadTree) Query(hitbox Rect) []QuadElement {
	return quadtree.query(quadtree.root, quadtree.globalRect, hitbox)
}

// Removes an element from the quad tree
func (quadtree *CollisionQuadTree) Remove(el QuadElement) {
	quadtree.remove(quadtree.root, quadtree.globalRect, el)
}

func (quadtree *CollisionQuadTree) query(node *QuadNode, nodeRect Rect, hitbox Rect) []QuadElement {
	var intersectingEls []QuadElement

	if node == nil {
		panic("node pointer to query was nill")
	}

	if !nodeRect.Intersects(hitbox) {
		return intersectingEls
	}

	if node.isLeaf() {
		for _, el := range node.els {
			if hitbox.Intersects(el.Rect) {
				intersectingEls = append(intersectingEls, el)
			}
		}
	} else {
		// call query on only the quadrants that intersect with the hitbox
		// children cannot be nil since the node is not a leaf
		for i, quadNode := range node.children {
			quadRect := ComputeQuadRect(nodeRect, i)
			if hitbox.Intersects(*quadRect) {
				intersectingQuadEls := quadtree.query(quadNode, *quadRect, hitbox)
				intersectingEls = append(intersectingEls, intersectingQuadEls...)
			}

		}

		return intersectingEls
	}

	return intersectingEls
}

func (quadtree *CollisionQuadTree) remove(node *QuadNode, nodeRect Rect, el QuadElement) bool {
	if node == nil {
		panic("node pointer was nil")
	}

	if !nodeRect.Contains(el.Rect) {
		panic("the given quad does not contain the element rect")
	}

	if node.isLeaf() {
		removeValue(node, el)
		// we have removed a value from a leaf node, so we may be able to merge in to the parent
		return true
	} else {
		quadrantIdx, quadRect := QuadrantContaining(nodeRect, el)

		if quadrantIdx == -1 {
			removeValue(node, el)
			return false
		} else {
			// if we end up removing a value from the child node
			if quadtree.remove(node.children[quadrantIdx], *quadRect, el) {

				// return whether we should try and merge again
				return quadtree.tryMerge(node)
			} else {
				return false
			}
		}
	}
}

// Attempts to merge child quads into parent.
func (quadtree *CollisionQuadTree) tryMerge(node *QuadNode) bool {
	if node == nil {
		panic("when merging the node given was a null pointer")
	}

	// waves the possiblity that the children are nil pointers
	if node.isLeaf() {
		panic("only interior nodes can be merged")
	}

	totalEls := len(node.els)
	for _, child := range node.children {
		if !child.isLeaf() {
			return false
		}

		totalEls += len(child.els)
	}

	if totalEls <= quadtree.threshold {
		for i, child := range node.children {
			for _, childEl := range child.els {
				node.els = append(node.els, childEl)
			}

			node.children[i] = nil
		}
		return true
	} else {
		return false
	}
}

func removeValue(node *QuadNode, el QuadElement) {
	for i, otherEl := range node.els {
		if el.Id == otherEl.Id {
			node.els[i] = node.els[len(node.els)-1]
			node.els = node.els[:len(node.els)-1]
			return
		}
	}

	panic("unable to find the given element with id: " + el.Id)
}

func (quadtree *CollisionQuadTree) insert(node *QuadNode, nodeRect Rect, depth int, el QuadElement) {
	if node == nil {
		panic("node pointer was nil")
	}

	if !nodeRect.Contains(el.Rect) {
		panic("the given quad does not contain the element rect")
	}

	if node.isLeaf() {
		// if were at max depth we want to insert the value to avoid infinite recursion
		if depth >= quadtree.maxDepth || len(node.els) < quadtree.threshold {
			node.els = append(node.els, el)
		} else {
			node.split(nodeRect)

			// now the node is no longer a leaf
			quadtree.insert(node, nodeRect, depth, el)
		}
	} else {
		// look for the child quad that contains the element
		quadrant, quadRect := QuadrantContaining(nodeRect, el)

		// if none contain it then give it to the current quad
		if quadrant == -1 {
			node.els = append(node.els, el)
		} else {
			// otherwise we may be able to partition even further to find the best fitting quad node
			quadtree.insert(node.children[quadrant], *quadRect, depth+1, el)
		}
	}
}
