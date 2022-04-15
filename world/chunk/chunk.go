package chunk

import (
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/block"
	"github.com/reonardoleis/fcg-glcraft/configs"
	math2 "github.com/reonardoleis/fcg-glcraft/math"

	"github.com/tbogdala/noisey"
)

type BiomeType byte

var (
	SerialChunkID uint64 = 0
)

type BlockInformation = byte

const (
	BlockInformationNone = iota
	BlockInformationCave
)

type Chunk struct {
	ID                uint64
	Offset            mgl32.Vec2 // identifies the chunk position in world
	BiomeType         BiomeType
	Blocks            [][][]*block.Block
	BlocksInformation [][][]BlockInformation
}

func NewChunk(offset mgl32.Vec2, biomeType BiomeType) *Chunk {
	defer (func() {
		SerialChunkID++
	})()
	return &Chunk{
		ID:        SerialChunkID,
		Offset:    offset,
		BiomeType: biomeType,
	}
}

func (c *Chunk) allocateBlockSlice() {
	for x := 0; x < int(configs.ChunkSize); x++ {
		c.Blocks = append(c.Blocks, [][]*block.Block{})
		c.BlocksInformation = append(c.BlocksInformation, [][]BlockInformation{})
		for y := 0; y < int(configs.WorldHeight); y++ {
			c.Blocks[x] = append(c.Blocks[x], []*block.Block{})
			c.BlocksInformation[x] = append(c.BlocksInformation[x], []BlockInformation{})
			for z := 0; z < int(configs.ChunkSize); z++ {
				c.Blocks[x][y] = append(c.Blocks[x][y], nil)
				c.BlocksInformation[x][y] = append(c.BlocksInformation[x][y], BlockInformationNone)
			}
		}
	}
}

