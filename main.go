package main

import (
	"image"

	"github.com/TheRaizer/GolangGame/core"
	"github.com/TheRaizer/GolangGame/core/collision"
	"github.com/TheRaizer/GolangGame/display"
	"github.com/TheRaizer/GolangGame/entities"
	"github.com/TheRaizer/GolangGame/util"
	"github.com/TheRaizer/GolangGame/util/quadtree"
)

func main() {
	globalRect := quadtree.Rect{X: 0, Y: 0, W: display.WIDTH, H: display.HEIGHT}
	collisionSys := collision.NewCollisionSystem(globalRect)

	// generate a gray image
	img := image.NewGray(image.Rectangle{Max: image.Point{X: display.WIDTH, Y: display.HEIGHT}})
	game := core.NewGame(*img, &collisionSys)

	player := entities.NewPlayer("player", util.Vec2[float32]{X: 0, Y: 0})
	player.AddChild(collision.NewCollider(
		"player_collider",
		quadtree.Rect{X: 0, Y: 0, W: 32, H: 32},
		&collisionSys,
	))

	game.AddGameObject(&player)
	game.Init()
}
