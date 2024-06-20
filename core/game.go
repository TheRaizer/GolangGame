package core

import (
	"fmt"
	"image"

	"github.com/TheRaizer/GolangGame/core/collision"
	"github.com/TheRaizer/GolangGame/core/objs"
	"github.com/TheRaizer/GolangGame/display"
	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	collisionSys System[*collision.Collider]

	gameObjects map[string]objs.GameObject
	screen      image.Gray
	surface     *sdl.Surface
	running     bool
	window      *sdl.Window
}

func NewGame(img image.Gray, collisionSys System[*collision.Collider]) Game {
	return Game{
		collisionSys: collisionSys,
		screen:       img,
		running:      false,
		gameObjects:  make(map[string]objs.GameObject),
	}
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

	surface, err := window.GetSurface()

	if err != nil {
		panic(err)
	}
	game.surface = surface

	for _, gameObject := range game.gameObjects {
		gameObject.OnInit(game.surface)
	}

	game.running = true
	game.loop()
}

func (game *Game) render() {
	game.window.UpdateSurface()
}

func (game *Game) loop() {
	var msPerUpdate int = 1000 / display.FRAMERATE
	var current, elapsed int

	previous := int(sdl.GetTicks64())
	var lag int = 0

	for game.running {
		current = int(sdl.GetTicks64())
		elapsed = current - previous
		previous = current

		if elapsed > 1000 {
			continue
		}

		// compute the amount of lag since the last update call
		lag += elapsed

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			game.handleEvent(event)
		}

		// once the lag has reached exceeded the expected update time
		// catch up on any lag by a constant delta time value (msPerUpdate)
		for lag >= msPerUpdate {
			for _, gameObject := range game.gameObjects {
				gameObject.OnUpdate(uint64(msPerUpdate), game.surface)
			}

			game.collisionSys.OnLoop()
			lag -= msPerUpdate
		}

		game.render()
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

func (game *Game) AddGameObject(gameObject objs.GameObject) {
	game.gameObjects[gameObject.GetID()] = gameObject
}

func (game *Game) RemoveGameObject(gameObject objs.GameObject) {
	delete(game.gameObjects, gameObject.GetID())
}
