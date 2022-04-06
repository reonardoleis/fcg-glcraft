package world

import (
	"math"
	"time"

	"github.com/aquilax/go-perlin"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/configs"
	"github.com/reonardoleis/fcg-glcraft/game_objects"
	math2 "github.com/reonardoleis/fcg-glcraft/math"
)

type WorldBlocks = map[int]map[int]map[int]*game_objects.Block

type World struct {
	Name                        string
	Size                        mgl32.Vec3
	Blocks                      WorldBlocks
	PopulatedBlocks             [][]*game_objects.Block
	ShouldUpdatePopulatedBlocks bool
	Seed                        int64
	Time                        int64
	Tick                        uint
	GlobalNoise                 *perlin.Perlin
}

func NewWorld(worldName string, size mgl32.Vec3, seed int64) *World {

	return &World{
		Name:        worldName,
		Size:        size,
		Seed:        seed,
		Time:        0,
		GlobalNoise: perlin.NewPerlin(2.0, 0.5, 10, time.Now().Unix()),
	}
}

func (w World) GetBlockFrom(wx, wy, wz int, playerSize float32) *game_objects.Block {
	return w.Blocks[wx][wy-int(playerSize)][wz]
}

func (w World) FindHighestBlock(wx, wz int) *game_objects.Block {
	var highestBlock *game_objects.Block
	keys := []int{}

	for k := range w.Blocks[wx] {
		keys = append(keys, k)
	}

	for wy := 0; wy < len(keys); wy++ {
		if highestBlock == nil {
			highestBlock = w.Blocks[wx][keys[wy]][wz]
		} else if w.Blocks[wx][keys[wy]][wz] != nil && w.Blocks[wx][keys[wy]][wz].Position.Y() > highestBlock.Position.Y() && w.Blocks[wx][keys[wy]][wz].BlockType != game_objects.BlockWater {
			highestBlock = w.Blocks[wx][keys[wy]][wz]
		}
	}

	return highestBlock
}

func (w *World) AddBlockAt(position mgl32.Vec3, ephemeral bool, color mgl32.Vec3) {
	// fmt.Println(position)
	x, y, z := position.Elem()
	if w.Blocks[int(x)] == nil {
		w.Blocks[int(x)] = make(map[int]map[int]*game_objects.Block)
	}

	if w.Blocks[int(x)][int(y)] == nil {
		w.Blocks[int(x)][int(y)] = make(map[int]*game_objects.Block)
	}

	newCube := game_objects.NewBlock(float32(x), float32(y), float32(z), 1, true, ephemeral, game_objects.BlockWood)
	newCube.WithEdges = false
	w.Blocks[int(x)][int(y)][int(z)] = &newCube
	w.ShouldUpdatePopulatedBlocks = true
}

func (w *World) RemoveBlockFrom(position mgl32.Vec4) {
	northNeighbor := math2.North(position, float32(configs.BlockSize)*2)
	southNeighbor := math2.South(position, float32(configs.BlockSize)*2)
	eastNeighbor := math2.East(position, float32(configs.BlockSize)*2)
	westNeighbor := math2.West(position, float32(configs.BlockSize)*2)
	upperNeighbor := math2.Upper(position, float32(configs.BlockSize)*2)
	lowerNeighbor := math2.Lower(position, float32(configs.BlockSize)*2)

	if w.Blocks[int(northNeighbor.X())][int(northNeighbor.Y())][int(northNeighbor.Z())] != nil {
		w.Blocks[int(northNeighbor.X())][int(northNeighbor.Y())][int(northNeighbor.Z())].Neighbors[1] = false
	}
	if w.Blocks[int(southNeighbor.X())][int(southNeighbor.Y())][int(southNeighbor.Z())] != nil {
		w.Blocks[int(southNeighbor.X())][int(southNeighbor.Y())][int(southNeighbor.Z())].Neighbors[0] = false
	}
	if w.Blocks[int(eastNeighbor.X())][int(eastNeighbor.Y())][int(eastNeighbor.Z())] != nil {
		w.Blocks[int(eastNeighbor.X())][int(eastNeighbor.Y())][int(eastNeighbor.Z())].Neighbors[3] = false
	}
	if w.Blocks[int(westNeighbor.X())][int(westNeighbor.Y())][int(westNeighbor.Z())] != nil {
		w.Blocks[int(westNeighbor.X())][int(westNeighbor.Y())][int(westNeighbor.Z())].Neighbors[2] = false
	}
	if w.Blocks[int(upperNeighbor.X())][int(upperNeighbor.Y())][int(upperNeighbor.Z())] != nil {
		w.Blocks[int(upperNeighbor.X())][int(upperNeighbor.Y())][int(upperNeighbor.Z())].Neighbors[5] = false
	}
	if w.Blocks[int(lowerNeighbor.X())][int(lowerNeighbor.Y())][int(lowerNeighbor.Z())] != nil {
		w.Blocks[int(lowerNeighbor.X())][int(lowerNeighbor.Y())][int(lowerNeighbor.Z())].Neighbors[4] = false
	}

	w.Blocks[int(position.X())][int(position.Y())][int(position.Z())] = nil
	w.ShouldUpdatePopulatedBlocks = true
}

