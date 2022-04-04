package game_objects

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/geometry"
	math2 "github.com/reonardoleis/fcg-glcraft/math"
)

type CastedVector = GameObject

func NewCastedVector(a, b mgl32.Vec3) CastedVector {
	vaoID, sceneObject := geometry.BuildLine(a, b)
	sceneObject.VaoID = vaoID

	model := math2.Matrix_Identity().Mul4(math2.Matrix_Translate(0, 0, 0))

	return CastedVector{
		Model:       model,
		SceneObject: sceneObject,
	}
}
