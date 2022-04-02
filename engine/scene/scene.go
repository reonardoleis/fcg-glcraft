package scene

import "github.com/reonardoleis/fcg-glcraft/game_objects"

type Scene struct {
	GameObjects []*game_objects.GameObject
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
