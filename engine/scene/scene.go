package scene

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/block"
	"github.com/reonardoleis/fcg-glcraft/camera"
	"github.com/reonardoleis/fcg-glcraft/collisions"
	"github.com/reonardoleis/fcg-glcraft/engine/controls"
	"github.com/reonardoleis/fcg-glcraft/engine/shaders"
	"github.com/reonardoleis/fcg-glcraft/geometry"
	math2 "github.com/reonardoleis/fcg-glcraft/math"
	"github.com/reonardoleis/fcg-glcraft/player"
	"github.com/reonardoleis/fcg-glcraft/world"
	"github.com/reonardoleis/fcg-glcraft/world/chunk"
)

type SceneType = uint

const (
	GameScene = iota
	UIScene
)

type SceneManager struct {
	Scenes      []*Scene
	ActiveScene int
}

func NewSceneManager() SceneManager {
	return SceneManager{
		Scenes:      []*Scene{},
		ActiveScene: -1,
	}
}

func (sm *SceneManager) AddScene(scene *Scene) {
	sm.Scenes = append(sm.Scenes, scene)
}

func (sm *SceneManager) SetActiveScene(sceneIndex int) {
	sm.ActiveScene = sceneIndex
}

func (sm *SceneManager) HandleActiveScene(window glfw.Window) {
	sm.Scenes[sm.ActiveScene].Update(window)
}

type Scene struct {
	World          *world.World
	MainCamera     *camera.Camera
	Type           SceneType
	Player         *player.Player
	ControlHandler *controls.Controls
	Objs           []*geometry.GeometryInformation
}

func NewScene(world *world.World, mainCamera *camera.Camera, player *player.Player, controlHandler controls.Controls, sceneType SceneType, objs []*geometry.GeometryInformation) *Scene {
	return &Scene{
		World:          world,
		MainCamera:     mainCamera,
		Player:         player,
		ControlHandler: &controlHandler,
		Type:           sceneType,
		Objs:           objs,
	}
}

// Updates the scene each frame, will handle world and player updates
func (s *Scene) Update(window glfw.Window) {

	cx, cz := s.Player.GetChunkOffset().Elem()

	currentChunk := s.World.Chunks[int(cx)][int(cz)]

	if currentChunk.GetBlockInformationAt(int(s.Player.Position.X()), int(s.Player.Position.Y()), int(s.Player.Position.Z())) == chunk.BlockInformationCave {
		gl.ClearColor(0.47, 0.65, 1.0, 1.0)
	} else {
		gl.ClearColor(0.47, 0.65, 1.0, 1.0)
	}

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(shaders.ShaderProgramCrosshair)

	geometry.DrawCrosshair()

	gl.UseProgram(shaders.ShaderProgramDefault)

	// verify if player changed chunk
	if currentChunk.ID != s.Player.LastChunk {
		s.World.FutureChunks = s.World.Chunks
		s.World.HandleChunkChange(int(currentChunk.Offset[0]), int(currentChunk.Offset[1]))
		s.World.SetPopulatedBlocks(currentChunk.Offset[0], currentChunk.Offset[1])
	}

	roundedPlayerX, roundedPlayerY, roundedPlayerZ := s.Player.GetRoundedPosition()
	realPlayerX, realPlayerY, realPlayerZ := s.Player.GetRealPosition()
	//playerY := float64(s.Player.Position.Y())

	// handle wireframe mode
	if s.ControlHandler.IsToggled(int(glfw.KeyZ)) {
		block.BlockEdgesOnly = true
	} else {
		block.BlockEdgesOnly = false
	}

	backOfPlayer, frontOfPlayer := s.Player.GetFrontAndBackDirections()
	gl.BindVertexArray(1)

	// draw all .obj objects
	for _, obj := range s.Objs {

		if obj.Animating {
			if obj.T >= 3 {
				obj.Tdir = -1
			}
			if obj.T <= 0 {
				if obj.BCurve.ControlPoints != nil {
					c := obj.BCurve.T(obj.T)
					obj.Position = mgl32.Vec3{obj.Position[0] + c.X(), obj.Position[1], obj.Position[2] + c.Z()}
				}
				obj.Tdir = 1
				curve := math2.NewBezierCurve()
				curve.GenerateRandomPoints()
				obj.BCurve = curve
			}

			obj.T = obj.T + (obj.Tdir * 1 * float32(math2.DeltaTime))
			c := obj.BCurve.T(obj.T)

			obj.DrawAt(nil, 1, c)
		} else {
			obj.Draw(nil, 1)
		}

		sphereCollider := collisions.SphereCollider{
			Center: obj.Position,
			Radius: 5,
		}

		if sphereCollider.CollidesWith(s.Player.Position.Vec3()) {
			obj.Animating = true
		} else {
			obj.Animating = false
		}

	}
	gl.BindVertexArray(0)
	s.World.Update(mgl32.Vec3{float32(roundedPlayerX), float32(roundedPlayerY), float32(roundedPlayerZ)}, backOfPlayer, frontOfPlayer, currentChunk)
	s.Player.Update(s.World, s.World.Chunks[int(cx)][int(cz)])

	/*window.SetTitle(fmt.Sprintf("FPS: %v - X: %v - Y: %v - Z: %v - wsX: %v - wsZ: %v", 1/math2.DeltaTime,
	roundedPlayerX, playerY, roundedPlayerZ, s.World.Size.X(), s.World.Size.Z()))*/
	gl.BindVertexArray(0)
	window.SetTitle(fmt.Sprintf("FPS: %v - X: %v - Y: %v - Z: %v - wsX: %v - wsZ: %v", 1/math2.DeltaTime,
		realPlayerX, realPlayerY, realPlayerZ, s.World.Size.X(), s.World.Size.Z()))

	s.ControlHandler.FinishMousePositionChanged()
	window.SwapBuffers()
	glfw.PollEvents()

	if s.World.ShouldUpdateChunks {
		s.World.Chunks = s.World.FutureChunks
		s.World.ShouldUpdateChunks = false
	}

	if s.World.NextPopulatedBlocksReady {
		s.World.PopulatedBlocks = s.World.NextPopulatedBlocks
	}

	go s.World.SetPopulatedBlocks(cx, cz)
}
