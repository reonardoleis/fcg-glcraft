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
		Near:           0.1,
		Far:            40.0,
		CameraTheta:    0.0,
		CameraPhi:      0.0,
		Type:           cameraType,
	}

	ActiveCamera = newCamera

	return newCamera
}

type Frustum struct {
	Ftl mgl32.Vec4
	Ftr mgl32.Vec4
	Fbl mgl32.Vec4
	Fbr mgl32.Vec4
	Ntl mgl32.Vec4
	Ntr mgl32.Vec4
	Nbl mgl32.Vec4
	Nbr mgl32.Vec4
}

func (c *Camera) GetFrustum() Frustum {
	vv := c.ViewVector
	/* float a = cam.nearClipPlane;//get length
	   float A = cam.fieldOfView * 0.5f;//get angle
	   A = A * Mathf.Deg2Rad;//convert tor radians
	   float h = (Mathf.Tan(A) * a);//calc height
	   float w = (h / cam.pixelHeight) * cam.pixelWidth;//deduct width
	*/

	Hfar := 2 * math.Tan(float64(c.Fov)/2) * float64(c.Far)
	Hnear := 2 * math.Tan(float64(c.Fov)/2) * float64(c.Near)

	Wnear := (Hnear / 720) * 1280
	Wfar := (Hfar / 720) * 1280

	fc := c.Position.Add(vv.Mul(c.Far))

	_, u := c.GetWU()

	ftl := fc.Add(c.UpVector.Mul(float32(Hfar) / 2)).Sub(u.Mul(float32(Wfar) / 2))
	ftr := fc.Add(c.UpVector.Mul(float32(Hfar) / 2)).Add(u.Mul(float32(Wfar) / 2))
	fbl := fc.Sub(c.UpVector.Mul(float32(Hfar) / 2)).Sub(u.Mul(float32(Wfar) / 2))
	fbr := fc.Sub(c.UpVector.Mul(float32(Hfar) / 2)).Add(u.Mul(float32(Wfar) / 2))

	nc := c.Position.Add(vv.Mul(c.Near))

	ntl := nc.Add(c.UpVector.Mul(float32(Hnear) / 2)).Sub(u.Mul(float32(Wnear) / 2))
	ntr := nc.Add(c.UpVector.Mul(float32(Hnear) / 2)).Add(u.Mul(float32(Wnear) / 2))
	nbl := nc.Sub(c.UpVector.Mul(float32(Hnear) / 2)).Sub(u.Mul(float32(Wnear) / 2))
	nbr := nc.Sub(c.UpVector.Mul(float32(Hnear) / 2)).Add(u.Mul(float32(Wnear) / 2))

	return Frustum{
		ftl,
		ftr,
		fbl,
		fbr,
		ntl,
		ntr,
		nbl,
		nbr,
	}
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

	nearplane := -c.Near
	farplane := -c.Far

	fov := c.Fov
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
