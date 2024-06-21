package main

import (
	"fmt"
	"image"

	"github.com/TheRaizer/GolangGame/core"
	"github.com/TheRaizer/GolangGame/core/collision"
	"github.com/TheRaizer/GolangGame/core/objs"
	"github.com/TheRaizer/GolangGame/display"
	"github.com/TheRaizer/GolangGame/entities"
	"github.com/TheRaizer/GolangGame/util"
	datastructures "github.com/TheRaizer/GolangGame/util/datastructures/quadtree"
)

func main() {
	globalRect := datastructures.Rect{X: 0, Y: 0, W: display.WIDTH, H: display.HEIGHT}
	collisionSys := collision.NewCollisionSystem(globalRect)

	// generate a gray image
	img := image.NewGray(image.Rectangle{Max: image.Point{X: display.WIDTH, Y: display.HEIGHT}})
	game := core.NewGame(*img, &collisionSys)

	player := entities.NewPlayer("player", util.Vec2[float32]{X: 0, Y: 0})
	player.AddChild(collision.NewCollider(
		"player_collider",
		datastructures.Rect{X: 0, Y: 0, W: 32, H: 32},
		&collisionSys,
		[]func(els []datastructures.QuadElement){
			func(els []datastructures.QuadElement) { fmt.Println(els) },
		},
	))

	wall := objs.NewWall("wall_1", util.Vec2[float32]{X: 50, Y: 50})
	wall.AddChild(collision.NewCollider(
		"wall_1_collider",
		datastructures.Rect{X: 0, Y: 0, W: 32, H: 32},
		&collisionSys,
		make([]func(els []datastructures.QuadElement), 0),
	))

	game.AddGameObject(&player)
	game.AddGameObject(&wall)
	game.Init()
}
