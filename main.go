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

// TODO: refactor this into separate files
func main() {
	globalRect := quadtree.Rect{X: 0, Y: 0, W: display.WIDTH, H: display.HEIGHT}
	collisionSys := collision.NewCollisionSystem(globalRect)

	// generate a gray image
	img := image.NewGray(image.Rectangle{Max: image.Point{X: display.WIDTH, Y: display.HEIGHT}})
	game := game.NewGame(*img, &collisionSys)

	playerCollider := collision.NewCollider(
		core.PLAYER_LAYER,
		"player_collider",
		quadtree.Rect{X: 0, Y: 0, W: 32, H: 32},
		&collisionSys,
		&collisionSys,
		[]func(els []quadtree.QuadElement){},
		&game,
	)

	rb := objs.NewRigidBody(core.PLAYER_LAYER, "player_rb", util.Vec2[float32]{}, &game, playerCollider, &collisionSys, true)
	player := entities.NewPlayer("player", util.Vec2[float32]{X: 0, Y: 0}, 200, &game, &rb)

	playerFloorCollider := collision.NewCollider(
		core.PLAYER_LAYER,
		"player_floor_collider",
		quadtree.Rect{X: 2, Y: 32, W: 30, H: 3}, // at the bottom of the player but not quite the entire width
		&collisionSys,
		&collisionSys,
		[]func(els []quadtree.QuadElement){
			func(els []quadtree.QuadElement) {
				for _, el := range els {
					obj := game.GetGameObject(el.Id)
					// if colliding with something not the player, then allow a jump
					if obj.Layer() != core.PLAYER_LAYER && rb.Velocity.Y > 0 {
						player.CanJump = true
					}
				}
			},
		},
		&game,
	)

	player.AddChild(playerCollider)
	player.AddChild(playerFloorCollider)
	player.AddChild(&rb)

	var wallWidth int32 = 300
	var wallHeight int32 = 32
	floor := objs.NewSolid("floor_1", util.Vec2[float32]{X: 0, Y: 500}, &game, wallWidth, wallHeight)
	floor.AddChild(collision.NewCollider(
		core.WALL_LAYER,
		"floor_1_collider",
		quadtree.Rect{X: 0, Y: 0, W: wallWidth, H: wallHeight},
		&collisionSys,
		&collisionSys,
		make([]func(els []quadtree.QuadElement), 0),
		&game,
	))
	wall := objs.NewSolid("wall_1", util.Vec2[float32]{X: 300, Y: 468}, &game, 32, 32)
	wall.AddChild(collision.NewCollider(
		core.WALL_LAYER,
		"wall_1_collider",
		quadtree.Rect{X: 0, Y: 0, W: 32, H: 32},
		&collisionSys,
		&collisionSys,
		make([]func(els []quadtree.QuadElement), 0),
		&game,
	))

	game.AddGameObject(&player)
	game.AddGameObject(&floor)
	game.AddGameObject(&wall)
	game.Init()
}
