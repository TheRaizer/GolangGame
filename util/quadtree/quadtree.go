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

	globalBox Rect
	root      *QuadNode
}

// Inserts an element into the quadtree
func (quadtree *CollisionQuadTree) Insert(el QuadElement) {
	quadtree.insert(quadtree.root, quadtree.globalBox, 0, el)
}

// Queries for elements that lie inside the given rect
func (quadtree *CollisionQuadTree) Query(rect Rect) {

}

// Removes an element from the quad tree
func (quadtree *CollisionQuadTree) Remove(el QuadElement) {

}

func (quadtree *CollisionQuadTree) insert(node *QuadNode, nodeBox Rect, depth int, el QuadElement) {
	if node == nil {
		panic("node pointer was nil")
	}

	if !nodeBox.Contains(&el.Rect) {
		panic("the given quad does not contain the element rect")
	}

	if node.isLeaf() {
		// if were at max depth we want to insert the value to avoid infinite recursion
		if depth >= quadtree.maxDepth || len(node.els) < quadtree.threshold {
			node.els = append(node.els, el)
		} else {
			node.split(nodeBox)

			// now the node is no longer a leaf
			quadtree.insert(node, nodeBox, depth, el)
		}
	} else {
		// look for the child quad that contains the element
		quadrant, quadBox := QuadrantContaining(nodeBox, el)

		// if none contain it then give it to the current quad
		if quadrant == -1 {
			node.els = append(node.els, el)
		} else {
			// otherwise we may be able to partition even further to find the best fitting quad node
			quadtree.insert(node.children[quadrant], *quadBox, depth+1, el)
		}
	}
}
