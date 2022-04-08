package player

import (
	"fmt"
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/camera"
	"github.com/reonardoleis/fcg-glcraft/collisions"
	"github.com/reonardoleis/fcg-glcraft/configs"

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
	// fmt.Println(vx, vy, vz)
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

func (p *Player) CheckCollisions(futureBoundingBox collisions.CubeBoundingBox, blocks world.WorldBlocks) (collidesSides, collidesBelow, collidesAbove bool) {
	//fmt.Println(futureBoundingBox)
	playerX, playerY, playerZ := p.GetRoundedPosition()
	//fmt.Println(playerX, playerY, playerZ)
	for x := playerX - 1; x <= playerX+1; x++ {
		for y := playerY - 1; y <= playerY+1; y++ {
			for z := playerZ - 1; z <= playerZ+1; z++ {
				if blocks[x][y][z] == nil {
					//fmt.Println("NULO..")
					continue
				}
				if y == playerY-1 && z != playerZ && x != playerX {
					continue
				}

				//fmt.Println("NAO NULO..")
				cubeVertices := blocks[x][y][z].GetFutureVertices()
				blockBoundingBox := collisions.NewCubeBoundingBox(blocks[x][y][z].Position.Vec3(), float32(configs.BlockSize), float32(configs.BlockSize))

				//fmt.Println("BLOCO Max: ", blockBoundingBox.Maxes)
				//fmt.Println("BLOCO Min: ", blockBoundingBox.Mins)
				//fmt.Println(x, playerX, y, playerY, z, playerZ)

				if y == playerY-1 && p.Collider.Collides(futureBoundingBox, *blockBoundingBox, p.BoundingBoxFutureVertices, cubeVertices) {
					fmt.Println("Colidindo abaixo...")
					collidesBelow = true
					continue
				}

				if p.Collider.Collides(futureBoundingBox, *blockBoundingBox, p.BoundingBoxFutureVertices, cubeVertices) {

					blocks[x][y][z].WithEdges = true

					if x == playerX && z == playerZ && y == playerY+1 {
						collidesAbove = true
					} else {
						collidesSides = true
						fmt.Println("Colidindo dos lados...")
						//fmt.Println("Bloco: ", x, y, z, "Player: ", playerX, playerY, playerZ)

					}
				}

			}

		}
	}
	return
}

func (p *Player) Update(world *world.World) {
	p.HandleLookDirection()
	/*
		Ax := 5*math.Cos(p.PlayerTheta-math.Pi/4) + float64(p.Position.X())
		Az := 5*math.Sin(-p.PlayerTheta+math.Pi/4) + float64(p.Position.Z())
		Bx := 5*math.Cos(p.PlayerTheta-math.Pi*3/4) + float64(p.Position.X())
		Bz := 5*math.Sin(-p.PlayerTheta+math.Pi*3/4) + float64(p.Position.Z())
		Cx := 5*math.Cos(p.PlayerTheta-math.Pi*5/4) + float64(p.Position.X())
		Cz := 5*math.Sin(-p.PlayerTheta+math.Pi*5/4) + float64(p.Position.Z())
		Dx := 5*math.Cos(p.PlayerTheta-math.Pi*7/4) + float64(p.Position.X())
		Dz := 5*math.Sin(-p.PlayerTheta+math.Pi*7/4) + float64(p.Position.Z())
		AxY := 5*math.Cos(p.PlayerTheta-math.Pi/4) + float64(p.Position.X())
		AzY := 5*math.Sin(-p.PlayerTheta+math.Pi/4) + float64(p.Position.Z())
		BxY := 5*math.Cos(p.PlayerTheta-math.Pi*3/4) + float64(p.Position.X())
		BzY := 5*math.Sin(-p.PlayerTheta+math.Pi*3/4) + float64(p.Position.Z())
		CxY := 5*math.Cos(p.PlayerTheta-math.Pi*5/4) + float64(p.Position.X())
		CzY := 5*math.Sin(-p.PlayerTheta+math.Pi*5/4) + float64(p.Position.Z())
		DxY := 5*math.Cos(p.PlayerTheta-math.Pi*7/4) + float64(p.Position.X())
		DzY := 5*math.Sin(-p.PlayerTheta+math.Pi*7/4) + float64(p.Position.Z())
		cube := game_objects.NewBlock(float32(Ax), p.Position.Y(), float32(Az), 1, false, false, game_objects.BlockDirt)
		cube2 := game_objects.NewBlock(float32(Bx), p.Position.Y(), float32(Bz), 1, false, false, game_objects.BlockGrass)
		cube3 := game_objects.NewBlock(float32(Cx), p.Position.Y(), float32(Cz), 1, false, false, game_objects.BlockWater)
		cube4 := game_objects.NewBlock(float32(Dx), p.Position.Y(), float32(Dz), 1, false, false, game_objects.BlockWood)
		cube5 := game_objects.NewBlock(float32(AxY), p.Position.Y(), float32(AzY), 1, false, false, game_objects.BlockDirt)
		cube6 := game_objects.NewBlock(float32(BxY), p.Position.Y(), float32(BzY), 1, false, false, game_objects.BlockGrass)
		cube7 := game_objects.NewBlock(float32(CxY), p.Position.Y(), float32(CzY), 1, false, false, game_objects.BlockWater)
		cube8 := game_objects.NewBlock(float32(DxY), p.Position.Y(), float32(DzY), 1, false, false, game_objects.BlockWood)
		minx := math.Min(Ax, math.Min(Bx, math.Min(Cx, Dx)))
		minz := math.Min(Az,math.Min(Bz,math.Min(Cz,Dz)))
		maxx := math.Max(Ax,math.Max(Bx,math.Max(Cx,Dx)))
		maxz := math.Max(Az,math.Max(Bz,math.Max(Cz,Dz)))
		miny := p.Position.Y()-1.0
		maxy := p.Position.Y()+1.0
		bb := collisions.CubeBoundingBox{
			Maxes: mgl32.Vec3{float32(minx),miny,float32(minz)},
			Mins: mgl32.Vec3{float32(maxx),maxy,float32(maxz)},
		}

		cube.Draw2()
		cube2.Draw2()
		cube3.Draw2()
		cube4.Draw2()
		cube5.Draw2()
		cube6.Draw2()
		cube7.Draw2()
		cube8.Draw2()
	*/
	//fmt.Println(p.PlayerTheta)
	//fmt.Println("A: ", Ax, Az)
	//fmt.Println("B: ", Bx, Bz)
	//fmt.Println("C: ", Cx, Az)
	//fmt.Println("D: ", Dx, Dz)

	w, u := p.GetMovementVector()

	newPosition := p.Position
	wneg := w.Mul(-1)
	if p.ControlHandler.IsDown(int(glfw.KeyW)) {

		//fmt.Println(wneg.X()+p.Position.X(), wneg.X()+p.Position.Z())
		//fmt.Println(p.Position.X()+1, p.Position.Z()+1, p.Position.X()-1, p.Position.Z()-1)

		if wneg.X()+p.Position.X() >= p.Position.X() && wneg.Z()+p.Position.Z() >= p.Position.Z() {

			if world.Blocks[int(p.Position.X())+1][int(p.Position.Y())][int(p.Position.Z())] != nil &&
				world.Blocks[int(p.Position.X())][int(p.Position.Y())][int(p.Position.Z())+1] != nil {
				fmt.Println(wneg)
				fmt.Println(p.Position)
				fmt.Println("Primeiro if... #1")

			} else if world.Blocks[int(p.Position.X())+1][int(p.Position.Y())][int(p.Position.Z())] == nil &&
				world.Blocks[int(p.Position.X())][int(p.Position.Y())][int(p.Position.Z())+1] != nil {
				fmt.Println("Segundo if... #1")
				wneg = wneg.Mul(p.WalkingSpeed * float32(math2.DeltaTime))
				newPosition = p.Position.Add(mgl32.Vec4{wneg.X(), 0.0, 0.0, 0.0})
			} else {
				fmt.Println("Terceiro if... #1")
				wneg = wneg.Mul(p.WalkingSpeed * float32(math2.DeltaTime))
				newPosition = p.Position.Add(mgl32.Vec4{wneg.X(), 0.0, wneg.Z(), 0.0})
			}
		} else if wneg.X()+p.Position.X() >= p.Position.X() && wneg.Z()+p.Position.Z() <= p.Position.Z() {

			if world.Blocks[int(p.Position.X())+1][int(p.Position.Y())][int(p.Position.Z())] != nil &&
				world.Blocks[int(p.Position.X())][int(p.Position.Y())][int(p.Position.Z())-1] != nil {
				fmt.Println(wneg)
				fmt.Println(p.Position)
				fmt.Println("Primeiro if... #2")

			} else if world.Blocks[int(p.Position.X())+1][int(p.Position.Y())][int(p.Position.Z())] == nil &&
				world.Blocks[int(p.Position.X())][int(p.Position.Y())][int(p.Position.Z())-1] != nil {
				fmt.Println("Segundo if...#2")
				wneg = wneg.Mul(p.WalkingSpeed * float32(math2.DeltaTime))
				newPosition = p.Position.Add(mgl32.Vec4{wneg.X(), 0.0, 0.0, 0.0})
			} else {
				fmt.Println("Terceiro if...#2")
				wneg = wneg.Mul(p.WalkingSpeed * float32(math2.DeltaTime))
				newPosition = p.Position.Add(mgl32.Vec4{wneg.X(), 0.0, wneg.Z(), 0.0})
			}
		} else if wneg.X()+p.Position.X() <= p.Position.X() && wneg.Z()+p.Position.Z() <= p.Position.Z() {
			if world.Blocks[int(p.Position.X())-1][int(p.Position.Y())][int(p.Position.Z())] != nil &&
				world.Blocks[int(p.Position.X())][int(p.Position.Y())][int(p.Position.Z())-1] != nil {
				fmt.Println(wneg)
				fmt.Println(p.Position)
				fmt.Println("Primeiro if...#3")

			} else if world.Blocks[int(p.Position.X())-1][int(p.Position.Y())][int(p.Position.Z())] == nil &&
				world.Blocks[int(p.Position.X())][int(p.Position.Y())][int(p.Position.Z())-1] != nil {
				fmt.Println("Segundo if...#3")
				wneg = wneg.Mul(p.WalkingSpeed * float32(math2.DeltaTime))
				newPosition = p.Position.Add(mgl32.Vec4{wneg.X(), 0.0, 0.0, 0.0})
			} else {
				fmt.Println("Terceiro if...#3")
				wneg = wneg.Mul(p.WalkingSpeed * float32(math2.DeltaTime))
				newPosition = p.Position.Add(mgl32.Vec4{wneg.X(), 0.0, wneg.Z(), 0.0})
			}
		} else if wneg.X()+p.Position.X() <= p.Position.X() && wneg.Z()+p.Position.Z() >= p.Position.Z() {
			if world.Blocks[int(p.Position.X())-1][int(p.Position.Y())][int(p.Position.Z())] != nil &&
				world.Blocks[int(p.Position.X())][int(p.Position.Y())][int(p.Position.Z())+1] != nil {
				fmt.Println(wneg)
				fmt.Println(p.Position)
				fmt.Println("Primeiro if...#4")
			} else if world.Blocks[int(p.Position.X())-1][int(p.Position.Y())][int(p.Position.Z())] == nil &&
				world.Blocks[int(p.Position.X())][int(p.Position.Y())][int(p.Position.Z())+1] != nil {
				fmt.Println("Segundo if...#4")
				wneg = wneg.Mul(p.WalkingSpeed * float32(math2.DeltaTime))
				newPosition = p.Position.Add(mgl32.Vec4{wneg.X(), 0.0, 0.0, 0.0})
			} else {
				fmt.Println("Terceiro if...#4")
				wneg = wneg.Mul(p.WalkingSpeed * float32(math2.DeltaTime))
				newPosition = p.Position.Add(mgl32.Vec4{wneg.X(), 0.0, wneg.Z(), 0.0})
			}
		} else {

			newPosition = p.Position.Add(w.Mul(-1).Mul(p.WalkingSpeed * float32(math2.DeltaTime)))

		}

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
	if p.ControlHandler.IsDown(int(glfw.MouseButtonLeft)) && !p._mouseLeftDownLastUpdate && p.HitAt != nil {
		p._mouseLeftDownLastUpdate = true
		world.RemoveBlockFrom(*p.HitAt)
		p.HitAt = nil
	}
	if !p.ControlHandler.IsDown(int(glfw.MouseButtonLeft)) {
		p._mouseLeftDownLastUpdate = false
	}
	if p.ControlHandler.IsDown(int(glfw.MouseButtonRight)) && p.ClosestEmptySpace != nil && !p._mouseRightDownLastUpdate {
		world.AddBlockAt(p.ClosestEmptySpace.Vec3(), false, mgl32.Vec3{})
		p.ClosestEmptySpace = nil
		p._mouseRightDownLastUpdate = true
	}
	if !p.ControlHandler.IsDown(int(glfw.MouseButtonRight)) {
		p._mouseRightDownLastUpdate = false
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

	/*if (world.Blocks[int(newPosition.X())][int(newPosition.Y())][int(newPosition.Z())] != nil ||
		world.Blocks[int(newPosition.X())][int(newPosition.Y())-1][int(newPosition.Z())] != nil ||
		world.Blocks[int(newPosition.X())][int(newPosition.Y())+1][int(newPosition.Z())] != nil ||
		world.Blocks[int(newPosition.X()+1)][int(newPosition.Y())][int(newPosition.Z()+1)] != nil ||
		world.Blocks[int(newPosition.X()+1)][int(newPosition.Y())][int(newPosition.Z()-1)] != nil ||
		world.Blocks[int(newPosition.X()-1)][int(newPosition.Y())][int(newPosition.Z()-1)] != nil ||
		world.Blocks[int(newPosition.X()-1)][int(newPosition.Y())][int(newPosition.Z()+1)] != nil) &&
		(world.Blocks[int(newPosition.X())][int(newPosition.Y())][int(newPosition.Z())] != nil && world.Blocks[int(newPosition.X())][int(newPosition.Y())][int(newPosition.Z())].BlockType != game_objects.BlockWater) {
		newPosition = p.Position
	}*/

	if !p._isJumping {
		fmt.Println("Caindo....")
		newPosition = newPosition.Sub(mgl32.Vec4{0.0, math2.GravityAccel * float32(math2.DeltaTime), 0.0, 0.0})
	}
	Ax := 0.3*math.Cos(p.PlayerTheta-math.Pi/4) + float64(newPosition.X())
	Az := 0.3*math.Sin(-p.PlayerTheta+math.Pi/4) + float64(newPosition.Z())
	Bx := 0.3*math.Cos(p.PlayerTheta-math.Pi*3/4) + float64(newPosition.X())
	Bz := 0.3*math.Sin(-p.PlayerTheta+math.Pi*3/4) + float64(newPosition.Z())
	Cx := 0.3*math.Cos(p.PlayerTheta-math.Pi*5/4) + float64(newPosition.X())
	Cz := 0.3*math.Sin(-p.PlayerTheta+math.Pi*5/4) + float64(newPosition.Z())
	Dx := 0.3*math.Cos(p.PlayerTheta-math.Pi*7/4) + float64(newPosition.X())
	Dz := 0.3*math.Sin(-p.PlayerTheta+math.Pi*7/4) + float64(newPosition.Z())

	minx := math.Min(Ax, math.Min(Bx, math.Min(Cx, Dx)))
	minz := math.Min(Az, math.Min(Bz, math.Min(Cz, Dz)))
	maxx := math.Max(Ax, math.Max(Bx, math.Max(Cx, Dx)))
	maxz := math.Max(Az, math.Max(Bz, math.Max(Cz, Dz)))
	miny := newPosition.Y() - 1.0
	maxy := newPosition.Y() + 1.0

	p.BoundingBoxFutureVertices[0] = mgl32.Vec3{float32(Ax), miny, float32(Az)}
	p.BoundingBoxFutureVertices[1] = mgl32.Vec3{float32(Bx), miny, float32(Bz)}
	p.BoundingBoxFutureVertices[2] = mgl32.Vec3{float32(Cx), miny, float32(Cz)}
	p.BoundingBoxFutureVertices[3] = mgl32.Vec3{float32(Dx), miny, float32(Dz)}
	p.BoundingBoxFutureVertices[4] = mgl32.Vec3{float32(Ax), maxy, float32(Az)}
	p.BoundingBoxFutureVertices[5] = mgl32.Vec3{float32(Bx), maxy, float32(Bz)}
	p.BoundingBoxFutureVertices[6] = mgl32.Vec3{float32(Cx), maxy, float32(Cz)}
	p.BoundingBoxFutureVertices[7] = mgl32.Vec3{float32(Dx), maxy, float32(Dz)}

	bb := collisions.CubeBoundingBox{
		Maxes: mgl32.Vec3{float32(minx), miny, float32(minz)},
		Mins:  mgl32.Vec3{float32(maxx), maxy, float32(maxz)},
	}
	//fmt.Println("PLAYER MAX: ", maxx, maxy, maxz)
	//fmt.Println("PLAYER MIN: ", minx, miny, minz)

	collided, collidedBelow, collidesAbove := p.CheckCollisions(bb, world.Blocks)

	if collidedBelow {
		newPosition = mgl32.Vec4{newPosition.X(), p.Position.Y(), newPosition.Z(), 1.0}
		collided, _, _ = p.CheckCollisions(bb, world.Blocks)
	}
	if !collided || true {
		p.Position = newPosition

	}

	p.BoundingBox.UpdateBounds(p.Position.Vec3())

	if collidesAbove {
		p._isJumping = false
	}

	if p.IsJumping() {
		p.HandleJump()
	}
	if p.ControlHandler.IsDown(int(glfw.KeySpace)) && !p.IsJumping() && collidedBelow {
		p.Jump()
	}

	p.HandleWorldLimits(world)

	p.Camera.Follow(p.Position.Add(mgl32.Vec4{0.0, float32(configs.BlockSize), 0.0, 0.0}))
	p.Camera.Update()

	p.HandleBlockInteractions(world)

}

func (p *Player) HandleBlockInteractions(world *world.World) {
	lookingAtPoint := p.Position.Add(mgl32.Vec4{0.0, float32(configs.BlockSize), 0.0, 1.0}).Add(p.Camera.ViewVector)
	lookingDirection := lookingAtPoint.Sub(p.Position.Add(mgl32.Vec4{0.0, float32(configs.BlockSize), 0.0, 1.0}))

	px, py, pz := p.GetRoundedPosition()
	shouldBreak := false

	// bounding box
	for s := 0.0; s < 5.0; s += 0.5 {
		if shouldBreak {
			break
		}
		ray := lookingDirection.Mul(float32(s))
		ray = mgl32.Vec4{ray.X() + p.Position.X(), ray.Y() + p.Position.Y() + float32(configs.BlockSize), ray.Z() + p.Position.Z(), 0.0}
		for x := px - 2; x <= px+2; x++ {
			if shouldBreak {
				break
			}
			for y := py - 2; y <= py+2; y++ {
				if shouldBreak {
					break
				}
				for z := pz - 2; z <= pz+2; z++ {
					if world.Blocks[x][y][z] == nil {
						continue
					}

					highestX := float32(x) + (world.Blocks[x][y][z].Size)/2
					highestY := float32(y) + (world.Blocks[x][y][z].Size)/2
					highestZ := float32(z) + (world.Blocks[x][y][z].Size)/2
					lowestX := float32(x) - (world.Blocks[x][y][z].Size)/2
					lowestY := float32(y) - (world.Blocks[x][y][z].Size)/2
					lowestZ := float32(z) - (world.Blocks[x][y][z].Size)/2

					if ray.X() <= highestX && ray.X() >= lowestX &&
						ray.Y() <= highestY && ray.Y() >= lowestY &&
						ray.Z() <= highestZ && ray.Z() >= lowestZ {
						p.HitAt = &mgl32.Vec4{float32(x), float32(y), float32(z), 1.0}
						//world.Blocks[x][y][z].WithEdges = true
						shouldBreak = true
						p.ClosestEmptySpace = world.FindPlacementPosition(*p.HitAt, ray, mgl32.Vec3{highestX, highestY, highestZ},
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

func (p Player) GetFlooredPosition() (int, int, int) {
	roundedX := int(math.Floor(float64(p.Position.X())))
	roundedY := int(math.Floor(float64(p.Position.Y())))
	roundedZ := int(math.Floor(float64(p.Position.Z())))

	return roundedX, roundedY, roundedZ
}
