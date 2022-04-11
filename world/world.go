package world

import (
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/block"
	"github.com/reonardoleis/fcg-glcraft/configs"
	math2 "github.com/reonardoleis/fcg-glcraft/math"
	"github.com/reonardoleis/fcg-glcraft/world/chunk"
	"github.com/tbogdala/noisey"
)

type WorldBlocks = map[int]map[int]map[int]*block.Block

type World struct {
	Name                        string
	Size                        mgl32.Vec3
	Blocks                      WorldBlocks
	Chunks                      map[int]map[int]*chunk.Chunk
	FutureChunks                map[int]map[int]*chunk.Chunk
	ShouldUpdateChunks          bool
	PopulatedBlocks             [][]*block.Block
	ShouldUpdatePopulatedBlocks bool
	Seed                        int64
	Time                        int64
	Tick                        float64
	GlobalNoise                 *noisey.OpenSimplexGenerator
}

func NewWorld(worldName string, size mgl32.Vec3, seed int64) *World {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	noiser := noisey.NewOpenSimplexGenerator(r)
	return &World{
		Name:               worldName,
		Size:               size,
		Seed:               seed,
		Time:               0,
		GlobalNoise:        &noiser,
		FutureChunks:       make(map[int]map[int]*chunk.Chunk),
		ShouldUpdateChunks: false,
	}
}

func (w World) GetBlockFrom(wx, wy, wz int, playerSize float32) *block.Block {
	return w.Blocks[wx][wy-int(playerSize)][wz]
}

func (w World) FindHighestBlock(wx, wz int) *block.Block {
	var highestBlock *block.Block
	keys := []int{}

	for k := range w.Blocks[wx] {
		keys = append(keys, k)
	}

	for wy := 0; wy < len(keys); wy++ {
		if highestBlock == nil {
			highestBlock = w.Blocks[wx][keys[wy]][wz]
		} else if w.Blocks[wx][keys[wy]][wz] != nil && w.Blocks[wx][keys[wy]][wz].Position.Y() > highestBlock.Position.Y() && w.Blocks[wx][keys[wy]][wz].BlockType != block.BlockWater {
			highestBlock = w.Blocks[wx][keys[wy]][wz]
		}
	}

	return highestBlock
}

func (w *World) HandleChunkChange(currentChunk *chunk.Chunk) {
	for i := currentChunk.Offset[0] - 2; i <= currentChunk.Offset[0]+2; i++ {
		if len(w.FutureChunks[int(i)]) == 0 {
			w.FutureChunks[int(i)] = make(map[int]*chunk.Chunk)
		}
		for j := currentChunk.Offset[1] - 2; j <= currentChunk.Offset[1]+2; j++ {
			if w.FutureChunks[int(i)][int(j)] == nil {
				w.FutureChunks[int(i)][int(j)] = chunk.NewChunk(mgl32.Vec2{float32(i), float32(j)}, 0)
				w.FutureChunks[int(i)][int(j)].GenerateChunk(w.GlobalNoise)
			}

		}
	}
}

func (w *World) GenerateWorld() {
	w.Chunks = make(map[int]map[int]*chunk.Chunk)
	for i := -10; i <= 10; i++ {
		w.Chunks[i] = make(map[int]*chunk.Chunk)
		for j := -10; j <= 10; j++ {
			w.Chunks[i][j] = chunk.NewChunk(mgl32.Vec2{float32(i), float32(j)}, 0)
			w.Chunks[i][j].GenerateChunk(w.GlobalNoise)
		}
	}
}

func (w World) GetChunk() {}

