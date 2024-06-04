package main

import (
	"image"

	"github.com/TheRaizer/GolangGame/core"
	"github.com/TheRaizer/GolangGame/display"
	"github.com/TheRaizer/GolangGame/entities"
)

// you're player is blue.
// you're enemy is red.
// you want to pick up green.
// you cannot pass black (walls).

func main() {
	// generate a gray image
	img := image.NewGray(image.Rectangle{Max: image.Point{X: display.WIDTH, Y: display.HEIGHT}})
	game := core.NewGame(*img)

	player := entities.NewPlayer(core.Vector{X: 0, Y: 0})

	game.AddGameObject(player)
	game.Init()
}
