package systems

import (
	"github.com/TheRaizer/GolangGame/core/objs"
)

func ApplyGravity(dt uint64, rb *objs.RigidBody) {
	dtSec := float32(dt)

	// control max downward velocity
	if rb.Velocity.Y <= 400 {
		rb.Velocity.Y += 9.81 * dtSec
	}
}