func (w World) FindPlacementPosition(hitAt mgl32.Vec4, nearFrom mgl32.Vec4, boundingBoxHighests, boundingBoxLowests mgl32.Vec3) *mgl32.Vec4 {
	block1 := mgl32.Vec4{hitAt.X(), hitAt.Y() + 1, hitAt.Z(), 1.0} // above
	block2 := mgl32.Vec4{hitAt.X(), hitAt.Y() - 1, hitAt.Z(), 1.0} // below
	block3 := mgl32.Vec4{hitAt.X() + 1, hitAt.Y(), hitAt.Z(), 1.0} // north
	block4 := mgl32.Vec4{hitAt.X() - 1, hitAt.Y(), hitAt.Z(), 1.0} // south
	block5 := mgl32.Vec4{hitAt.X(), hitAt.Y(), hitAt.Z() + 1, 1.0} // east
	block6 := mgl32.Vec4{hitAt.X(), hitAt.Y(), hitAt.Z() - 1, 1.0} // west
	// fmt.Println(hitAt, nearFrom)
	bb := []float32{boundingBoxHighests.X() - nearFrom.X(), boundingBoxHighests.Y() - nearFrom.Y(), boundingBoxHighests.Z() - nearFrom.Z(),
		nearFrom.X() - boundingBoxLowests.X(), nearFrom.Y() - boundingBoxLowests.Y(), nearFrom.Z() - boundingBoxLowests.Z()}

	lowest := bb[0]
	lowestPosition := 0
	for i := 1; i < len(bb); i++ {
		if bb[i] < lowest {
			lowest = bb[i]
			lowestPosition = i
		}
	}

	block := mgl32.Vec4{}
	switch lowestPosition {
	case 0:
		block = block3
	case 1:
		block = block1
	case 2:
		block = block5
	case 3:
		block = block4
	case 4:
		block = block2
	case 5:
		block = block6
	}

	if w.Blocks[int(block.X())][int(block.Y())][int(block.Z())] != nil {
		return nil
	}

	return &block
}

