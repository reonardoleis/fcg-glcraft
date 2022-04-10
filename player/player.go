package player

import (
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/block"
	"github.com/reonardoleis/fcg-glcraft/camera"
	"github.com/reonardoleis/fcg-glcraft/collisions"
	"github.com/reonardoleis/fcg-glcraft/configs"
	"github.com/reonardoleis/fcg-glcraft/world/chunk"

	"github.com/reonardoleis/fcg-glcraft/engine/controls"
	math2 "github.com/reonardoleis/fcg-glcraft/math"
	"github.com/reonardoleis/fcg-glcraft/world"
)

var (
	color = mgl32.Vec3{1.0, 0.0, 0.0}
)

type Player struct {
	Position                  mgl32.Vec4
	Camera                    *camera.Camera
	IsGrounded                bool
	WalkingSpeed              float32
	JumpHeight                float32
	JumpSpeed                 float32
	_isJumping                bool
	_originalY                float32
	defaultSpeed              float32
	RunningMultiplier         float32
	ControlHandler            controls.Controls
	PlayerPhi                 float64
	PlayerTheta               float64
	MovementVector            mgl32.Vec4
	Height                    float32
	HitAt                     *mgl32.Vec4
	ClosestEmptySpace         *mgl32.Vec4
	_mouseLeftDownLastUpdate  bool
	_mouseRightDownLastUpdate bool
	Collider                  *collisions.CubeCollider
	BoundingBox               *collisions.CubeBoundingBox
	BoundingBox2              *collisions.CubeBoundingBox
	BoundingBoxFutureVertices [8]mgl32.Vec3
	SelectedBlock             block.BlockType
	LastChunk                 uint64
}

func NewPlayer(playerPosition mgl32.Vec4, controlHandler controls.Controls, walkingSpeed, runningMultiplier, jumpHeight, jumpSpeed, height float32) Player {
	return Player{
		Position:                  playerPosition,
		ControlHandler:            controlHandler,
		Camera:                    nil,
		IsGrounded:                false,
		WalkingSpeed:              walkingSpeed,
		RunningMultiplier:         runningMultiplier,
		defaultSpeed:              walkingSpeed,
		PlayerPhi:                 0.0,
		PlayerTheta:               0.0,
		MovementVector:            mgl32.Vec4{},
		JumpHeight:                jumpHeight,
		JumpSpeed:                 jumpSpeed,
		_isJumping:                false,
		_originalY:                playerPosition.Y(),
		Height:                    height,
		HitAt:                     &mgl32.Vec4{},
		ClosestEmptySpace:         &mgl32.Vec4{},
		_mouseLeftDownLastUpdate:  true,
		_mouseRightDownLastUpdate: true,
		Collider:                  collisions.NewCubeCollider(),
		BoundingBox:               collisions.NewCubeBoundingBox(playerPosition.Vec3(), configs.PlayerWidth*0.5, configs.PlayerHeight*0.5),
		BoundingBox2:              collisions.NewCubeBoundingBox(playerPosition.Vec3(), configs.PlayerWidth, configs.PlayerHeight),
		BoundingBoxFutureVertices: [8]mgl32.Vec3{},
		SelectedBlock:             block.BlockDirt,
		LastChunk:                 0,
	}
}

func (p Player) IsJumping() bool {
	return p._isJumping
}

func (p *Player) BeFollowedByCamera(camera *camera.Camera) {
	p.Camera = camera

}

func (p *Player) SetPosition(newPosition mgl32.Vec4) {
	p.Position = newPosition
}

func (p *Player) HandleLookDirection() {
	if p.ControlHandler.MousePositionChanged() {
		dx, _ := p.ControlHandler.GetMouseDeltas()
		p.PlayerTheta -= 0.01 * dx
		p.PlayerTheta = p.PlayerTheta
		// p.PlayerPhi += 0.01 * dy
	}

	vx := float32(math.Cos(float64(p.PlayerPhi)) * math.Sin(float64(p.PlayerTheta)))
	vy := float32(math.Sin(float64(0)))
	vz := float32(math.Cos(float64(p.PlayerPhi)) * math.Cos(float64(p.PlayerTheta)))
	// // fmt.Println(vx, vy, vz)
	p.MovementVector = mgl32.Vec4{vx, vy, vz, 0.0}
}

