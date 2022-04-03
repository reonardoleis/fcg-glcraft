package game_objects

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/engine/renderer"
	"github.com/reonardoleis/fcg-glcraft/engine/shaders"
)

type GameObject struct {
	Position    mgl32.Vec4
	Size        float32
	Model       mgl32.Mat4
	SceneObject renderer.SceneObject
	WithEdges   bool
	Edges       renderer.SceneObject
}

func (g GameObject) Draw() {
	model_uniform := gl.GetUniformLocation(shaders.ShaderProgram, gl.Str("model\000"))                     // Variável da matriz "model"
	render_as_black_uniform := gl.GetUniformLocation(shaders.ShaderProgram, gl.Str("render_as_black\000")) // Variável booleana em shader_vertex.glsl

	if !CubeEdgesOnly {
		gl.BindVertexArray(g.SceneObject.VaoID)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, shaders.CubeTextureTest)

		gl.UniformMatrix4fv(model_uniform, 1, false, &g.Model[0])
		gl.Uniform1i(render_as_black_uniform, 0)
		gl.DrawElements(
			uint32(g.SceneObject.RenderingMode), // Veja slides 124-130 do documento Aula_04_Modelagem_Geometrica_3D.pdf
			int32(g.SceneObject.NumIndices),
			gl.UNSIGNED_INT,
			g.SceneObject.FirstIndex,
		)
	}

	if g.WithEdges {
		gl.BindVertexArray(g.Edges.VaoID)
		gl.UniformMatrix4fv(model_uniform, 1, false, &g.Model[0])
		gl.Uniform1i(render_as_black_uniform, 1)
		gl.DrawElements(
			uint32(g.Edges.RenderingMode), // Veja slides 124-130 do documento Aula_04_Modelagem_Geometrica_3D.pdf
			int32(g.Edges.NumIndices),
			gl.UNSIGNED_INT,
			g.Edges.FirstIndex,
		)
	}
}

func (g GameObject) GetPosition() mgl32.Vec4 {
	return g.Position
}
