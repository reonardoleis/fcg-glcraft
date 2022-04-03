package world

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/game_objects"
	math2 "github.com/reonardoleis/fcg-glcraft/math"
)

type WorldBlocks = map[int]map[int]map[int]*game_objects.Block

type World struct {
	Name   string
	Size   mgl32.Vec3
	Blocks WorldBlocks
	Seed   int64
	Time   int64
}

func NewWorld(worldName string, size mgl32.Vec3, seed int64) *World {
	return &World{
		Name: worldName,
		Size: size,
		Seed: seed,
		Time: 0,
	}
}

func (w World) FindHighestBlock(wx, wz int) *game_objects.GameObject {
	var highestBlock *game_objects.GameObject
	keys := []int{}

	for k := range w.Blocks[wx] {
		keys = append(keys, k)
	}

	for wy := 0; wy < len(keys); wy++ {
		if highestBlock == nil {
			highestBlock = w.Blocks[wx][keys[wy]][wz]
		} else if w.Blocks[wx][keys[wy]][wz].Position.Y() > highestBlock.Position.Y() {
			highestBlock = w.Blocks[wx][keys[wy]][wz]
		}
	}

	return highestBlock
}

func (w *World) GenerateWorld() {
	math2.SetSeed(w.Seed)

	generatedWorld := make(WorldBlocks)
	for x := int(-w.Size.X()); x < int(w.Size.X()); x++ {
		generatedWorld[x] = make(map[int]map[int]*game_objects.GameObject)
		for y := int(-w.Size.Y()); y < int(w.Size.Y()); y++ {
			randY := math2.RandInt(0, 1)
			generatedWorld[x][y] = make(map[int]*game_objects.GameObject)
			for z := int(-w.Size.Z()); z < int(w.Size.Z()); z++ {
				if z == 2 {
					continue
				}

				newCube := game_objects.NewBlock(float32(x), float32(randY), float32(z), 1, true)
				generatedWorld[x][y][z] = &newCube
			}
		}
	}

	w.Blocks = generatedWorld
}
