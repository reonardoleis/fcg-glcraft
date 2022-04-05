package scene

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/camera"
	"github.com/reonardoleis/fcg-glcraft/engine/controls"
	"github.com/reonardoleis/fcg-glcraft/engine/shaders"
	"github.com/reonardoleis/fcg-glcraft/geometry"
	math2 "github.com/reonardoleis/fcg-glcraft/math"
	"github.com/reonardoleis/fcg-glcraft/player"
	"github.com/reonardoleis/fcg-glcraft/world"
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
}

func NewScene(world *world.World, mainCamera *camera.Camera, player *player.Player, controlHandler controls.Controls, sceneType SceneType) *Scene {
	return &Scene{
		World:          world,
		MainCamera:     mainCamera,
		Player:         player,
		ControlHandler: &controlHandler,
		Type:           sceneType,
	}
}

func (s *Scene) Update(window glfw.Window) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	gl.UseProgram(shaders.ShaderProgramCrosshair)

	geometry.DrawCrosshair()

	gl.UseProgram(shaders.ShaderProgramDefault)
	s.MainCamera.Update()
	s.Player.Update(s.World)

	roundedPlayerX, roundedPlayerY, roundedPlayerZ := s.Player.GetRoundedPosition()
	playerY := float64(s.Player.Position.Y())
	s.World.Update(mgl32.Vec3{float32(roundedPlayerX), float32(roundedPlayerY), float32(roundedPlayerZ)})

	window.SetTitle(fmt.Sprintf("FPS: %v - X: %v - Y: %v - Z: %v - wsX: %v - wsZ: %v", 1/math2.DeltaTime,
		roundedPlayerX, playerY, roundedPlayerZ, s.World.Size.X(), s.World.Size.Z()))

	blockBelow := s.World.FindHighestBlock(roundedPlayerX, roundedPlayerZ)

	s.Player.Fall(blockBelow)

	s.ControlHandler.FinishMousePositionChanged()
	window.SwapBuffers()
	glfw.PollEvents()
}
