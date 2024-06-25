package core

type System[T GameObject] interface {
	RegisterObject(obj T)
	DeregisterObject(obj T)
	OnLoop()
}
