package player

import (
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/camera"
	"github.com/reonardoleis/fcg-glcraft/engine/controls"
)

type Player struct {
	Position          mgl32.Vec4
	Camera            *camera.Camera
	IsGrounded        bool
	WalkingSpeed      float32
	defaultSpeed      float32
	RunningMultiplier float32
	ControlHandler    controls.Controls
}

func NewPlayer(playerPosition mgl32.Vec4, controlHandler controls.Controls, walkingSpeed, runningMultiplier float32) Player {
	return Player{
		Position:          playerPosition,
		ControlHandler:    controlHandler,
		Camera:            nil,
		IsGrounded:        false,
		WalkingSpeed:      walkingSpeed,
		RunningMultiplier: runningMultiplier,
		defaultSpeed:      walkingSpeed,
	}
}

func (p *Player) SetCamera(camera *camera.Camera) {
	p.Camera = camera
}

func (p *Player) Update() {
	w, u := p.Camera.GetWU()

	newPosition := p.Position

	if p.ControlHandler.IsDown(int(glfw.KeyW)) {
		newPosition = p.Position.Add(w.Mul(-1).Mul(p.WalkingSpeed))
		fmt.Println(w)
	}
	if p.ControlHandler.IsDown(int(glfw.KeyS)) {
		newPosition = p.Position.Add(w.Mul(p.WalkingSpeed))

	}
	if p.ControlHandler.IsDown(int(glfw.KeyD)) {
		newPosition = p.Position.Add(u.Mul(p.WalkingSpeed))

	}
	if p.ControlHandler.IsDown(int(glfw.KeyA)) {
		newPosition = p.Position.Add(u.Mul(-1).Mul(p.WalkingSpeed))

	}

	if p.ControlHandler.IsToggled(int(glfw.KeyLeftShift)) {
		p.WalkingSpeed = p.defaultSpeed * p.RunningMultiplier
	} else {
		p.WalkingSpeed = p.defaultSpeed
	}

	p.Position = newPosition
	p.Camera.Follow(p.Position)
}
