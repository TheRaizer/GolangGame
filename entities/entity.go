package entities

import (
	"github.com/TheRaizer/GolangGame/display"
)

type Entity interface {
	canMove(p display.Pixel) bool
	move(dt uint64)
}