func (w *World) SetInitialNeighbors() {
	for x := int(-w.Size.X()); x < int(w.Size.X()); x++ {
		for y := int(-w.Size.Y()); y < int(w.Size.Y()); y++ {
			for z := int(-w.Size.Z()); z < int(w.Size.Z()); z++ {
				blockPositionX, blockPositionY, blockPositionZ := int(x), int(y), int(z)

				if w.Blocks[blockPositionX][blockPositionY][blockPositionZ] == nil {
					continue
				}

				if w.Blocks[blockPositionX+1][blockPositionY][blockPositionZ] != nil {
					if w.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[0] = 1

					} else if w.Blocks[blockPositionX+1][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[0] = 0
					} else {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[0] = 1
					}

				}
				if w.Blocks[blockPositionX-1][blockPositionY][blockPositionZ] != nil {
					if w.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[1] = 1

					} else if w.Blocks[blockPositionX-1][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[1] = 0

					} else {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[1] = 1
					}

				}
				if w.Blocks[blockPositionX][blockPositionY][blockPositionZ+1] != nil {
					if w.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[2] = 1

					} else if w.Blocks[blockPositionX][blockPositionY][blockPositionZ+1].BlockType == block.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[2] = 0

					} else {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[2] = 1
					}

				}
				if w.Blocks[blockPositionX][blockPositionY][int(blockPositionZ-1)] != nil {
					if w.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[3] = 1

					} else if w.Blocks[blockPositionX][blockPositionY][blockPositionZ-1].BlockType == block.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[3] = 0

					} else {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[3] = 1
					}

				}
				if w.Blocks[blockPositionX][blockPositionY+1][blockPositionZ] != nil {
					if w.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[4] = 1

					} else if w.Blocks[blockPositionX][blockPositionY+1][blockPositionZ].BlockType == block.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[4] = 0

					} else {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[4] = 1
					}

				}
				if w.Blocks[blockPositionX][blockPositionY-1][blockPositionZ] != nil {
					if w.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[5] = 1

					} else if w.Blocks[blockPositionX][blockPositionY-1][blockPositionZ].BlockType == block.BlockWater {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[5] = 0

					} else {
						w.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[5] = 1
					}

				}
			}
		}
	}
}

func (w *World) InitPopulatedBlocks() {
	blockTypes := block.GetBlockTypes()
	for _, _ = range blockTypes {
		w.PopulatedBlocks = append(w.PopulatedBlocks, make([]*block.Block, 0))
	}
}

func (w *World) PopulateIfEmpty(position mgl32.Vec3) {
	if len(w.Blocks[int(position.X())]) == 0 {
		w.Blocks[int(position.X())] = make(map[int]map[int]*block.Block)
	}
	if len(w.Blocks[int(position.X())][int(position.Y())]) == 0 {
		w.Blocks[int(position.X())][int(position.Y())] = make(map[int]*block.Block)
	}
	if w.Blocks[int(position.X())][int(position.Y())][int(position.Z())] == nil {
		w.Blocks[int(position.X())][int(position.Y())][int(position.Z())] = &block.Block{}
	}
}

func (w *World) PlaceTree(position mgl32.Vec3) {
	treeHeight := math2.RandInt(2, 4)

	for i := 0; i <= treeHeight; i++ {
		blockPosition := mgl32.Vec3{position.X(), position.Y() + float32(i), position.Z()}
		w.PopulateIfEmpty(blockPosition)
		treeBlock := block.NewBlock(position.X(), position.Y()+float32(i), position.Z(), float32(configs.BlockSize), false, false, block.BlockWood)
		w.Blocks[int(position.X())][int(position.Y())+i][int(position.Z())] = &treeBlock
	}

	treeMax := position.Y() + float32(treeHeight)

	for x := -1 + int(position.X()); x <= int(position.X())+1; x++ {
		for y := int(treeMax); y <= int(treeMax)+1; y++ {
			for z := -1 + int(position.Z()); z <= int(position.Z())+1; z++ {
				if float32(x) == position.X() && float32(z) == position.Z() {
					continue
				}
				blockPosition := mgl32.Vec3{float32(x), float32(y), float32(z)}
				w.PopulateIfEmpty(blockPosition)
				treeBlock := block.NewBlock(float32(x), float32(y), float32(z), float32(configs.BlockSize), false, false, block.BlockLeaves)
				w.Blocks[x][y][z] = &treeBlock
			}
		}
	}

	blockPosition := mgl32.Vec3{position.X(), treeMax + 2, position.Z()}
	w.PopulateIfEmpty(blockPosition)
	treeBlock := block.NewBlock(blockPosition.X(), blockPosition.Y(), blockPosition.Z(), float32(configs.BlockSize), false, false, block.BlockLeaves)
	w.Blocks[int(blockPosition.X())][int(blockPosition.Y())][int(blockPosition.Z())] = &treeBlock
}

