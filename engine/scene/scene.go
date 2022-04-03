package scene

import (
	"github.com/reonardoleis/fcg-glcraft/camera"
	"github.com/reonardoleis/fcg-glcraft/game_objects"
	"github.com/reonardoleis/fcg-glcraft/player"
	"github.com/reonardoleis/fcg-glcraft/world"
)

type SceneType uint8

const (
	MenuScene uint32 = iota
	GameScene
	OtherScene
)

type Scene struct {
	Name        string
	Type        SceneType
	GameObjects []*game_objects.GameObject
	World       *world.World
	MainCamera  *camera.Camera
	Player      *player.Player
}

func NewScene() *Scene {
	return &Scene{
		GameObjects: []*game_objects.GameObject{},
	}
}

func (vs *Scene) Add(gameObject game_objects.GameObject) {
	vs.GameObjects = append(vs.GameObjects, &gameObject)
}

func (vs Scene) GetGameObjects() []*game_objects.GameObject {
	return vs.GameObjects
}
