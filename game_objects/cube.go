package game_objects

import (
	"github.com/reonardoleis/fcg-glcraft/engine/renderer"
	"github.com/reonardoleis/fcg-glcraft/geometry"
	math2 "github.com/reonardoleis/fcg-glcraft/math"
)

var (
	CubeEdgesOnly bool = false
)

type Cube = GameObject

func NewCube(x, y, z, size float32, withEdges bool) Cube {
	vaoID, sceneObject := geometry.BuildCube(0, 0, 0, size)
	sceneObject.VaoID = vaoID
	edges := renderer.SceneObject{}

	model := math2.Matrix_Identity().Mul4(math2.Matrix_Translate(x, y, z))

	if withEdges {
		vaoID, edges = geometry.BuildCubeEdges(0, 0, 0, size)
		edges.VaoID = vaoID
	}
	return Cube{
		X:           x,
		Y:           y,
		Z:           z,
		Size:        size,
		Model:       model,
		SceneObject: sceneObject,
		WithEdges:   withEdges,
		Edges:       edges,
	}
}

func (c *Cube) Translate(x, y, z float32) {
	c.Model = math2.Matrix_Identity().Mul4(math2.Matrix_Translate(x, y, z))
}