func (w *World) UpdatePopulatedBlocks(fromX, toX, fromY, toY, fromZ, toZ float64, playerPosition mgl32.Vec4) {
	blockTypes := block.GetBlockTypes()
	for _, blockType := range blockTypes {
		w.PopulatedBlocks[blockType] = []*block.Block{}
	}

	for x := fromX; x < toX; x++ {
		if (len(w.Blocks[int(x)])) == 0 {
			continue
		}
		intX := int(x)
		for y := fromY; y < toY; y++ {
			if (len(w.Blocks[int(x)][int(y)])) == 0 {
				continue
			}
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

	for i := 0; i < len(w.PopulatedBlocks[block.BlockGlass]); i++ {
		for j := i; j < len(w.PopulatedBlocks[block.BlockGlass]); j++ {
			if math2.Distance(w.PopulatedBlocks[block.BlockGlass][i].Position, playerPosition) < math2.Distance(w.PopulatedBlocks[block.BlockGlass][j].Position, playerPosition) {
				temp := w.PopulatedBlocks[block.BlockGlass][i]
				w.PopulatedBlocks[block.BlockGlass][i] = w.PopulatedBlocks[block.BlockGlass][j]
				w.PopulatedBlocks[block.BlockGlass][j] = temp
			}
		}
	}

	w.ShouldUpdatePopulatedBlocks = false
}

func (w *World) Update(roundedPlayerPosition mgl32.Vec3, backOfPlayer, frontOfPlayer mgl32.Vec3, currentChunk *chunk.Chunk) {

	/*maxDist := float64(configs.ViewDistance)

	roundedPlayerX, roundedPlayerY, roundedPlayerZ := roundedPlayerPosition.Elem()
	fromX, toX := math.Max(-float64(w.Size.X()), float64(roundedPlayerX)-maxDist), math.Min(float64(w.Size.X()), float64(roundedPlayerX)+maxDist)
	fromY, toY := math.Max(-float64(w.Size.Y()), float64(roundedPlayerY)-maxDist), math.Min(float64(w.Size.Y()), float64(roundedPlayerY)+maxDist)
	fromZ, toZ := math.Max(-float64(w.Size.Z()), float64(roundedPlayerZ)-maxDist), math.Min(float64(w.Size.Z()), float64(roundedPlayerZ)+maxDist)

	// // fmt.Println(len(w.PopulatedBlocks))

	sumX := frontOfPlayer.X() - backOfPlayer.X()
	sumZ := frontOfPlayer.Z() - backOfPlayer.Z()

	for _, populatedBlockTypes := range w.PopulatedBlocks {
		if len(populatedBlockTypes) == 0 {
			continue
		}
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
			populatedBlock.WithEdges = false
			populatedBlock.Colliding = false
		}

	}

	// // fmt.Println(playerBehind)

	w.Tick++
	if w.Tick >= configs.TickRate {
		w.Tick = 0
		w.ShouldUpdatePopulatedBlocks = true
	}

	if w.ShouldUpdatePopulatedBlocks {
		w.UpdatePopulatedBlocks(fromX, toX, fromY, toY, fromZ, toZ, mgl32.Vec4{roundedPlayerPosition.X(), roundedPlayerPosition.Y(), roundedPlayerPosition.Z(), 1.0})
	}*/

	for i := currentChunk.Offset[0] - configs.ViewDistance; i <= currentChunk.Offset[0]+configs.ViewDistance; i++ {
		for j := currentChunk.Offset[1] - configs.ViewDistance; j <= currentChunk.Offset[1]+configs.ViewDistance; j++ {
			if w.Tick >= configs.TickRate {
				go w.Chunks[int(i)][int(j)].Update()
				go w.Chunks[int(i)][int(j)].SetWatersUpdate()
			}
			chunkRenderableBlocks := w.Chunks[int(i)][int(j)].GetBlocksToRender()
			for _, renderableBlock := range chunkRenderableBlocks {
				renderableBlock.Draw2()
			}
		}
	}

	if w.Tick >= configs.TickRate {
		w.Tick = 0
	}

	w.Tick += math2.DeltaTime

}
