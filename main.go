package main

import (
	"image"

	"github.com/TheRaizer/GolangGame/core"
	"github.com/TheRaizer/GolangGame/core/collision"
	"github.com/TheRaizer/GolangGame/core/game"
	"github.com/TheRaizer/GolangGame/core/objs"
	"github.com/TheRaizer/GolangGame/display"
	"github.com/TheRaizer/GolangGame/entities"
	"github.com/TheRaizer/GolangGame/util"
	"github.com/TheRaizer/GolangGame/util/datastructures/quadtree"
)

func main() {
	globalRect := quadtree.Rect{X: 0, Y: 0, W: display.WIDTH, H: display.HEIGHT}
	collisionSys := collision.NewCollisionSystem(globalRect)

	// generate a gray image
	img := image.NewGray(image.Rectangle{Max: image.Point{X: display.WIDTH, Y: display.HEIGHT}})
	game := game.NewGame(*img, &collisionSys)

	collider := collision.NewCollider(
		core.PLAYER_LAYER,
		"player_collider",
		quadtree.Rect{X: 0, Y: 0, W: 32, H: 32},
		&collisionSys,
		&collisionSys,
		[]func(els []quadtree.QuadElement){},
		&game,
	)
	rb := objs.NewRigidBody(core.PLAYER_LAYER, "player_rb", util.Vec2[float32]{}, &game, collider, &collisionSys, true)
	player := entities.NewPlayer("player", util.Vec2[float32]{X: 0, Y: 0}, 120, &game, &rb)

	player.AddChild(collider)
	player.AddChild(&rb)

	var wallWidth int32 = 300
	var wallHeight int32 = 32
	wall := objs.NewWall("wall_1", util.Vec2[float32]{X: 0, Y: 500}, &game, wallWidth, wallHeight)
	wall.AddChild(collision.NewCollider(
		core.WALL_LAYER,
		"wall_1_collider",
		quadtree.Rect{X: 0, Y: 0, W: wallWidth, H: wallHeight},
		&collisionSys,
		&collisionSys,
		make([]func(els []quadtree.QuadElement), 0),
		&game,
	))

	game.AddGameObject(&player)
	game.AddGameObject(&wall)
	game.Init()
}
