package game

import (
	"image"

	"github.com/TheRaizer/GolangGame/core"
	"github.com/TheRaizer/GolangGame/core/collision"
	"github.com/TheRaizer/GolangGame/display"
	"github.com/TheRaizer/GolangGame/util"
	"github.com/veandco/go-sdl2/sdl"
)

type Game struct {
	collisionSys core.System[*collision.Collider]

	gameObjects map[string]core.GameObject
	screen      image.Gray
	surface     *sdl.Surface
	running     bool
	window      *sdl.Window
}

func NewGame(img image.Gray, collisionSys core.System[*collision.Collider]) Game {
	return Game{
		collisionSys: collisionSys,
		screen:       img,
		running:      false,
		gameObjects:  make(map[string]core.GameObject),
	}
}

func (game *Game) Init() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow(display.TITLE, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, display.WIDTH, display.HEIGHT, sdl.WINDOW_SHOWN)
	util.CheckErr(err)
	game.window = window

	surface, err := window.GetSurface()
	util.CheckErr(err)
	game.surface = surface
	renderer, err := window.GetRenderer()
	util.CheckErr(err)

	for _, gameObject := range game.gameObjects {
		gameObject.OnInit(game.surface, renderer)
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

		// handle input
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			game.handleEvent(event)
		}

		// once the lag has reached exceeded the expected update time
		// catch up on any lag by a constant delta time value (msPerUpdate)
		for lag >= msPerUpdate {
			game.collisionSys.OnLoop()

			for _, gameObject := range game.gameObjects {
				gameObject.OnUpdate(uint64(msPerUpdate), game.surface)
			}

			lag -= msPerUpdate
		}

		game.render()
	}

	game.Quit()
}

func (game *Game) handleEvent(event sdl.Event) {
	switch event.(type) {
	case *sdl.QuitEvent:
		game.Quit()
		break

	}
	for _, gameObject := range game.gameObjects {
		gameObject.OnInput(event)
	}
}

func (game *Game) Quit() {
	defer sdl.Quit()
	defer game.window.Destroy()
	game.running = false
}

func (game *Game) AddGameObject(gameObject core.GameObject) {
	if game.gameObjects[gameObject.ID()] != nil {
		panic("duplicate id: " + gameObject.ID())
	}
	game.gameObjects[gameObject.ID()] = gameObject
}

func (game *Game) RemoveGameObject(id string) {
	delete(game.gameObjects, id)
}

func (game *Game) GetGameObject(id string) core.GameObject {
	return game.gameObjects[id]
}
