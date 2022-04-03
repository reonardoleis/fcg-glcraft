package world

import (
	"fmt"

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
		} else if w.Blocks[wx][keys[wy]][wz] != nil && w.Blocks[wx][keys[wy]][wz].Position.Y() > highestBlock.Position.Y() {
			highestBlock = w.Blocks[wx][keys[wy]][wz]
		}
	}

	return highestBlock
}

func (w *World) AddBlockAt(position mgl32.Vec3, ephemeral bool) {
	fmt.Println(position)
	x, y, z := position.Elem()
	if w.Blocks[int(x)] == nil {
		w.Blocks[int(x)] = make(map[int]map[int]*game_objects.GameObject)
	}

	if w.Blocks[int(x)][int(y)] == nil {
		w.Blocks[int(x)][int(y)] = make(map[int]*game_objects.GameObject)
	}

	newCube := game_objects.NewBlock(float32(x), float32(y), float32(z), 1, true, ephemeral)
	newCube.WithEdges = false
	w.Blocks[int(x)][int(y)][int(z)] = &newCube

	if y > w.Size.Y() {
		w.Size = mgl32.Vec3{w.Size.X(), y, w.Size.Z()}
	}
}

func (w World) FindPlacementPosition(hitAt mgl32.Vec4, nearFrom mgl32.Vec4) *mgl32.Vec4 {
	block1 := mgl32.Vec4{hitAt.X(), hitAt.Y() + 1, hitAt.Z(), 1.0}     // above
	block2 := mgl32.Vec4{hitAt.X(), hitAt.Y() - 1, hitAt.Z(), 1.0}     // below
	block3 := mgl32.Vec4{hitAt.X() + 1, hitAt.Y(), hitAt.Z(), 1.0}     // north
	block4 := mgl32.Vec4{hitAt.X() - 1, hitAt.Y() + 1, hitAt.Z(), 1.0} // south
	block5 := mgl32.Vec4{hitAt.X(), hitAt.Y() + 1, hitAt.Z() + 1, 1.0} // east
	block6 := mgl32.Vec4{hitAt.X(), hitAt.Y() + 1, hitAt.Z() - 1, 1.0} // west

	possibilities := []*mgl32.Vec4{&block1, &block2, &block3, &block4, &block5, &block6}
	best := possibilities[0]
	for i := 1; i < len(possibilities); i++ {
		if math2.Distance(*possibilities[i], nearFrom) < math2.Distance(*best, nearFrom) {
			best = possibilities[i]
		}
	}

	if w.Blocks[int(hitAt.X())][int(hitAt.Y())][int(hitAt.Z())] == nil {
		return nil
	}

	if w.Blocks[int(best.X())][int(best.Y())][int(best.Z())] != nil {
		return nil
	}

	return best
}

func (w *World) GenerateWorld() {
	math2.SetSeed(w.Seed)

	generatedWorld := make(WorldBlocks)
	for x := int(-w.Size.X()); x < int(w.Size.X()); x++ {
		generatedWorld[x] = make(map[int]map[int]*game_objects.GameObject)
		for y := 0; y < int(w.Size.Y()); y++ {
			generatedWorld[x][y] = make(map[int]*game_objects.GameObject)
			for z := int(-w.Size.Z()); z < int(w.Size.Z()); z++ {
				if z == 2 {
					continue
				}

				newCube := game_objects.NewBlock(float32(x), float32(y), float32(z), 1, true, false)
				newCube.WithEdges = false
				generatedWorld[x][y][z] = &newCube
			}
		}
	}

	w.Blocks = generatedWorld
}
