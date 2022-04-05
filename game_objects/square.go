package game_objects

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/geometry"
)

type Square struct {
	Position      mgl32.Vec4
	Size          float32
	Model         mgl32.Mat4
	ModelGeometry geometry.GeometryInformation
	WithEdges     bool
	EdgesGeometry geometry.GeometryInformation
}