func (c *Chunk) SetNeighbors() {
	for x := 0; x < configs.ChunkSize; x++ {
		for y := 0; y < configs.WorldHeight; y++ {
			for z := 0; z < configs.ChunkSize; z++ {
				blockPositionX, blockPositionY, blockPositionZ := int(x), int(y), int(z)

				if c.Blocks[blockPositionX][blockPositionY][blockPositionZ] == nil {
					continue
				}

				if blockPositionX+1 < configs.ChunkSize && c.Blocks[blockPositionX+1][blockPositionY][blockPositionZ] != nil {
					if c.Blocks[blockPositionX+1][blockPositionY][blockPositionZ].Transparent && c.Blocks[blockPositionX+1][blockPositionY][blockPositionZ].BlockType != block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[0] = 0
					} else if c.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						// handle weak water neighbors
						if c.Blocks[blockPositionX+1][blockPositionY][blockPositionZ].WaterForce != 8 {
							c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[0] = 0
						}

					} else if c.Blocks[blockPositionX+1][blockPositionY][blockPositionZ].BlockType == block.BlockWater || c.Blocks[blockPositionX+1][blockPositionY][blockPositionZ].BlockType == block.BlockAir {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[0] = 0
					} else {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[0] = 1
					}

				}
				if blockPositionX-1 >= 0 && c.Blocks[blockPositionX-1][blockPositionY][blockPositionZ] != nil {
					if c.Blocks[blockPositionX-1][blockPositionY][blockPositionZ].Transparent && c.Blocks[blockPositionX-1][blockPositionY][blockPositionZ].BlockType != block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[1] = 0
					} else if c.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[1] = 1

						// handle weak water neighbors
						if c.Blocks[blockPositionX-1][blockPositionY][blockPositionZ].WaterForce != 8 {
							c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[1] = 0
						}

					} else if c.Blocks[blockPositionX-1][blockPositionY][blockPositionZ].BlockType == block.BlockWater || c.Blocks[blockPositionX-1][blockPositionY][blockPositionZ].BlockType == block.BlockAir {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[1] = 0

					} else {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[1] = 1
					}

				}
				if blockPositionZ+1 < configs.ChunkSize && c.Blocks[blockPositionX][blockPositionY][blockPositionZ+1] != nil {
					if c.Blocks[blockPositionX][blockPositionY][blockPositionZ+1].Transparent && c.Blocks[blockPositionX][blockPositionY][blockPositionZ+1].BlockType != block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[2] = 0
					} else if c.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[2] = 1

						// handle weak water neighbors
						if c.Blocks[blockPositionX][blockPositionY][blockPositionZ+1].WaterForce != 8 {
							c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[2] = 0
						}

					} else if c.Blocks[blockPositionX][blockPositionY][blockPositionZ+1].BlockType == block.BlockWater || c.Blocks[blockPositionX][blockPositionY][blockPositionZ+1].BlockType == block.BlockAir {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[2] = 0

					} else {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[2] = 1
					}

				}
				if blockPositionZ-1 >= 0 && c.Blocks[blockPositionX][blockPositionY][int(blockPositionZ-1)] != nil {
					if c.Blocks[blockPositionX][blockPositionY][blockPositionZ-1].Transparent && c.Blocks[blockPositionX][blockPositionY][blockPositionZ-1].BlockType != block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[3] = 0
					} else if c.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[3] = 1

						// handle weak water neighbors
						if c.Blocks[blockPositionX][blockPositionY][blockPositionZ-1].WaterForce != 8 {
							c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[3] = 0
						}

					} else if c.Blocks[blockPositionX][blockPositionY][blockPositionZ-1].BlockType == block.BlockWater || c.Blocks[blockPositionX][blockPositionY][blockPositionZ-1].BlockType == block.BlockAir {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[3] = 0

					} else {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[3] = 1
					}

				}
				if blockPositionY+1 < configs.WorldHeight && c.Blocks[blockPositionX][blockPositionY+1][blockPositionZ] != nil {
					if c.Blocks[blockPositionX][blockPositionY+1][blockPositionZ].Transparent && c.Blocks[blockPositionX][blockPositionY+1][blockPositionZ].BlockType != block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[4] = 0
					} else if c.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[4] = 1

						// handle weak water neighbors
						if c.Blocks[blockPositionX][blockPositionY+1][blockPositionZ].WaterForce != 8 {
							c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[4] = 0
						}

						// If the water block has a water block above it then it should not be scaled down
						if c.Blocks[blockPositionX][blockPositionY+1][blockPositionZ].BlockType == block.BlockWater {
							c.Blocks[blockPositionX][blockPositionY][blockPositionZ].HasWaterAbove = true
						}

					} else if c.Blocks[blockPositionX][blockPositionY+1][blockPositionZ].BlockType == block.BlockWater || c.Blocks[blockPositionX][blockPositionY+1][blockPositionZ].BlockType == block.BlockAir {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[4] = 0

					} else {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[4] = 1
					}

				}
				if blockPositionY-1 >= 0 && c.Blocks[blockPositionX][blockPositionY-1][blockPositionZ] != nil {
					if c.Blocks[blockPositionX][blockPositionY-1][blockPositionZ].Transparent && c.Blocks[blockPositionX][blockPositionY-1][blockPositionZ].BlockType != block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[5] = 0
					} else if c.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[5] = 1

						// handle weak water neighbors
						if c.Blocks[blockPositionX][blockPositionY-1][blockPositionZ].WaterForce != 8 {
							c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[5] = 0
						}

					} else if c.Blocks[blockPositionX][blockPositionY-1][blockPositionZ].BlockType == block.BlockWater || c.Blocks[blockPositionX][blockPositionY-1][blockPositionZ].BlockType == block.BlockAir {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[5] = 0

					} else {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[5] = 1
					}

				}
			}
		}
	}
}

