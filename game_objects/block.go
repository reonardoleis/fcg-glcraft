package game_objects

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/engine/renderer"
	"github.com/reonardoleis/fcg-glcraft/geometry"
	math2 "github.com/reonardoleis/fcg-glcraft/math"
)

var (
	BlockEdgesOnly bool = false
)

type Block = GameObject

func NewBlock(x, y, z, size float32, withEdges, ephemeral bool, color mgl32.Vec3) Block {
	vaoID, sceneObject := geometry.BuildCube(0, 0, 0, size, color.X(), color.Y(), color.Z())
	sceneObject.VaoID = vaoID
	edges := renderer.SceneObject{}

	model := math2.Matrix_Identity().Mul4(math2.Matrix_Translate(x, y, z))

	if withEdges {
		vaoID, edges = geometry.BuildCubeEdges(0, 0, 0, size)
		edges.VaoID = vaoID
	}
	return Block{
		Position:    mgl32.Vec4{x, y, z, 0.0},
		Size:        size,
		Model:       model,
		SceneObject: sceneObject,
		WithEdges:   withEdges,
		Edges:       edges,
		Ephemeral:   ephemeral,
	}
}

func (c *Block) Translate(x, y, z float32) {
	c.Model = math2.Matrix_Identity().Mul4(math2.Matrix_Translate(x, y, z))
}