func (p Player) GetFrontAndBackDirections() (behind mgl32.Vec3, front mgl32.Vec3) {
	return p.MovementVector.Mul(-2).Add(p.Position).Vec3(), p.MovementVector.Vec3().Add(p.Position.Vec3())
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
}

func (p *Player) HandleJump() {
	p.SetPosition(mgl32.Vec4{p.Position.X(), p.Position.Y() + p.JumpSpeed*float32(math2.DeltaTime), p.Position.Z(), 1.0})
	if p.Position.Y()-p._originalY >= p.JumpHeight {
		p._isJumping = false
	}
}

func (p *Player) CheckCollisions(roundedNewPositionX int, roundedNewPositionY int, roundedNewPositionZ int, futureBoundingBox collisions.CubeBoundingBox, chunk *chunk.Chunk) (collidesSides, collidesBelow, collidesAbove bool) {
	//fmt.Println("Checando colisoes no chunk ", chunk.ID)
	//fmt.Println(futureBoundingBox)
	//fmt.Println("Iterando...")
	//fmt.Println(playerX, playerY, playerZ)
	for x := roundedNewPositionX - 1; x <= roundedNewPositionX+1; x++ {
		for y := roundedNewPositionY - 1; y <= roundedNewPositionY+1; y++ {
			for z := roundedNewPositionZ - 1; z <= roundedNewPositionZ+1; z++ {
				//fmt.Println("Y: ", y)
				blockToVerify := chunk.GetBlockAt(x, y, z)
				if blockToVerify == nil {
					//fmt.Println("NULO:", x, y, z)
					continue
				}

				//// fmt.Println("NAO NULO..")
				cubeVertices := blockToVerify.GetFutureVertices()

				blockBoundingBox := collisions.NewCubeBoundingBox(blockToVerify.Position.Vec3(), float32(configs.BlockSize), float32(configs.BlockSize))
				//fmt.Println("Maximos: ", blockBoundingBox.Maxes, " Minimos: ", blockBoundingBox.Mins)
				//fmt.Println("Maximos: ", futureBoundingBox.Maxes, " Minimos: ", futureBoundingBox.Mins)

				//// fmt.Println("BLOCO Max: ", blockBoundingBox.Maxes)
				//// fmt.Println("BLOCO Min: ", blockBoundingBox.Mins)
				//// fmt.Println(x, playerX, y, playerY, z, playerZ)

				if y == roundedNewPositionY-1 && p.Collider.Collides(futureBoundingBox, *blockBoundingBox, p.BoundingBoxFutureVertices, cubeVertices) {
					//fmt.Println("Colidindo abaixo...")
					//fmt.Println("Maximos: ", blockBoundingBox.Maxes, " Minimos: ", blockBoundingBox.Mins)
					//fmt.Println("Maximos: ", futureBoundingBox.Maxes, " Minimos: ", futureBoundingBox.Mins)
					collidesBelow = true
					//blocks[x][y][z].Colliding = true
					continue
				}

				if p.Collider.Collides(futureBoundingBox, *blockBoundingBox, p.BoundingBoxFutureVertices, cubeVertices) {

					if y == roundedNewPositionY+1 {

						collidesAbove = true
						collidesSides = true
						continue
					} else {
						collidesSides = true

						continue

						//fmt.Println("Bloco: ", x, y, z, "Player: ", roundedNewPositionX, roundedNewPositionY, roundedNewPositionZ)
						//
					}
				}

			}

		}
	}
	//fmt.Println("Parando...")
	return
}

