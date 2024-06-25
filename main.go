package main

import (
	"fmt"
	"image"

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

	player := entities.NewPlayer("player", util.Vec2[float32]{X: 0, Y: 0}, &game)
	player.AddChild(collision.NewCollider(
		"player_collider",
		quadtree.Rect{X: 0, Y: 0, W: 32, H: 32},
		&collisionSys,
		&collisionSys,
		[]func(els []quadtree.QuadElement){
			func(els []quadtree.QuadElement) {
				for _, el := range els {
					if el.Id == "wall_1_collider" {
						obj := game.GetGameObject(el.Id)
						fmt.Println(obj)
					}
				}
			},
		},
		&game,
	))

	wall := objs.NewWall("wall_1", util.Vec2[float32]{X: 50, Y: 50}, &game)
	wall.AddChild(collision.NewCollider(
		"wall_1_collider",
		quadtree.Rect{X: 0, Y: 0, W: 32, H: 32},
		&collisionSys,
		&collisionSys,
		make([]func(els []quadtree.QuadElement), 0),
		&game,
	))

	game.AddGameObject(&player)
	game.AddGameObject(&wall)
	game.Init()
}