func (w *World) GenerateWorld() {

	generatedWorld := make(WorldBlocks)
	for x := int(-w.Size.X()); x < int(w.Size.X()); x++ {
		if len(generatedWorld[x]) == 0 {
			generatedWorld[x] = make(map[int]map[int]*game_objects.Block)
		}
		for z := int(-w.Size.Z()); z < int(w.Size.Z()); z++ {

			y := float64(w.Size.Y()/2) + math.Round(float64(w.Size.Y()-w.Size.Y()/2)*w.GlobalNoise.Noise2D(float64(x)/float64(w.Size.X()), float64(z)/float64(w.Size.Z())))

			for i := int(y); i >= int(math.Max(0, float64(int(y)-int(w.Size.Y()/2)))); i-- {
				if generatedWorld[x][i] == nil {
					generatedWorld[x][i] = make(map[int]*game_objects.Block)
				}

				newCube := game_objects.NewBlock(float32(x), float32(i), float32(z), 1, true, false, game_objects.BlockStone)
				newCube.WithEdges = false

				generatedWorld[x][i][z] = &newCube
			}

		}

	}

	for x := int(-w.Size.X()); x < int(w.Size.X()); x++ {
		for z := int(-w.Size.Z()); z < int(w.Size.Z()); z++ {
			for y := int(-w.Size.Y()); y < int(w.Size.Y()); y++ {

				if y < 50 && generatedWorld[x][y][z] == nil {

					index := y

					for {
						if generatedWorld[x][index][z] == nil {
							if len(generatedWorld[x][index]) == 0 {
								generatedWorld[x][index] = make(map[int]*game_objects.Block)
							}

							newCube := game_objects.NewBlock(float32(x), float32(index), float32(z), 1, true, false, game_objects.BlockWater)
							newCube.WithEdges = false

							generatedWorld[x][index][z] = &newCube
							index--
							if index <= int(w.Size.Y()) {
								break
							}
						} else {
							break
						}

					}

				}

			}

		}

	}

	for x := int(-w.Size.X()); x < int(w.Size.X()); x++ {
		for z := int(-w.Size.Z()); z < int(w.Size.Z()); z++ {
			for y := int(-w.Size.Y()); y < int(w.Size.Y()); y++ {
				if generatedWorld[x][y][z] == nil {
					continue
				}

				if generatedWorld[x][y+1][z] == nil && generatedWorld[x][y][z].BlockType != game_objects.BlockWater {
					generatedWorld[x][y][z].BlockType = game_objects.BlockGrass
					continue
				}

				if generatedWorld[x][y][z].BlockType == game_objects.BlockWater {
					continue
				}

				if generatedWorld[x+1][y][z] != nil && generatedWorld[x+1][y][z].BlockType == game_objects.BlockWater ||
					generatedWorld[x-1][y][z] != nil && generatedWorld[x-1][y][z].BlockType == game_objects.BlockWater ||
					generatedWorld[x][y+1][z] != nil && generatedWorld[x][y+1][z].BlockType == game_objects.BlockWater ||
					generatedWorld[x][y-1][z] != nil && generatedWorld[x][y-1][z].BlockType == game_objects.BlockWater ||
					generatedWorld[x][y][z+1] != nil && generatedWorld[x][y][z+1].BlockType == game_objects.BlockWater ||
					generatedWorld[x][y][z-1] != nil && generatedWorld[x][y][z-1].BlockType == game_objects.BlockWater {
					generatedWorld[x][y][z].BlockType = game_objects.BlockSand
					continue
				}

			}

		}

	}

	w.Blocks = generatedWorld
	w.SetInitialNeighbors()
	w.ShouldUpdatePopulatedBlocks = true
}

func (w *World) SetInitialNeighbors() {
	for x := int(-w.Size.X()); x < int(w.Size.X()); x++ {
		for y := int(-w.Size.Y()); y < int(w.Size.Y()); y++ {
			for z := int(-w.Size.Z()); z < int(w.Size.Z()); z++ {
				blockPositionX, blockPositionY, blockPositionZ := int(x), int(y), int(z)

				if w.Blocks[blockPositionX][blockPositionY][blockPositionZ] == nil {
					continue
				}

				if w.Blocks[blockPositionX+1][blockPositionY][blockPositionZ] != nil {
					if w.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == game_objects.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[0] = true

					} else if w.Blocks[blockPositionX+1][blockPositionY][blockPositionZ].BlockType == game_objects.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[0] = false
					} else {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[0] = true
					}

				}
				if w.Blocks[blockPositionX-1][blockPositionY][blockPositionZ] != nil {
					if w.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == game_objects.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[1] = true

					} else if w.Blocks[blockPositionX-1][blockPositionY][blockPositionZ].BlockType == game_objects.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[1] = false

					} else {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[1] = true
					}

				}
				if w.Blocks[blockPositionX][blockPositionY][blockPositionZ+1] != nil {
					if w.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == game_objects.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[2] = true

					} else if w.Blocks[blockPositionX][blockPositionY][blockPositionZ+1].BlockType == game_objects.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[2] = false

					} else {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[2] = true
					}

				}
				if w.Blocks[blockPositionX][blockPositionY][int(blockPositionZ-1)] != nil {
					if w.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == game_objects.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[3] = true

					} else if w.Blocks[blockPositionX][blockPositionY][blockPositionZ-1].BlockType == game_objects.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[3] = false

					} else {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[3] = true
					}

				}
				if w.Blocks[blockPositionX][blockPositionY+1][blockPositionZ] != nil {
					if w.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == game_objects.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[4] = true

					} else if w.Blocks[blockPositionX][blockPositionY+1][blockPositionZ].BlockType == game_objects.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[4] = false

					} else {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[4] = true
					}

				}
				if w.Blocks[blockPositionX][blockPositionY-1][blockPositionZ] != nil {
					if w.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == game_objects.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[5] = true

					} else if w.Blocks[blockPositionX][blockPositionY-1][blockPositionZ].BlockType == game_objects.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[5] = false

					} else {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[5] = true
					}

				}
			}
		}
	}
}

