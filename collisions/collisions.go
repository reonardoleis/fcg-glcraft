package collisions

import (
	"github.com/go-gl/mathgl/mgl32"
)

type CubeBoundingBox struct {
	Maxes  mgl32.Vec3
	Mins   mgl32.Vec3
	width  float32
	height float32
}

func NewCubeBoundingBox(position mgl32.Vec3, width, height float32) *CubeBoundingBox {
	maxes := mgl32.Vec3{position.X() + width/2, position.Y() + height/2, position.Z() + width/2}
	mins := mgl32.Vec3{position.X() - width/2, position.Y() - height/2, position.Z() - width/2}
	return &CubeBoundingBox{
		Maxes:  maxes,
		Mins:   mins,
		width:  width,
		height: height,
	}
}

func (cbb *CubeBoundingBox) UpdateBounds(position mgl32.Vec3) {
	cbb.Maxes = mgl32.Vec3{position.X() + cbb.width/2, position.Y() + cbb.height/2, position.Z() + cbb.width/2}
	cbb.Mins = mgl32.Vec3{position.X() - cbb.width/2, position.Y() - cbb.height/2, position.Z() - cbb.width/2}
}