func (c *Chunk) PlaceTree(x, y, z int) {
	position := mgl32.Vec3{float32(x), float32(y), float32(z)}
	treeHeight := math2.RandInt(2, 4)

	for i := 0; i <= treeHeight; i++ {
		blockPosition := position.Add(mgl32.Vec3{0.0, float32(i), 0.0})
		treeBlock := block.NewBlock(float32(x)+(c.Offset[0]*float32(configs.ChunkSize)), blockPosition.Y(), float32(z)+(c.Offset[1]*float32(configs.ChunkSize)), float32(configs.BlockSize), false, false, block.BlockWood)
		c.AddBlockAtNotOffsetted(int(blockPosition.X()), int(blockPosition.Y()), int(blockPosition.Z()), &treeBlock)
	}

	treeMax := position.Y() + float32(treeHeight)

	for x := -1 + int(position.X()); x <= int(position.X())+1; x++ {
		for y := int(treeMax); y <= int(treeMax)+1; y++ {
			for z := -1 + int(position.Z()); z <= int(position.Z())+1; z++ {
				if float32(x) == position.X() && float32(z) == position.Z() {
					continue
				}
				blockPosition := mgl32.Vec3{float32(x), float32(y), float32(z)}
				treeBlock := block.NewBlock(float32(x)+(c.Offset[0]*float32(configs.ChunkSize)), blockPosition.Y(), float32(z)+(c.Offset[1]*float32(configs.ChunkSize)), float32(configs.BlockSize), false, false, block.BlockLeaves)
				c.AddBlockAtNotOffsetted(int(blockPosition.X()), int(blockPosition.Y()), int(blockPosition.Z()), &treeBlock)
			}
		}
	}

	blockPosition := mgl32.Vec3{position.X(), treeMax + 2, position.Z()}
	treeBlock := block.NewBlock(float32(x)+(c.Offset[0]*float32(configs.ChunkSize)), blockPosition.Y(), float32(z)+(c.Offset[1]*float32(configs.ChunkSize)), float32(configs.BlockSize), false, false, block.BlockLeaves)
	c.AddBlockAtNotOffsetted(int(blockPosition.X()), int(blockPosition.Y()), int(blockPosition.Z()), &treeBlock)
}

