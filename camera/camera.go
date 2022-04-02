package camera

import (
	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/engine/controls"
	"github.com/reonardoleis/fcg-glcraft/engine/shaders"
	"github.com/reonardoleis/fcg-glcraft/engine/window"
	math2 "github.com/reonardoleis/fcg-glcraft/math"
)

type Camera struct {
	Position       mgl32.Vec4
	ViewVector     mgl32.Vec4
	UpVector       mgl32.Vec4
	ControlHandler controls.Controls
	CameraDistance float64
	Fov            float32
	Near           float32
	Far            float32
	CameraTheta    float64
	CameraPhi      float64
	View           mgl32.Mat4
	Projection     mgl32.Mat4
}

func (c *Camera) Init() {
	c.ViewVector = mgl32.Vec4{0.0, 0.0, 0.0, 0.0}
	c.UpVector = mgl32.Vec4{0.0, 1.0, 0.0, 0.0}
}

func (c *Camera) Update() {
	if c.ControlHandler.MousePositionChanged() {
		dx, dy := c.ControlHandler.GetMouseDeltas()
		c.CameraTheta -= 0.01 * dx
		c.CameraPhi -= 0.01 * dy

		phiMax := math.Pi / 2
		phiMin := -phiMax

		if c.CameraPhi > phiMax {
			c.CameraPhi = phiMax
		}

		if c.CameraPhi < phiMin {
			c.CameraPhi = phiMin
		}
	}

	r := float32(c.CameraDistance)
	vx := r * float32(math.Cos(float64(c.CameraPhi))*math.Sin(float64(c.CameraTheta)))
	vy := r * float32(math.Sin(float64(c.CameraPhi)))
	vz := r * float32(math.Cos(float64(c.CameraPhi))*math.Cos(float64(c.CameraTheta)))

	c.ViewVector = mgl32.Vec4{vx, vy, vz, 0.0}
	c.UpVector = mgl32.Vec4{0.0, 1.0, 0.0, 0.0}
}

func (c *Camera) Handle() {
	viewUniform := gl.GetUniformLocation(shaders.ShaderProgram, gl.Str("view\000"))             // Variável da matriz "view" em shader_vertex.glsl
	projectionUniform := gl.GetUniformLocation(shaders.ShaderProgram, gl.Str("projection\000")) // Variável da matriz "projection" em shader_vertex.glsl

	c.View = math2.Matrix_Camera_View(c.Position, c.ViewVector, c.UpVector)

	nearplane := float32(-0.1)
	farplane := float32(-30.0)

	fov := float32(math.Pi / 3.0)
	c.Projection = math2.Matrix_Perspective(fov, float32(window.ScreenRatio), nearplane, farplane)

	gl.UniformMatrix4fv(viewUniform, 1, false, &c.View[0])
	gl.UniformMatrix4fv(projectionUniform, 1, false, &c.Projection[0])

}

func (c *Camera) SetPosition(position mgl32.Vec4) {
	c.Position = position
}

func (c Camera) GetWU() (w, u mgl32.Vec4) {
	w = c.ViewVector.Mul(-1).Normalize()
	u = math2.Crossproduct(c.UpVector, w)
	u = u.Normalize()
	return w, u
}
