package core

import (
	"fmt"
	"image"

	"github.com/TheRaizer/GolangGame/display"
	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	gameObjects map[string]GameObject
	screen      image.Gray
	running     bool
	window      *sdl.Window
}

func NewGame(img image.Gray) *Game {
	return &Game{screen: img, running: false, gameObjects: make(map[string]GameObject)}
}

func (game *Game) Init() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(display.TITLE, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, display.WIDTH, display.HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	game.window = window

	for _, gameObject := range game.gameObjects {
		gameObject.OnInit()
	}

	game.running = true
	game.loop()
}

func (game *Game) render() {
	game.window.UpdateSurface()
}

func (game *Game) loop() {
	var lastTime uint64 = sdl.GetTicks64()
	for game.running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			game.handleEvent(event)
		}

		currentTime := sdl.GetTicks64()
		dt := currentTime - lastTime
		lastTime = currentTime

		for _, gameObject := range game.gameObjects {
			gameObject.OnUpdate(dt)
		}
		game.render()

		delay := (1000 / display.FRAMERATE) - (sdl.GetTicks64() - currentTime)

		fmt.Println(delay)
		sdl.Delay(uint32(delay))
	}

	game.Quit()
}

func (game *Game) handleEvent(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.QuitEvent:
		game.Quit()
		break
	case *sdl.KeyboardEvent:
		fmt.Println(t)
		break
	}
}

func (game *Game) Quit() {
	defer sdl.Quit()
	defer game.window.Destroy()
	game.running = false
}

func (game *Game) AddGameObject(gameObject GameObject) {
	game.gameObjects[gameObject.GetID()] = gameObject
}

func (game *Game) RemoveGameObject(gameObject GameObject) {
	delete(game.gameObjects, gameObject.GetID())
}