func (p *Player) Update(world *world.World, chunk *chunk.Chunk) {
	p.HandleLookDirection()

	w, u := p.GetMovementVector()

	newPosition := p.Position

	if !p._isJumping {
		//fmt.Println("Caindo....")
		newPosition = newPosition.Sub(mgl32.Vec4{0.0, math2.GravityAccel * float32(math2.DeltaTime), 0.0, 0.0})

	}

	bb := p.UpdateBoundingBox(newPosition)

	roundedNewPositionX := int(math.Round(float64(newPosition.X())))
	roundedNewPositionY := int(math.Round(float64(newPosition.Y())))
	roundedNewPositionZ := int(math.Round(float64(newPosition.Z())))

	_, collidedBelow, collidesAbove := p.CheckCollisions(roundedNewPositionX, roundedNewPositionY, roundedNewPositionZ, bb, chunk)

	//wneg := w.Mul(-1)
	if p.ControlHandler.IsDown(int(glfw.KeyW)) {
		newPosition = newPosition.Add(w.Mul(-1).Mul(p.WalkingSpeed * float32(math2.DeltaTime)))

		bb = p.UpdateBoundingBox(newPosition)

		roundedNewPositionX = int(math.Round(float64(newPosition.X())))
		roundedNewPositionY = int(math.Round(float64(newPosition.Y())))
		roundedNewPositionZ = int(math.Round(float64(newPosition.Z())))

		collided, _, _ := p.CheckCollisions(roundedNewPositionX, roundedNewPositionY, roundedNewPositionZ, bb, chunk)

		if collided {
			newPosition = newPosition.Sub(w.Mul(-1).Mul(p.WalkingSpeed * float32(math2.DeltaTime)))

			newPosition = newPosition.Add(mgl32.Vec4{w.Mul(-1).Mul(p.WalkingSpeed * float32(math2.DeltaTime)).X(), 0.0, 0.0, 0.0})

			bb = p.UpdateBoundingBox(newPosition)

			roundedNewPositionX = int(math.Round(float64(newPosition.X())))
			roundedNewPositionY = int(math.Round(float64(newPosition.Y())))
			roundedNewPositionZ = int(math.Round(float64(newPosition.Z())))

			collided, _, _ = p.CheckCollisions(roundedNewPositionX, roundedNewPositionY, roundedNewPositionZ, bb, chunk)

			if collided {
				newPosition = newPosition.Sub(mgl32.Vec4{w.Mul(-1).Mul(p.WalkingSpeed * float32(math2.DeltaTime)).X(), 0.0, 0.0, 0.0})

				newPosition = newPosition.Add(mgl32.Vec4{0.0, 0.0, w.Mul(-1).Mul(p.WalkingSpeed * float32(math2.DeltaTime)).Z(), 0.0})

				bb = p.UpdateBoundingBox(newPosition)

				roundedNewPositionX = int(math.Round(float64(newPosition.X())))
				roundedNewPositionY = int(math.Round(float64(newPosition.Y())))
				roundedNewPositionZ = int(math.Round(float64(newPosition.Z())))

				collided, _, _ = p.CheckCollisions(roundedNewPositionX, roundedNewPositionY, roundedNewPositionZ, bb, chunk)

				if collided {
					newPosition = newPosition.Sub(mgl32.Vec4{0.0, 0.0, w.Mul(-1).Mul(p.WalkingSpeed * float32(math2.DeltaTime)).Z(), 0.0})
				}
			}

		}

	}
	if p.ControlHandler.IsDown(int(glfw.KeyS)) {
		newPosition = newPosition.Add(w.Mul(p.WalkingSpeed * float32(math2.DeltaTime)))

		bb = p.UpdateBoundingBox(newPosition)

		roundedNewPositionX = int(math.Round(float64(newPosition.X())))
		roundedNewPositionY = int(math.Round(float64(newPosition.Y())))
		roundedNewPositionZ = int(math.Round(float64(newPosition.Z())))

		collided, _, _ := p.CheckCollisions(roundedNewPositionX, roundedNewPositionY, roundedNewPositionZ, bb, chunk)

		if collided {
			newPosition = newPosition.Sub(w.Mul(p.WalkingSpeed * float32(math2.DeltaTime)))

			newPosition = newPosition.Add(mgl32.Vec4{w.Mul(p.WalkingSpeed * float32(math2.DeltaTime)).X(), 0.0, 0.0, 0.0})

			bb = p.UpdateBoundingBox(newPosition)

			roundedNewPositionX = int(math.Round(float64(newPosition.X())))
			roundedNewPositionY = int(math.Round(float64(newPosition.Y())))
			roundedNewPositionZ = int(math.Round(float64(newPosition.Z())))

			collided, _, _ = p.CheckCollisions(roundedNewPositionX, roundedNewPositionY, roundedNewPositionZ, bb, chunk)

			if collided {
				newPosition = newPosition.Sub(mgl32.Vec4{w.Mul(p.WalkingSpeed * float32(math2.DeltaTime)).X(), 0.0, 0.0, 0.0})

				newPosition = newPosition.Add(mgl32.Vec4{0.0, 0.0, w.Mul(p.WalkingSpeed * float32(math2.DeltaTime)).Z(), 0.0})

				bb = p.UpdateBoundingBox(newPosition)

				roundedNewPositionX = int(math.Round(float64(newPosition.X())))
				roundedNewPositionY = int(math.Round(float64(newPosition.Y())))
				roundedNewPositionZ = int(math.Round(float64(newPosition.Z())))

				collided, _, _ = p.CheckCollisions(roundedNewPositionX, roundedNewPositionY, roundedNewPositionZ, bb, chunk)

				if collided {
					newPosition = newPosition.Sub(mgl32.Vec4{0.0, 0.0, w.Mul(p.WalkingSpeed * float32(math2.DeltaTime)).Z(), 0.0})
				}
			}

		}

	}
	if p.ControlHandler.IsDown(int(glfw.KeyD)) {
		newPosition = newPosition.Add(u.Mul(p.WalkingSpeed * float32(math2.DeltaTime)))

		bb = p.UpdateBoundingBox(newPosition)

		roundedNewPositionX = int(math.Round(float64(newPosition.X())))
		roundedNewPositionY = int(math.Round(float64(newPosition.Y())))
		roundedNewPositionZ = int(math.Round(float64(newPosition.Z())))

		collided, _, _ := p.CheckCollisions(roundedNewPositionX, roundedNewPositionY, roundedNewPositionZ, bb, chunk)

		if collided {
			newPosition = newPosition.Sub(u.Mul(p.WalkingSpeed * float32(math2.DeltaTime)))

			newPosition = newPosition.Add(mgl32.Vec4{u.Mul(p.WalkingSpeed * float32(math2.DeltaTime)).X(), 0.0, 0.0, 0.0})

			bb = p.UpdateBoundingBox(newPosition)

			roundedNewPositionX = int(math.Round(float64(newPosition.X())))
			roundedNewPositionY = int(math.Round(float64(newPosition.Y())))
			roundedNewPositionZ = int(math.Round(float64(newPosition.Z())))

			collided, _, _ = p.CheckCollisions(roundedNewPositionX, roundedNewPositionY, roundedNewPositionZ, bb, chunk)

			if collided {
				newPosition = newPosition.Sub(mgl32.Vec4{u.Mul(p.WalkingSpeed * float32(math2.DeltaTime)).X(), 0.0, 0.0, 0.0})

				newPosition = newPosition.Add(mgl32.Vec4{0.0, 0.0, u.Mul(p.WalkingSpeed * float32(math2.DeltaTime)).Z(), 0.0})

				bb = p.UpdateBoundingBox(newPosition)

				roundedNewPositionX = int(math.Round(float64(newPosition.X())))
				roundedNewPositionY = int(math.Round(float64(newPosition.Y())))
				roundedNewPositionZ = int(math.Round(float64(newPosition.Z())))

				collided, _, _ = p.CheckCollisions(roundedNewPositionX, roundedNewPositionY, roundedNewPositionZ, bb, chunk)

				if collided {
					newPosition = newPosition.Sub(mgl32.Vec4{0.0, 0.0, u.Mul(p.WalkingSpeed * float32(math2.DeltaTime)).Z(), 0.0})
				}
			}

		}

	}
	if p.ControlHandler.IsDown(int(glfw.KeyA)) {
		newPosition = newPosition.Add(u.Mul(-1).Mul(p.WalkingSpeed * float32(math2.DeltaTime)))

		bb = p.UpdateBoundingBox(newPosition)

		roundedNewPositionX = int(math.Round(float64(newPosition.X())))
		roundedNewPositionY = int(math.Round(float64(newPosition.Y())))
		roundedNewPositionZ = int(math.Round(float64(newPosition.Z())))

		collided, _, _ := p.CheckCollisions(roundedNewPositionX, roundedNewPositionY, roundedNewPositionZ, bb, chunk)

		if collided {
			newPosition = newPosition.Sub(u.Mul(-1).Mul(p.WalkingSpeed * float32(math2.DeltaTime)))

			newPosition = newPosition.Add(mgl32.Vec4{u.Mul(-1).Mul(p.WalkingSpeed * float32(math2.DeltaTime)).X(), 0.0, 0.0, 0.0})

			bb = p.UpdateBoundingBox(newPosition)

			roundedNewPositionX = int(math.Round(float64(newPosition.X())))
			roundedNewPositionY = int(math.Round(float64(newPosition.Y())))
			roundedNewPositionZ = int(math.Round(float64(newPosition.Z())))

			collided, _, _ = p.CheckCollisions(roundedNewPositionX, roundedNewPositionY, roundedNewPositionZ, bb, chunk)

			if collided {
				newPosition = newPosition.Sub(mgl32.Vec4{u.Mul(-1).Mul(p.WalkingSpeed * float32(math2.DeltaTime)).X(), 0.0, 0.0, 0.0})

				newPosition = newPosition.Add(mgl32.Vec4{0.0, 0.0, u.Mul(-1).Mul(p.WalkingSpeed * float32(math2.DeltaTime)).Z(), 0.0})

				bb = p.UpdateBoundingBox(newPosition)

				roundedNewPositionX = int(math.Round(float64(newPosition.X())))
				roundedNewPositionY = int(math.Round(float64(newPosition.Y())))
				roundedNewPositionZ = int(math.Round(float64(newPosition.Z())))

				collided, _, _ = p.CheckCollisions(roundedNewPositionX, roundedNewPositionY, roundedNewPositionZ, bb, chunk)

				if collided {
					newPosition = newPosition.Sub(mgl32.Vec4{0.0, 0.0, u.Mul(-1).Mul(p.WalkingSpeed * float32(math2.DeltaTime)).Z(), 0.0})
				}
			}

		}

	}
	if collidedBelow {

		newPosition = mgl32.Vec4{newPosition.X(), p.Position.Y(), newPosition.Z(), 1.0}

	}

	p.Position = newPosition

	if p.ControlHandler.IsDown(int(glfw.MouseButtonLeft)) && !p._mouseLeftDownLastUpdate && p.HitAt != nil {
		p._mouseLeftDownLastUpdate = true
		chunk.RemoveBlockFrom(*p.HitAt)
		p.HitAt = nil
	}
	if !p.ControlHandler.IsDown(int(glfw.MouseButtonLeft)) {
		p._mouseLeftDownLastUpdate = false
	}
	if p.ControlHandler.IsDown(int(glfw.MouseButtonRight)) && p.ClosestEmptySpace != nil && !p._mouseRightDownLastUpdate {
		chunk.AddBlockAt(p.ClosestEmptySpace.Vec3(), false, p.SelectedBlock)
		p.ClosestEmptySpace = nil
		p._mouseRightDownLastUpdate = true
	}
	if !p.ControlHandler.IsDown(int(glfw.MouseButtonRight)) {
		p._mouseRightDownLastUpdate = false
	}

	if p.ControlHandler.IsDown(int(glfw.Key1)) {
		p.SelectedBlock = block.GetBlockTypes()[0]
	}
	if p.ControlHandler.IsDown(int(glfw.Key2)) {
		p.SelectedBlock = block.GetBlockTypes()[1]
	}
	if p.ControlHandler.IsDown(int(glfw.Key3)) {
		p.SelectedBlock = block.GetBlockTypes()[2]
	}
	if p.ControlHandler.IsDown(int(glfw.Key4)) {
		p.SelectedBlock = block.GetBlockTypes()[3]
	}
	if p.ControlHandler.IsDown(int(glfw.Key5)) {
		p.SelectedBlock = block.GetBlockTypes()[4]
	}
	if p.ControlHandler.IsDown(int(glfw.Key6)) {
		p.SelectedBlock = block.GetBlockTypes()[5]
	}
	if p.ControlHandler.IsDown(int(glfw.Key7)) {
		p.SelectedBlock = block.GetBlockTypes()[6]
	}
	if p.ControlHandler.IsDown(int(glfw.Key8)) {
		p.SelectedBlock = block.GetBlockTypes()[7]
	}

	if p.ControlHandler.IsToggled(int(glfw.KeyLeftShift)) {
		p.WalkingSpeed = p.defaultSpeed * p.RunningMultiplier
	} else {
		p.WalkingSpeed = p.defaultSpeed
	}

	if p.ControlHandler.IsDown(int(glfw.Key1)) {
		color = mgl32.Vec3{1.0, 0.0, 0.0}
	}

	if p.ControlHandler.IsDown(int(glfw.Key2)) {
		color = mgl32.Vec3{1.0, 1.0, 0.0}
	}

	if p.ControlHandler.IsDown(int(glfw.Key3)) {
		color = mgl32.Vec3{0.0, 1.0, 0.0}
	}

	if p.ControlHandler.IsDown(int(glfw.Key4)) {
		color = mgl32.Vec3{0.0, 1.0, 1.0}
	}

	if p.ControlHandler.IsDown(int(glfw.Key4)) {
		color = mgl32.Vec3{0.0, 0.0, 1.0}
	}

	if p.ControlHandler.IsDown(int(glfw.Key4)) {
		color = mgl32.Vec3{1.0, 0.0, 1.0}
	}

	if collidesAbove {
		p._isJumping = false
	}

	if p.ControlHandler.IsDown(int(glfw.KeySpace)) && !p.IsJumping() && collidedBelow {
		p.Jump()
	}

	if p.IsJumping() {
		p.HandleJump()
	}

	p.HandleWorldLimits(world)

	p.Camera.Follow(p.Position.Add(mgl32.Vec4{0.0, float32(configs.PlayerHeight) / 2, 0.0, 0.0}))
	p.Camera.Update()

	p.LastChunk = chunk.ID
	p.HandleBlockInteractions(chunk)

}

