package quadtree

// https://pvigier.github.io/2019/08/04/quadtree-collision-detection.html

// NOTE: in this implementation every node will contain only the elements that fit within their bounds.
// This is compared with an implementation where only leaf nodes contain their contents
// This is because hit boxes can vary greatly in size so a large hitbox that is contained in multiple
// leaf nodes would require it to be stored multiple times.

type BaseQuadTree struct {
	nodes []QuadNode
}

type QuadTree struct {
	*BaseQuadTree
	threshold int // max number of elements before we split the quad
	maxDepth  int // max number of times we will allow quads to be split

	globalRect Rect
	root       *QuadNode
}

// Inserts an element into the quadtree
func (quadtree *QuadTree) Insert(el QuadElement) {
	quadtree.insert(quadtree.root, quadtree.globalRect, 0, el)
}

// Queries for elements that lie inside the given rect
func (quadtree *QuadTree) Query(hitbox Rect) []QuadElement {
	return quadtree.query(quadtree.root, quadtree.globalRect, hitbox)
}

// Removes an element from the quad tree
func (quadtree *QuadTree) Remove(el QuadElement) {
	quadtree.remove(quadtree.root, quadtree.globalRect, el)
}

func (quadtree *QuadTree) query(node *QuadNode, nodeRect Rect, hitbox Rect) []QuadElement {
	var intersectingEls []QuadElement

	if node == nil {
		panic("node was nil")
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

func (quadtree *QuadTree) remove(node *QuadNode, nodeRect Rect, el QuadElement) bool {
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
func (quadtree *QuadTree) tryMerge(node *QuadNode) bool {
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

func (quadtree *QuadTree) insert(node *QuadNode, nodeRect Rect, depth int, el QuadElement) {
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

// // NOTE: untested and inefficient
// func (quadtree *QuadTree) FindAllIntersections() [][2]QuadElement {
// 	return quadtree.findAllIntersections(quadtree.root, quadtree.globalRect)
// }
//
// // return every pairwise interscetion that occurs in the tree
// func (quadtree *QuadTree) findAllIntersections(node *QuadNode, nodeRect Rect) [][2]QuadElement {
// 	var intersections [][2]QuadElement
//
// 	if node == nil {
// 		panic("node pointer was nil")
// 	}
//
// 	// interate through the given node and compare each value to values within itself
// 	// avoid rechecking intersections
// 	for i := 0; i < len(node.els); i++ {
// 		for j := 0; j < i; j++ {
// 			if node.els[i].Rect.Intersects(node.els[j].Rect) {
// 				intersections = append(intersections, [2]QuadElement{node.els[i], node.els[j]})
// 			}
// 		}
// 	}
//
// 	// if its not a leaf check its descendants for more intersections
// 	if !node.isLeaf() {
// 		for _, el := range node.els {
// 			// for each element check intersections between itself and its descendents
// 			intersectionsInDescendants := quadtree.findAllIntersectionsInDescendants(node, nodeRect, el)
// 			intersections = append(intersections, intersectionsInDescendants...)
// 		}
//
// 		for i, child := range node.children {
// 			// find all the intersections in any values in the children now recursively
// 			childIntersections := quadtree.findAllIntersections(child, *ComputeQuadRect(nodeRect, i))
// 			intersections = append(intersections, childIntersections...)
// 		}
// 	}
//
// 	return intersections
// }
//
// // find all intersections between an element stored in a branch node, and the elements of all of its descendants
// func (quadtree *QuadTree) findAllIntersectionsInDescendants(parentNode *QuadNode, parentRect Rect, parentEl QuadElement) [][2]QuadElement {
// 	var intersections [][2]QuadElement
//
// 	if parentNode == nil {
// 		panic("node was nil")
// 	}
//
// 	if parentNode.isLeaf() {
// 		panic("node cannot be a leaf")
// 	}
//
// 	for i, child := range parentNode.children {
// 		// find every intersection with the parent element and any child elements
// 		for _, el := range child.els {
// 			if el.Rect.Intersects(parentEl.Rect) {
// 				intersections = append(intersections, [2]QuadElement{parentEl, el})
// 			}
// 		}
//
// 		if !child.isLeaf() {
// 			intersectionsInChild := quadtree.findAllIntersectionsInDescendants(child, *ComputeQuadRect(parentRect, i), parentEl)
// 			intersections = append(intersections, intersectionsInChild...)
// 		}
// 	}
//
// 	return intersections
// }
