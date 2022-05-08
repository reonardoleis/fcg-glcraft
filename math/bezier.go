package math2

import (
	"github.com/go-gl/mathgl/mgl32"
)

type BezierCurve struct {
	ControlPoints []mgl32.Vec3
}

func NewBezierCurve() BezierCurve {
	return BezierCurve{
		ControlPoints: make([]mgl32.Vec3, 4),
	}
}

// Generate 4 random control points
func (bc *BezierCurve) GenerateRandomPoints() {
	for i := 0; i < 4; i++ {
		bc.ControlPoints[i] = mgl32.Vec3{float32(RandInt(0, 3)), 0, float32(RandInt(0, 3))}
	}
}

// Compute the bezier curve at T
func (bc BezierCurve) T(t float32) mgl32.Vec3 {
	// T in [0, 3]
	cp1 := bc.ControlPoints[0]
	cp2 := bc.ControlPoints[1]
	cp3 := bc.ControlPoints[2]
	cp4 := bc.ControlPoints[3]
	pz := (Pow3(1-t) * cp1.X() / 3) + (3 * t * Pow2(1-t) * cp2.X() / 3) + (3 * Pow2(t) * (1 - t) * cp3.X() / 3) + Pow3(t)*cp4.X()/3

	return mgl32.Vec3{t, 0, pz}
}