func (p *Player) UpdateBoundingBox(newPosition mgl32.Vec4) collisions.CubeBoundingBox {
	r := 0.3
	Ax := r*math.Cos(p.PlayerTheta-math.Pi/4) + float64(newPosition.X())
	Az := r*math.Sin(-p.PlayerTheta+math.Pi/4) + float64(newPosition.Z())
	Bx := r*math.Cos(p.PlayerTheta-math.Pi*3/4) + float64(newPosition.X())
	Bz := r*math.Sin(-p.PlayerTheta+math.Pi*3/4) + float64(newPosition.Z())
	Cx := r*math.Cos(p.PlayerTheta-math.Pi*5/4) + float64(newPosition.X())
	Cz := r*math.Sin(-p.PlayerTheta+math.Pi*5/4) + float64(newPosition.Z())
	Dx := r*math.Cos(p.PlayerTheta-math.Pi*7/4) + float64(newPosition.X())
	Dz := r*math.Sin(-p.PlayerTheta+math.Pi*7/4) + float64(newPosition.Z())

	minx := math.Min(Ax, math.Min(Bx, math.Min(Cx, Dx)))
	minz := math.Min(Az, math.Min(Bz, math.Min(Cz, Dz)))
	maxx := math.Max(Ax, math.Max(Bx, math.Max(Cx, Dx)))
	maxz := math.Max(Az, math.Max(Bz, math.Max(Cz, Dz)))
	miny := newPosition.Y() - 1.0
	maxy := newPosition.Y() + 0.8

	p.BoundingBoxFutureVertices[0] = mgl32.Vec3{float32(Ax), miny, float32(Az)}
	p.BoundingBoxFutureVertices[1] = mgl32.Vec3{float32(Bx), miny, float32(Bz)}
	p.BoundingBoxFutureVertices[2] = mgl32.Vec3{float32(Cx), miny, float32(Cz)}
	p.BoundingBoxFutureVertices[3] = mgl32.Vec3{float32(Dx), miny, float32(Dz)}
	p.BoundingBoxFutureVertices[4] = mgl32.Vec3{float32(Ax), maxy, float32(Az)}
	p.BoundingBoxFutureVertices[5] = mgl32.Vec3{float32(Bx), maxy, float32(Bz)}
	p.BoundingBoxFutureVertices[6] = mgl32.Vec3{float32(Cx), maxy, float32(Cz)}
	p.BoundingBoxFutureVertices[7] = mgl32.Vec3{float32(Dx), maxy, float32(Dz)}

	bb := collisions.CubeBoundingBox{
		Mins:  mgl32.Vec3{float32(minx), miny, float32(minz)},
		Maxes: mgl32.Vec3{float32(maxx), maxy, float32(maxz)},
	}

	return bb
}

