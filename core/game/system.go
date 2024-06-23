package game

import "github.com/TheRaizer/GolangGame/core"

type System[T core.GameObject] interface {
	RegisterObject(obj T)
	DeregisterObject(obj T)
	OnLoop()
}
