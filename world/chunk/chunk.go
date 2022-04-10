package chunk

import (
	"math"

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

type Chunk struct {
	ID        uint64
	Offset    mgl32.Vec2 // identifies the chunk position in world
	BiomeType BiomeType
	Blocks    [][][]*block.Block
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
		for y := 0; y < int(configs.WorldHeight); y++ {
			c.Blocks[x] = append(c.Blocks[x], []*block.Block{})
			for z := 0; z < int(configs.ChunkSize); z++ {
				c.Blocks[x][y] = append(c.Blocks[x][y], nil)
			}
		}
	}
}

func (c *Chunk) SetInitialNeighbors() {
	for x := 0; x < configs.ChunkSize; x++ {
		for y := 0; y < configs.WorldHeight; y++ {
			for z := 0; z < configs.ChunkSize; z++ {
				blockPositionX, blockPositionY, blockPositionZ := int(x), int(y), int(z)

				if c.Blocks[blockPositionX][blockPositionY][blockPositionZ] == nil {
					continue
				}

				if blockPositionX+1 < configs.ChunkSize && c.Blocks[blockPositionX+1][blockPositionY][blockPositionZ] != nil {
					if c.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[0] = 1

					} else if c.Blocks[blockPositionX+1][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[0] = 0
					} else {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[0] = 1
					}

				}
				if blockPositionX-1 >= 0 && c.Blocks[blockPositionX-1][blockPositionY][blockPositionZ] != nil {
					if c.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[1] = 1

					} else if c.Blocks[blockPositionX-1][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[1] = 0

					} else {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[1] = 1
					}

				}
				if blockPositionZ+1 < configs.ChunkSize && c.Blocks[blockPositionX][blockPositionY][blockPositionZ+1] != nil {
					if c.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[2] = 1

					} else if c.Blocks[blockPositionX][blockPositionY][blockPositionZ+1].BlockType == block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[2] = 0

					} else {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[2] = 1
					}

				}
				if blockPositionZ-1 >= 0 && c.Blocks[blockPositionX][blockPositionY][int(blockPositionZ-1)] != nil {
					if c.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[3] = 1

					} else if c.Blocks[blockPositionX][blockPositionY][blockPositionZ-1].BlockType == block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[3] = 0

					} else {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[3] = 1
					}

				}
				if blockPositionY+1 < configs.WorldHeight && c.Blocks[blockPositionX][blockPositionY+1][blockPositionZ] != nil {
					if c.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[4] = 1

					} else if c.Blocks[blockPositionX][blockPositionY+1][blockPositionZ].BlockType == block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[4] = 0

					} else {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[4] = 1
					}

				}
				if blockPositionY-1 >= 0 && c.Blocks[blockPositionX][blockPositionY-1][blockPositionZ] != nil {
					if c.Blocks[blockPositionX][blockPositionY][blockPositionZ].BlockType == block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[5] = 1

					} else if c.Blocks[blockPositionX][blockPositionY-1][blockPositionZ].BlockType == block.BlockWater {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[5] = 0

					} else {
						c.Blocks[blockPositionX][blockPositionY][blockPositionZ].Neighbors[5] = 1
					}

				}
			}
		}
	}
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

			normalizingQuotientX := float64(configs.BiomeChunks * configs.ChunkSize)
			normalizingQuotientZ := float64(configs.BiomeChunks * configs.ChunkSize)

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
				if y < 32 && c.Blocks[x][y][z] == nil {

					for index := y; index > 1; index-- {

						if c.Blocks[x][index][z] == nil {
							//fmt.Println("gerando agua")
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

	c.SetInitialNeighbors()
}

func (c *Chunk) GetBlockAt(x, y, z int) *block.Block {
	_x := x - (int(c.Offset[0] * float32(configs.ChunkSize)))
	_z := z - (int(c.Offset[1] * float32(configs.ChunkSize)))
	if _x < 0 || _x >= configs.ChunkSize || y < 0 || y >= configs.WorldHeight || _z < 0 || _z >= configs.ChunkSize {
		return nil
	}

	return c.Blocks[_x][y][_z]
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
	}

	c.Blocks[int(position.X())][int(position.Y())][int(position.Z())] = nil
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
}

func (c *Chunk) GetBlocksToRender() []*block.Block {
	blocksToRender := make([]*block.Block, 0)
	for _, x := range c.Blocks {
		for _, y := range x {
			for _, block := range y {
				if block != nil && block.CountNeighbors() != 6 {
					blocksToRender = append(blocksToRender, block)
					continue
				}

			}
		}
	}

	return blocksToRender
}
