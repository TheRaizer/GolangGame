package core

import "github.com/TheRaizer/GolangGame/core/objs"

type System[T objs.GameObject] interface {
	RegisterObject(obj T)
	DeregisterObject(obj T)
	OnLoop()
}
