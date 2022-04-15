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

type CameraType int32

const (
	FirstPersonCamera CameraType = iota
)

var (
	ActiveCamera *Camera
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
	Type           CameraType
	view           mgl32.Mat4
	projection     mgl32.Mat4
}

func NewCamera(cameraPosition mgl32.Vec4, controlHandler controls.Controls, fov float32, cameraType CameraType) *Camera {
	newCamera := &Camera{
		Position:       cameraPosition,
		ViewVector:     mgl32.Vec4{0.0, 0.0, 0.0, 0.0},
		UpVector:       mgl32.Vec4{0.0, 1.0, 0.0, 0.0},
		ControlHandler: controlHandler,
		CameraDistance: 2.5,
		Fov:            fov,
		Near:           -0.1,
		Far:            -60.0,
		CameraTheta:    0.0,
		CameraPhi:      0.0,
		Type:           cameraType,
	}

	ActiveCamera = newCamera

	return newCamera
}

type Frustum struct {
}

func (c *Camera) GetFrustum() Frustum {
	/*nearCenter := c.Position.Sub(c.ViewVector.Mul(c.Near))
	farCenter := c.Position.Sub(c.ViewVector.Mul(c.Far))
	nearHeight := 2 * float32(math.Tan(float64(c.Fov)/2)) * c.Near
	farHeight := 2 * float32(math.Tan(float64(c.Fov)/2)) * c.Far
	nearWidth := nearHeight * window.ScreenRatio
	farWidth := farHeight * window.ScreenRatio

	_, u := c.GetWU()

	farTopLeft := farCenter.Add(c.UpVector.Mul(farHeight * 0.5).Sub(u.Mul(farWidth * 0.5)))
	farTopRight := farCenter.Add(c.UpVector.Mul(farHeight * 0.5).Add(u.Mul(farWidth * 0.5)))
	farBottomLeft := farCenter.Sub(c.UpVector.Mul(farHeight * 0.5).Sub(u.Mul(farWidth * 0.5)))
	farBottomRight := farCenter.Sub(c.UpVector.Mul(farHeight * 0.5).Add(u.Mul(farWidth * 0.5)))

	nearTopLeft := nearCenter + camY*(nearHeight*0.5) - camX*(nearWidth*0.5)
	nearTopRight := nearCenter + camY*(nearHeight*0.5) + camX*(nearWidth*0.5)
	nearBottomLeft := nearCenter - camY*(nearHeight*0.5) - camX*(nearWidth*0.5)
	nearBottomRight := nearCenter - camY*(nearHeight*0.5) + camX*(nearWidth*0.5)*/
	return Frustum{}
}

func (c *Camera) HandleFirstPersonCamera() {
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

	c.ViewVector = mgl32.Vec4{vx, vy, vz, 0.0}.Normalize()
	c.UpVector = mgl32.Vec4{0.0, 1.0, 0.0, 0.0}
}

func (c *Camera) Update() {
	switch c.Type {
	case FirstPersonCamera:
		c.HandleFirstPersonCamera()
	default:
		c.HandleFirstPersonCamera()
	}
	c.Handle()
}

func (c *Camera) Handle() {
	viewUniform := gl.GetUniformLocation(shaders.ShaderProgramDefault, gl.Str("view\000"))             // Variável da matriz "view" em shader_vertex.glsl
	projectionUniform := gl.GetUniformLocation(shaders.ShaderProgramDefault, gl.Str("projection\000")) // Variável da matriz "projection" em shader_vertex.glsl

	c.view = math2.Matrix_Camera_View(c.Position, c.ViewVector, c.UpVector)

	nearplane := float32(-0.1)
	farplane := float32(-60.0)

	fov := float32(math.Pi / 3.0)
	c.projection = math2.Matrix_Perspective(fov, float32(window.ScreenRatio), nearplane, farplane)

	gl.UniformMatrix4fv(viewUniform, 1, false, &c.view[0])
	gl.UniformMatrix4fv(projectionUniform, 1, false, &c.projection[0])
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

func (c *Camera) Follow(p mgl32.Vec4) {
	c.Position = p
}
