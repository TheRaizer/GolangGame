package quadtree

import (
	"github.com/TheRaizer/GolangGame/util"
	"github.com/veandco/go-sdl2/sdl"
)

type Rect sdl.Rect

func (rect *Rect) Center() util.Vec2[int32] {
	centerX := rect.X + rect.W/2
	centerY := rect.Y + rect.H/2

	return util.Vec2[int32]{X: centerX, Y: centerY}
}

func (rect *Rect) Contains(otherRect *Rect) bool {
	return rect.X <= otherRect.X && rect.Y <= otherRect.Y && rect.Right() >= otherRect.Right() && rect.Bottom() >= otherRect.Bottom()
}

func (rect *Rect) Intersects(otherRect *Rect) bool {
	// does not over-lap if:
	// the other rect's ride side is to the left of the rect
	// the other rect's top is below the rects bottom
	// the other rect's left side is to the right of the rect
	// the other rect's bottom is above the rects top

	// we take the negation
	return !(otherRect.Right() < rect.X || otherRect.Y > rect.Bottom() || otherRect.X > rect.Right() || otherRect.Bottom() < rect.Y)
}

func (rect *Rect) Right() int32 {
	return rect.X + rect.W
}

func (rect *Rect) Bottom() int32 {
	return rect.Y + rect.H
}
