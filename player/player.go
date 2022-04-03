package player

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/camera"
	"github.com/reonardoleis/fcg-glcraft/engine/controls"
	"github.com/reonardoleis/fcg-glcraft/game_objects"
	math2 "github.com/reonardoleis/fcg-glcraft/math"
	"github.com/reonardoleis/fcg-glcraft/world"
)

type Player struct {
	Position          mgl32.Vec4
	Camera            *camera.Camera
	IsGrounded        bool
	WalkingSpeed      float32
	JumpHeight        float32
	JumpSpeed         float32
	_isJumping        bool
	_originalY        float32
	defaultSpeed      float32
	RunningMultiplier float32
	ControlHandler    controls.Controls
	PlayerPhi         float64
	PlayerTheta       float64
	MovementVector    mgl32.Vec4
	Height            float32
}

func NewPlayer(playerPosition mgl32.Vec4, controlHandler controls.Controls, walkingSpeed, runningMultiplier, jumpHeight, jumpSpeed, height float32) Player {
	return Player{
		Position:          playerPosition,
		ControlHandler:    controlHandler,
		Camera:            nil,
		IsGrounded:        false,
		WalkingSpeed:      walkingSpeed,
		RunningMultiplier: runningMultiplier,
		defaultSpeed:      walkingSpeed,
		PlayerPhi:         0.0,
		PlayerTheta:       0.0,
		MovementVector:    mgl32.Vec4{},
		JumpHeight:        jumpHeight,
		JumpSpeed:         jumpSpeed,
		_isJumping:        false,
		_originalY:        playerPosition.Y(),
		Height:            height,
	}
}

func (p Player) IsJumping() bool {
	return p._isJumping
}

func (p *Player) SetCamera(camera *camera.Camera) {
	p.Camera = camera
}

func (p *Player) SetPosition(newPosition mgl32.Vec4) {
	p.Position = newPosition
}

func (p *Player) HandleLookDirection() {
	if p.ControlHandler.MousePositionChanged() {
		dx, _ := p.ControlHandler.GetMouseDeltas()
		p.PlayerTheta -= 0.01 * dx
		// p.PlayerPhi += 0.01 * dy
	}

	vx := float32(math.Cos(float64(p.PlayerPhi)) * math.Sin(float64(p.PlayerTheta)))
	vy := float32(math.Sin(float64(0)))
	vz := float32(math.Cos(float64(p.PlayerPhi)) * math.Cos(float64(p.PlayerTheta)))
	// fmt.Println(vx, vy, vz)
	p.MovementVector = mgl32.Vec4{vx, vy, vz, 0.0}
}

func (p Player) GetMovementVector() (w, u mgl32.Vec4) {
	w = p.MovementVector.Mul(-1).Normalize()
	u = math2.Crossproduct(math2.UpVector, w)
	u = u.Normalize()
	return w, u
}

func (p *Player) Jump() {
	p._originalY = p.Position.Y()
	p._isJumping = true
	p.IsGrounded = false
}

func (p *Player) HandleJump() {
	p.SetPosition(mgl32.Vec4{p.Position.X(), p.Position.Y() + p.JumpSpeed*float32(math2.DeltaTime), p.Position.Z(), 1.0})
	if p.Position.Y()-p._originalY >= p.JumpHeight {
		p._isJumping = false
	}
}

func (p *Player) Update(world *world.World) {
	p.HandleLookDirection()
	w, u := p.GetMovementVector()

	newPosition := p.Position

	if p.ControlHandler.IsDown(int(glfw.KeyW)) {
		newPosition = p.Position.Add(w.Mul(-1).Mul(p.WalkingSpeed * float32(math2.DeltaTime)))
	}
	if p.ControlHandler.IsDown(int(glfw.KeyS)) {
		newPosition = p.Position.Add(w.Mul(p.WalkingSpeed * float32(math2.DeltaTime)))

	}
	if p.ControlHandler.IsDown(int(glfw.KeyD)) {
		newPosition = p.Position.Add(u.Mul(p.WalkingSpeed * float32(math2.DeltaTime)))

	}
	if p.ControlHandler.IsDown(int(glfw.KeyA)) {
		newPosition = p.Position.Add(u.Mul(-1).Mul(p.WalkingSpeed * float32(math2.DeltaTime)))
	}

	if p.ControlHandler.IsToggled(int(glfw.KeyLeftShift)) {
		p.WalkingSpeed = p.defaultSpeed * p.RunningMultiplier
	} else {
		p.WalkingSpeed = p.defaultSpeed
	}

	p.Position = newPosition

	if p.IsJumping() {
		p.HandleJump()
	}
	if p.ControlHandler.IsDown(int(glfw.KeySpace)) && !p.IsJumping() && p.IsGrounded {
		p.Jump()
	}

	p.HandleWorldLimits(world)

	p.Camera.Follow(p.Position)
}

func (p *Player) HandleWorldLimits(world *world.World) {
	worldSizeX, _, worldSizeZ := world.Size.X(), world.Size.Y(), world.Size.Z()
	roundedPlayerX := int(math.Round(float64(p.Position.X())))

	roundedPlayerZ := int(math.Round(float64(p.Position.Z())))
	if roundedPlayerX < -int(worldSizeX) {
		roundedPlayerX = -int(worldSizeX)
		p.SetPosition(mgl32.Vec4{-float32(worldSizeX), p.Position.Y(), p.Position.Z(), 1.0})
	}

	if roundedPlayerX > int(worldSizeX)-1 {
		roundedPlayerX = int(worldSizeX - 1)
		p.SetPosition(mgl32.Vec4{float32(worldSizeX) - 1, p.Position.Y(), p.Position.Z(), 1.0})
	}

	if roundedPlayerZ < int(-worldSizeZ) {
		roundedPlayerZ = int(-worldSizeZ)
		p.SetPosition(mgl32.Vec4{p.Position.X(), p.Position.Y(), -float32(worldSizeZ), 1.0})
	}

	if roundedPlayerZ > int(worldSizeZ-1) {
		roundedPlayerZ = int(worldSizeZ - 1)
		p.SetPosition(mgl32.Vec4{p.Position.X(), p.Position.Y(), float32(worldSizeZ) - 1, 1.0})
	}
}

func (p Player) GetRoundedPosition() (int, int, int) {
	roundedX := int(math.Round(float64(p.Position.X())))
	roundedY := int(math.Round(float64(p.Position.Y())))
	roundedZ := int(math.Round(float64(p.Position.Z())))

	return roundedX, roundedY, roundedZ
}

func (p *Player) Fall(blockBelow *game_objects.GameObject) {
	if blockBelow != nil && float64(p.Position.Y()-p.Height) <= float64(blockBelow.Position.Y()+(blockBelow.Size/2)) {
		p.IsGrounded = true
		// player.Position = (mgl32.Vec4{player.Position.X(), blockBelow.Position.Y() + (blockBelow.Size / 2) + float32(playerHeight), player.Position.Z(), 1.0})
	} else if !p.IsJumping() || blockBelow == nil {
		p.IsGrounded = false
		p.SetPosition(mgl32.Vec4{p.Position.X(), p.Position.Y() - math2.GravityAccel*float32(math2.DeltaTime), p.Position.Z(), 1.0})
	}

}
