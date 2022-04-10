package block

import (
	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/engine/shaders"
	"github.com/reonardoleis/fcg-glcraft/geometry"
	math2 "github.com/reonardoleis/fcg-glcraft/math"
)

type BoundingBoxObject struct {
	Position  mgl32.Vec4
	Model     mgl32.Mat4
	Width     float32
	Height    float32
	RotationY float32
	VaoID     uint32
}

func NewBoundingBoxObject(position mgl32.Vec4, width, height float32) BoundingBoxObject {
	cube := geometry.BuildCube(0.0, 0.0, 0.0, height, 1.0, 0.0, 0.0)
	return BoundingBoxObject{
		Position: position,
		Model:    math2.Matrix_Identity(),
		Width:    width,
		Height:   height,
		VaoID:    cube.VaoID,
	}
}

func (bbo BoundingBoxObject) Draw(rotY float32) {
	model_uniform := gl.GetUniformLocation(shaders.ShaderProgramDefault, gl.Str("model\000")) // Variável da matriz "model"
	black := gl.GetUniformLocation(shaders.ShaderProgramDefault, gl.Str("black\000"))         // Variável da matriz "model"

	rotmat := math2.Matrix_Rotate_Y(rotY)
	faceMat := math2.Matrix_Identity().Mul4(math2.Matrix_Translate(bbo.Position.X(), bbo.Position.Y(), bbo.Position.Z())).Mul4(rotmat)
	gl.BindVertexArray(bbo.VaoID)
	gl.UniformMatrix4fv(model_uniform, 1, false, &faceMat[0])
	gl.Uniform1i(black, 1)
	gl.DrawElements(
		uint32(geometry.CommonFaceEdgeGeometry.RenderingMode), // Veja slides 124-130 do documento Aula_04_Modelagem_Geometrica_3D.pdf
		int32(56),
		gl.UNSIGNED_INT,
		geometry.CommonFaceEdgeGeometry.FirstIndex,
	)
	math.Min(0, 3)

}
