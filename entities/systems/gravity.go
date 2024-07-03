package systems

import (
	"github.com/TheRaizer/GolangGame/core/objs"
)

func ApplyGravity(dt uint64, rb *objs.RigidBody) {
	dtSec := float32(dt)

	// max downward velocity of 500
	if rb.Velocity.Y <= 300 {
		rb.Velocity.Y += 9.81 * dtSec
	}
}