func (p *Player) HandleBlockInteractions(chunk *chunk.Chunk) {
	lookingAtPoint := p.Position.Add(mgl32.Vec4{0.0, float32(configs.PlayerHeight) / 2, 0.0, 1.0}).Add(p.Camera.ViewVector)
	lookingDirection := lookingAtPoint.Sub(p.Position.Add(mgl32.Vec4{0.0, float32(configs.PlayerHeight) / 2, 0.0, 1.0}))

	px, py, pz := p.GetRoundedPosition()
	shouldBreak := false

	// bounding box
	for s := 0.0; s < 10.0; s += 0.01 {
		if shouldBreak {
			break
		}
		ray := lookingDirection.Mul(float32(s))
		ray = mgl32.Vec4{ray.X() + p.Position.X(), ray.Y() + p.Position.Y() + float32(configs.PlayerHeight)/2, ray.Z() + p.Position.Z(), 0.0}
		for x := px - 2; x <= px+2; x++ {
			if shouldBreak {
				break
			}
			for y := py - 2; y <= py+2; y++ {
				if shouldBreak {
					break
				}
				for z := pz - 2; z <= pz+2; z++ {
					if chunk.GetBlockAt(x, y, z) == nil {
						continue
					}

					highestX := float32(x) + (float32(configs.BlockSize))/2
					highestY := float32(y) + (float32(configs.BlockSize))/2
					highestZ := float32(z) + (float32(configs.BlockSize))/2
					lowestX := float32(x) - (float32(configs.BlockSize))/2
					lowestY := float32(y) - (float32(configs.BlockSize))/2
					lowestZ := float32(z) - (float32(configs.BlockSize))/2

					if ray.X() <= highestX && ray.X() >= lowestX &&
						ray.Y() <= highestY && ray.Y() >= lowestY &&
						ray.Z() <= highestZ && ray.Z() >= lowestZ {
						p.HitAt = &mgl32.Vec4{float32(x), float32(y), float32(z), 1.0}
						//chunk[x][y][z].WithEdges = true
						shouldBreak = true
						p.ClosestEmptySpace = chunk.FindPlacementPosition(*p.HitAt, ray, mgl32.Vec3{highestX, highestY, highestZ},
							mgl32.Vec3{lowestX, lowestY, lowestZ})
						break
					}
				}
			}
		}
	}
}