func (c *Chunk) GenerateChunk(noiseSource *noisey.OpenSimplexGenerator) {
	c.allocateBlockSlice()
	for x := 0; x < int(configs.ChunkSize); x++ {
		for z := 0; z < int(configs.ChunkSize); z++ {
			/*blockHeight := float64(configs.WorldHeight/2) + math.Round(float64(configs.WorldHeight)/2)*
			noiseSource.Get2D((float64(x)+float64((float32(configs.ChunkSize)*c.Offset[0])))/float64((c.Offset[0]+1)*float32(configs.ChunkSize)),
				(float64(z)+float64(float32(configs.ChunkSize)*c.Offset[1])/float64((c.Offset[1]+1)*float32(configs.ChunkSize))))*/

			blockWithOffsetX := float64((x + (configs.ChunkSize * int(c.Offset[0]))))
			blockWithOffsetZ := float64((z + (configs.ChunkSize * int(c.Offset[1]))))

			normalizingQuotientX := float64(configs.ChunkSmoothness * configs.ChunkSize)
			normalizingQuotientZ := float64(configs.ChunkSmoothness * configs.ChunkSize)

			noiseParamX := blockWithOffsetX / normalizingQuotientX
			noiseParamZ := blockWithOffsetZ / normalizingQuotientZ

			// fmt.Println("Gerando bloco", blockWithOffsetX, blockWithOffsetZ, normalizingQuotientX, normalizingQuotientZ)

			blockHeight := noiseSource.Get2D(float64(noiseParamX), float64(noiseParamZ))
			blockHeight = math.Round(float64(configs.WorldHeight/2) + (math.Round(float64(configs.WorldHeight)/2) * blockHeight))
			for y := blockHeight; y >= 0; y-- {
				newBlock := block.NewBlock(float32(x)+(float32(configs.ChunkSize)*c.Offset[0]), float32(int(y)), float32(z)+(float32(configs.ChunkSize)*c.Offset[1]), float32(configs.BlockSize), false, false, block.BlockStone)
				c.Blocks[x][int(y)][z] = &newBlock
			}
		}
	}

	for x := 0; x < int(configs.ChunkSize); x++ {
		for y := 0; y < int(configs.WorldHeight); y++ {
			for z := 0; z < int(configs.ChunkSize); z++ {
				blockWithOffsetX := float64((x + (configs.ChunkSize * int(c.Offset[0]))))
				blockWithOffsetZ := float64((z + (configs.ChunkSize * int(c.Offset[1]))))

				normalizingQuotientX := float64(configs.ChunkSmoothness)
				normalizingQuotientY := float64(configs.ChunkSmoothness)
				normalizingQuotientZ := float64(configs.ChunkSmoothness)

				noiseParamX := blockWithOffsetX / normalizingQuotientX
				noiseParamY := float64(y) / normalizingQuotientY
				noiseParamZ := blockWithOffsetZ / normalizingQuotientZ

				noise := noiseSource.Get3D(noiseParamX, noiseParamY, noiseParamZ)
				if noise >= float64(configs.CaveThreshold) && y < configs.CaveMinHeight && y != 0 && c.Blocks[x][y][z] != nil && c.Blocks[x][y][z].BlockType != block.BlockWater {
					c.Blocks[x][y][z] = nil
					c.BlocksInformation[x][y][z] = BlockInformationCave
				}
			}
		}
	}

	for x := 0; x < int(configs.ChunkSize); x++ {
		for y := 0; y < int(configs.WorldHeight); y++ {
			for z := 0; z < int(configs.ChunkSize); z++ {

				shouldPlaceTree := math2.RandInt(0, 100) >= 95
				blockBelow := c.GetBlockAtNotOffsetted(x, y-1, z)
				if blockBelow != nil && blockBelow.BlockType == block.BlockGrass && shouldPlaceTree && x > 0 && x < configs.ChunkSize-1 && z > 0 && z < configs.ChunkSize-1 {
					c.PlaceTree(x, y, z)
				}

				if y < 32 && c.Blocks[x][y][z] == nil {

					for index := y; index > 1 && c.BlocksInformation[x][y][z] != BlockInformationCave; index-- {

						if c.Blocks[x][index][z] == nil {
							waterBlock := block.NewBlock(float32(x)+(float32(configs.ChunkSize)*c.Offset[0]), float32(index), float32(z)+(float32(configs.ChunkSize)*c.Offset[1]), float32(configs.BlockSize), false, false, block.BlockWater)
							c.Blocks[x][index][z] = &waterBlock
							if index <= 1 {
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

	r1 := rand.New(rand.NewSource(int64(math2.RandInt(0, 1000000000000000))))
	r2 := rand.New(rand.NewSource(int64(math2.RandInt(0, 1000000000000000))))
	coalNoiser := noisey.NewOpenSimplexGenerator(r1)
	ironNoiser := noisey.NewOpenSimplexGenerator(r2)
	for x := 0; x < configs.ChunkSize; x++ {
		for y := 0; y < configs.WorldHeight-1; y++ {
			for z := 0; z < configs.ChunkSize; z++ {
				currentBlock := c.Blocks[x][y][z]
				if currentBlock != nil {

					shouldPlaceGrass := true

					// surface-grass handling & ore generation
					if currentBlock.BlockType != block.BlockWater {
						for height := y + 1; height < configs.WorldHeight; height++ {
							if c.Blocks[x][height][z] != nil || c.BlocksInformation[x][height][z] == BlockInformationCave {
								shouldPlaceGrass = false
								if c.BlocksInformation[x][height][z] == BlockInformationCave {
									if noiseSource.Get3D(float64(x), float64(y), float64(z)) >= float64(configs.CaveDirtThreshold) {
										currentBlock.BlockType = block.BlockDirt
									}

								}

								coalNoise := coalNoiser.Get3D(float64(x)/float64(configs.CaveContentSmoothness), float64(y), float64(z)/float64(configs.CaveContentSmoothness))
								ironNoise := ironNoiser.Get3D(float64(x)/float64(configs.CaveContentSmoothness), float64(y), float64(z)/float64(configs.CaveContentSmoothness))
								if y < 50 && coalNoise >= configs.CaveCoalThreshold[0] && coalNoise <= configs.CaveCoalThreshold[1] {
									currentBlock.BlockType = block.BlockCoal
								}
								if y < 40 && ironNoise >= configs.CaveIronThreshold[0] && ironNoise <= configs.CaveIronThreshold[1] {
									currentBlock.BlockType = block.BlockIron
								}
								break
							}
						}

						if shouldPlaceGrass {
							currentBlock.BlockType = block.BlockGrass
						}

					}
				}
			}
		}
	}

	for x := 0; x < int(configs.ChunkSize); x++ {
		for y := 0; y < int(configs.WorldHeight); y++ {
			for z := 0; z < int(configs.ChunkSize); z++ {

				shouldPlaceTree := math2.RandInt(0, 100) >= 99
				blockBelow := c.GetBlockAtNotOffsetted(x, y-1, z)
				if blockBelow != nil && blockBelow.BlockType == block.BlockGrass && shouldPlaceTree && x > 0 && x < configs.ChunkSize-1 && z > 0 && z < configs.ChunkSize-1 {
					c.PlaceTree(x, y, z)
				}
			}
		}
	}

	c.SetNeighbors()
}

func (c *Chunk) GetBlockAtNotOffsetted(x, y, z int) *block.Block {
	_x := x
	_z := z
	if _x < 0 || _x >= configs.ChunkSize || y < 0 || y >= configs.WorldHeight || _z < 0 || _z >= configs.ChunkSize {
		return nil
	}

	return c.Blocks[_x][y][_z]
}

func (c *Chunk) GetBlockAt(x, y, z int) *block.Block {
	_x := x - (int(c.Offset[0] * float32(configs.ChunkSize)))
	_z := z - (int(c.Offset[1] * float32(configs.ChunkSize)))
	if _x < 0 || _x >= configs.ChunkSize || y < 0 || y >= configs.WorldHeight || _z < 0 || _z >= configs.ChunkSize {
		return nil
	}

	return c.Blocks[_x][y][_z]
}

func (c *Chunk) GetBlockInformationAt(x, y, z int) BlockInformation {
	_x := x - (int(c.Offset[0] * float32(configs.ChunkSize)))
	_z := z - (int(c.Offset[1] * float32(configs.ChunkSize)))
	if _x < 0 || _x >= configs.ChunkSize || y < 0 || y >= configs.WorldHeight || _z < 0 || _z >= configs.ChunkSize {
		return BlockInformationNone
	}

	return c.BlocksInformation[_x][y][_z]
}

func (c Chunk) FindPlacementPosition(hitAt mgl32.Vec4, nearFrom mgl32.Vec4, boundingBoxHighests, boundingBoxLowests mgl32.Vec3) *mgl32.Vec4 {
	block1 := mgl32.Vec4{hitAt.X(), hitAt.Y() + 1, hitAt.Z(), 1.0} // above
	block2 := mgl32.Vec4{hitAt.X(), hitAt.Y() - 1, hitAt.Z(), 1.0} // below
	block3 := mgl32.Vec4{hitAt.X() + 1, hitAt.Y(), hitAt.Z(), 1.0} // north
	block4 := mgl32.Vec4{hitAt.X() - 1, hitAt.Y(), hitAt.Z(), 1.0} // south
	block5 := mgl32.Vec4{hitAt.X(), hitAt.Y(), hitAt.Z() + 1, 1.0} // east
	block6 := mgl32.Vec4{hitAt.X(), hitAt.Y(), hitAt.Z() - 1, 1.0} // west
	// // fmt.Println(hitAt, nearFrom)
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

	foundBlock := c.GetBlockAt(int(block.X())-(int(c.Offset[0]*float32(configs.ChunkSize))), int(block.Y()), int(block.Z()-(c.Offset[1]*float32(configs.ChunkSize))))
	if foundBlock != nil {
		return nil
	}

	return &block
}

func (c *Chunk) RemoveBlockFrom(position mgl32.Vec4) {
	northNeighbor := math2.North(position, float32(configs.BlockSize)*2)
	southNeighbor := math2.South(position, float32(configs.BlockSize)*2)
	eastNeighbor := math2.East(position, float32(configs.BlockSize)*2)
	westNeighbor := math2.West(position, float32(configs.BlockSize)*2)
	upperNeighbor := math2.Upper(position, float32(configs.BlockSize)*2)
	lowerNeighbor := math2.Lower(position, float32(configs.BlockSize)*2)

	northBlock := c.GetBlockAt(int(northNeighbor.X()), int(northNeighbor.Y()), int(northNeighbor.Z()))

	southBlock := c.GetBlockAt(int(southNeighbor.X()), int(southNeighbor.Y()), int(southNeighbor.Z()))

	eastBlock := c.GetBlockAt(int(eastNeighbor.X()), int(eastNeighbor.Y()), int(eastNeighbor.Z()))

	westBlock := c.GetBlockAt(int(westNeighbor.X()), int(westNeighbor.Y()), int(westNeighbor.Z()))

	upperBlock := c.GetBlockAt(int(upperNeighbor.X()), int(upperNeighbor.Y()), int(upperNeighbor.Z()))

	lowerBlock := c.GetBlockAt(int(lowerNeighbor.X()), int(lowerNeighbor.Y()), int(lowerNeighbor.Z()))

	offsettedX, _, offsettedZ := c.GetOffsettedPositions(position.X(), position.Y(), position.Z())

	if northBlock != nil {
		northBlock.Neighbors[1] = 0
	}
	if southBlock != nil {
		southBlock.Neighbors[0] = 0
	}
	if eastBlock != nil {
		eastBlock.Neighbors[3] = 0
	}
	if westBlock != nil {
		westBlock.Neighbors[2] = 0
	}
	if upperBlock != nil {
		upperBlock.Neighbors[5] = 0
	}
	if lowerBlock != nil {
		lowerBlock.Neighbors[4] = 0
		if c.Blocks[int(offsettedX)][int(position.Y())][int(offsettedZ)].BlockType == block.BlockWater {
			lowerBlock.HasWaterAbove = false
		}

	}

	if c.Blocks[int(offsettedX)][int(position.Y())][int(offsettedZ)].IsBreakable {
		c.Blocks[int(offsettedX)][int(position.Y())][int(offsettedZ)] = nil
		c.SetNeighbors()
	}
}

func (c Chunk) GetOffsettedPositions(x, y, z float32) (float32, float32, float32) {
	return x - (c.Offset[0] * float32(configs.ChunkSize)), y, z - (c.Offset[1] * float32(configs.ChunkSize))
}

func (c *Chunk) AddBlockAt(position mgl32.Vec3, ephemeral bool, blockType block.BlockType) {
	x, y, z := position.Elem()

	offsettedX, _, offsettedZ := c.GetOffsettedPositions(x, y, z)

	newBlock := block.NewBlock(x, float32(y), z, 1, true, ephemeral, blockType)
	newBlock.WithEdges = false
	c.Blocks[int(offsettedX)][int(y)][int(offsettedZ)] = &newBlock

	c.SetNeighbors()
}

func (c *Chunk) AddBlockAtNotOffsetted(x, y, z int, block *block.Block) {
	_x := x
	_z := z
	if _x < 0 || _x >= configs.ChunkSize || y < 0 || y >= configs.WorldHeight || _z < 0 || _z >= configs.ChunkSize {
		return
	}

	c.Blocks[_x][y][_z] = block

	c.SetNeighbors()
}

func (c *Chunk) GetBlocksToRender() []*block.Block {
	blocksToRender := make([]*block.Block, 0)
	for _, x := range c.Blocks {
		for _, y := range x {
			for _, currentBlock := range y {
				if currentBlock != nil && currentBlock.CountNeighbors() != 6 && currentBlock.BlockType != block.BlockAir {
					blocksToRender = append(blocksToRender, currentBlock)
					continue
				}

			}
		}
	}

	return blocksToRender
}

func (c *Chunk) Update() {
	for x := 0; x < configs.ChunkSize; x++ {
		for y := 0; y < configs.WorldHeight; y++ {
			for z := 0; z < configs.ChunkSize; z++ {

				currentBlock := c.GetBlockAtNotOffsetted(x, y, z)

				if currentBlock == nil {
					continue
				}

				// handle sand & gravel blocks
				if currentBlock.BlockType == block.BlockSand && !currentBlock.IsFalling {
					blockBelow := c.GetBlockAtNotOffsetted(x, y-1, z)
					if blockBelow == nil {
						currentBlock.IsFalling = true
						currentBlock.IsBreakable = false
					}
				}

				if currentBlock.BlockType == block.BlockSand && currentBlock.IsFalling {
					blockBelow := c.GetBlockAtNotOffsetted(x, y-1, z)
					if blockBelow == nil {
						currentBlock.Position = currentBlock.Position.Sub(mgl32.Vec4{0.0, configs.BlockFallingSpeed * float32(math2.DeltaTime), 0.0, 1.0})
					}

					currentWorldPositionY := math.Ceil(float64(currentBlock.Position.Y()))
					blockBelow2 := c.GetBlockAtNotOffsetted(x, int(currentWorldPositionY-1), z)
					if blockBelow2 != nil {
						copy := *currentBlock
						copy.IsBreakable = true
						copy.IsFalling = false
						copy.Position = mgl32.Vec4{float32(x), blockBelow2.Position.Y() + 1, float32(z), 1.0}
						c.Blocks[x][y][z] = nil
						c.Blocks[x][int(blockBelow2.Position.Y()+1)][z] = &copy
						c.SetNeighbors()
					}
				}

				// end handle sand & gravel blocks

				// handle water blocks
				if currentBlock.BlockType == block.BlockWater {
					if currentBlock != nil && currentBlock.BlockType == block.BlockWater && currentBlock.SpreadThisTick {
						blockBelow := c.GetBlockAtNotOffsetted(x, y-1, z)
						if blockBelow != nil && blockBelow.BlockType == block.BlockWater {
							continue
						}
						if blockBelow == nil {
							// if it does not have a block below
							// then add a maximum-force water block
							newWaterBlock := block.NewBlock(float32(x)+(c.Offset[0]*float32(configs.ChunkSize)), float32(y-1), float32(z)+(c.Offset[1]*float32(configs.ChunkSize)), float32(configs.BlockSize), false, false, block.BlockWater)
							newWaterBlock.SpreadThisTick = false
							c.AddBlockAtNotOffsetted(x, y-1, z, &newWaterBlock)
							continue
						}

						if currentBlock.WaterForce > 1 {
							for dx := -1; dx <= 1; dx++ {
								for dz := -1; dz <= 1; dz++ {
									if (dx != 0 && dz != 0) || math.Abs(float64(dx)) != math.Abs(float64(dz)) {
										currentVerifying := c.GetBlockAtNotOffsetted(x+dx, y, z+dz)
										if currentVerifying == nil {
											belowCurrentVeryfing := c.GetBlockAtNotOffsetted(x+dx, y-1, z+dz)
											if belowCurrentVeryfing != nil && belowCurrentVeryfing.BlockType == block.BlockWater {
												continue
											}
											newWaterBlock := block.NewBlock(float32(x+dx)+(c.Offset[0]*float32(configs.ChunkSize)), float32(y), float32(z+dz)+(c.Offset[1]*float32(configs.ChunkSize)), float32(configs.BlockSize), false, false, block.BlockWater)
											newWaterBlock.WaterForce = currentBlock.WaterForce - 1
											newWaterBlock.SpreadThisTick = false
											c.AddBlockAtNotOffsetted(x+dx, y, z+dz, &newWaterBlock)
										}
									}
								}
							}
						}

					}
				}

				// end handle water blocks

			}
		}
	}

	c.SetNeighbors()
}

func (c *Chunk) SetWatersUpdate() {
	for x := 0; x < configs.ChunkSize; x++ {
		for y := 0; y < configs.WorldHeight; y++ {
			for z := 0; z < configs.ChunkSize; z++ {
				if c.Blocks[x][y][z] != nil && c.Blocks[x][y][z].BlockType == block.BlockWater {
					c.Blocks[x][y][z].SpreadThisTick = true
				}
			}
		}
	}
}
