package main

import (
	"image"

	"github.com/TheRaizer/GolangGame/core"
	"github.com/TheRaizer/GolangGame/display"
	"github.com/TheRaizer/GolangGame/entities"
)

func main() {
	// generate a gray image
	img := image.NewGray(image.Rectangle{Max: image.Point{X: display.WIDTH, Y: display.HEIGHT}})
	game := core.NewGame(*img)

	player := entities.NewPlayer(core.Vector{X: 0, Y: 0})

	game.AddGameObject(player)
	game.Init()
}