func (p *Player) HandleWorldLimits(world *world.World) {
	worldSizeX, worldSizeZ := int(world.Size.X()), int(world.Size.Z())
	roundedPlayerX := int(math.Round(float64(p.Position.X())))

	roundedPlayerZ := int(math.Round(float64(p.Position.Z())))
	if roundedPlayerX < -worldSizeX {
		roundedPlayerX = -worldSizeX
		p.SetPosition(mgl32.Vec4{-float32(worldSizeX), p.Position.Y(), p.Position.Z(), 1.0})
	}

	if roundedPlayerX > worldSizeX-1 {
		roundedPlayerX = worldSizeX - 1
		p.SetPosition(mgl32.Vec4{float32(worldSizeX) - 1, p.Position.Y(), p.Position.Z(), 1.0})
	}

	if roundedPlayerZ < -worldSizeZ {
		roundedPlayerZ = -worldSizeZ
		p.SetPosition(mgl32.Vec4{p.Position.X(), p.Position.Y(), -float32(worldSizeZ), 1.0})
	}

	if roundedPlayerZ > worldSizeZ-1 {
		roundedPlayerZ = worldSizeZ - 1
		p.SetPosition(mgl32.Vec4{p.Position.X(), p.Position.Y(), float32(worldSizeZ) - 1, 1.0})
	}
}

func (p Player) GetRoundedPosition() (int, int, int) {
	roundedX := int(math.Round(float64(p.Position.X())))
	roundedY := int(math.Round(float64(p.Position.Y())))
	roundedZ := int(math.Round(float64(p.Position.Z())))

	return roundedX, roundedY, roundedZ
}

func (p Player) GetChunkOffset() mgl32.Vec2 {
	rx, _, rz := p.GetRoundedPosition()
	return mgl32.Vec2{float32(math.Floor(float64(float32(rx) / float32(configs.ChunkSize)))), float32(math.Floor(float64(float32(rz) / float32(configs.ChunkSize))))}
}

func (p Player) GetRealPosition() (float32, float32, float32) {

	return p.Position.X(), p.Position.Y(), p.Position.Z()
}

func (p Player) GetFlooredPosition() (int, int, int) {
	roundedX := int(math.Floor(float64(p.Position.X())))
	roundedY := int(math.Floor(float64(p.Position.Y())))
	roundedZ := int(math.Floor(float64(p.Position.Z())))

	return roundedX, roundedY, roundedZ
}
