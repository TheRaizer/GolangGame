package systems

import (
	"github.com/TheRaizer/GolangGame/core/objs"
)

func ApplyGravity(dt uint64, rb *objs.RigidBody, multiplier float32) {
	// control max downward velocity
	if rb.Velocity.Y <= 300 {
		rb.Velocity.Y += 9.81 * float32(dt) * multiplier
	}
}