func (w *World) InitPopulatedBlocks() {
	blockTypes := game_objects.GetBlockTypes()
	for _, _ = range blockTypes {
		w.PopulatedBlocks = append(w.PopulatedBlocks, make([]*game_objects.Block, 0))
	}
}

func (w *World) UpdatePopulatedBlocks(fromX, toX, fromY, toY, fromZ, toZ float64) {
	blockTypes := game_objects.GetBlockTypes()
	for _, blockType := range blockTypes {
		w.PopulatedBlocks[blockType] = []*game_objects.Block{}
	}

	for x := fromX; x < toX; x++ {
		intX := int(x)
		for y := 0; y < int(w.Size.Y()); y++ {

			intY := int(y)
			for z := fromZ; z < toZ; z++ {
				intZ := int(z)

				if w.Blocks[intX][intY][intZ] == nil || w.Blocks[intX][intY][intZ].CountNeighbors() == 6 {
					continue
				}

				w.PopulatedBlocks[w.Blocks[intX][intY][intZ].BlockType] = append(w.PopulatedBlocks[w.Blocks[intX][intY][intZ].BlockType], w.Blocks[intX][intY][intZ])

			}
		}
	}

	w.ShouldUpdatePopulatedBlocks = false
}

func (w *World) Update(roundedPlayerPosition mgl32.Vec3, backOfPlayer, frontOfPlayer mgl32.Vec3) {

	maxDist := float64(configs.ViewDistance)

	roundedPlayerX, roundedPlayerY, roundedPlayerZ := roundedPlayerPosition.Elem()
	fromX, toX := math.Max(-float64(w.Size.X()), float64(roundedPlayerX)-maxDist), math.Min(float64(w.Size.X()), float64(roundedPlayerX)+maxDist)
	fromY, toY := math.Max(-float64(w.Size.Y()), float64(roundedPlayerY)-maxDist), math.Min(float64(w.Size.Y()), float64(roundedPlayerY)+maxDist)
	fromZ, toZ := math.Max(-float64(w.Size.Z()), float64(roundedPlayerZ)-maxDist), math.Min(float64(w.Size.Z()), float64(roundedPlayerZ)+maxDist)

	// fmt.Println(len(w.PopulatedBlocks))

	sumX := frontOfPlayer.X() - backOfPlayer.X()
	sumZ := frontOfPlayer.Z() - backOfPlayer.Z()

	for _, populatedBlockTypes := range w.PopulatedBlocks {
		for _, populatedBlock := range populatedBlockTypes {
			if sumX < 0 && sumZ < 0 {
				if populatedBlock.Position.X() > backOfPlayer.X() && populatedBlock.Position.Z() > backOfPlayer.Z() {
					continue
				}
			} else if sumX > 0 && sumZ > 0 {
				if populatedBlock.Position.X() < backOfPlayer.X() && populatedBlock.Position.Z() < backOfPlayer.Z() {
					continue
				}
			} else if sumX > 0 && sumZ < 0 {
				if populatedBlock.Position.X() < backOfPlayer.X() && populatedBlock.Position.Z() > backOfPlayer.Z() {
					continue
				}
			} else if sumX < 0 && sumZ > 0 {
				if populatedBlock.Position.X() > backOfPlayer.X() && populatedBlock.Position.Z() < backOfPlayer.Z() {
					continue
				}
			}
			populatedBlock.Draw2()
		}

	}

	// fmt.Println(playerBehind)

	w.Tick++
	if w.Tick >= configs.TickRate {
		w.Tick = 0
		w.ShouldUpdatePopulatedBlocks = true
	}

	if w.ShouldUpdatePopulatedBlocks {
		w.UpdatePopulatedBlocks(fromX, toX, fromY, toY, fromZ, toZ)
	}

}
