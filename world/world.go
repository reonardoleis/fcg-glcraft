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
	PopulatedBlocks             []*game_objects.Block
	ShouldUpdatePopulatedBlocks bool
	Seed                        int64
	Time                        int64
	Tick                        uint
}

func NewWorld(worldName string, size mgl32.Vec3, seed int64) *World {

	return &World{
		Name: worldName,
		Size: size,
		Seed: seed,
		Time: 0,
	}
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
		} else if w.Blocks[wx][keys[wy]][wz] != nil && w.Blocks[wx][keys[wy]][wz].Position.Y() > highestBlock.Position.Y() {
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

var maxPrimeIndex = 10
var primeIndex = 0
var numX = 512
var numY = 512
var numOctaves = 7
var persistence = 0.5

var primes = [][]int{
	{995615039, 600173719, 701464987},
	{831731269, 162318869, 136250887},
	{174329291, 946737083, 245679977},
	{362489573, 795918041, 350777237},
	{457025711, 880830799, 909678923},
	{787070341, 177340217, 593320781},
	{405493717, 291031019, 391950901},
	{458904767, 676625681, 424452397},
	{531736441, 939683957, 810651871},
	{997169939, 842027887, 423882827},
}

func Noise(i, x, y int) float64 {
	n := x + y*57
	n = (n << 13) ^ n
	a := primes[i][0]
	b := primes[i][1]
	c := primes[i][2]
	t := (n*(n*n*a+b) + c) & 0x7fffffff
	return 1.0 - float64(t)/1073741824.0
}

func SmoothedNoise(i, x, y int) float64 {
	corners := (Noise(i, x-1, y-1) + Noise(i, x+1, y-1) +
		Noise(i, x-1, y+1) + Noise(i, x+1, y+1)) / 16
	sides := (Noise(i, x-1, y) + Noise(i, x+1, y) + Noise(i, x, y-1) +
		Noise(i, x, y+1)) / 8
	center := Noise(i, x, y) / 4
	return corners + sides + center
}

func Interpolate(a, b, x float64) float64 { // cosine interpolation
	ft := x * 3.1415927
	f := (1 - math.Cos(ft)) * 0.5
	return a*(1-f) + b*f
}

func InterpolatedNoise(i int, x, y float64) float64 {
	integer_X := x
	fractional_X := x - integer_X
	integer_Y := y
	fractional_Y := y - integer_Y

	v1 := SmoothedNoise(i, int(integer_X), int(integer_Y))
	v2 := SmoothedNoise(i, int(integer_X)+1, int(integer_Y))
	v3 := SmoothedNoise(i, int(integer_X), int(integer_Y)+1)
	v4 := SmoothedNoise(i, int(integer_X)+1, int(integer_Y)+1)
	i1 := Interpolate(v1, v2, fractional_X)
	i2 := Interpolate(v3, v4, fractional_X)
	return Interpolate(i1, i2, fractional_Y)
}

func ValueNoise_2D(x, y float64) float64 {
	total := 0.0
	frequency := math.Pow(2, float64(numOctaves))
	amplitude := 1.0
	for i := 0; i < numOctaves; i++ {
		frequency /= 2
		amplitude *= persistence
		total += InterpolatedNoise((primeIndex+i)%maxPrimeIndex,
			x/frequency, y/frequency) * amplitude
	}
	return total / frequency
}

func (w *World) GenerateWorld() {
	p := perlin.NewPerlin(2, 2, 10, time.Now().Unix())

	generatedWorld := make(WorldBlocks)
	for x := int(-w.Size.X()); x < int(w.Size.X()); x++ {
		if len(generatedWorld[x]) == 0 {
			generatedWorld[x] = make(map[int]map[int]*game_objects.Block)
		}
		for z := int(-w.Size.Z()); z < int(w.Size.Z()); z++ {

			y := float64(w.Size.Y()/2) + math.Round(float64(w.Size.Y()-w.Size.Y()/2)*p.Noise2D(float64(x)/float64(w.Size.X()), float64(z)/float64(w.Size.Z())))

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

func (w *World) UpdatePopulatedBlocks(fromX, toX, fromY, toY, fromZ, toZ float64) {

	w.PopulatedBlocks = []*game_objects.Block{}

	for x := fromX; x < toX; x++ {
		intX := int(x)
		for y := 0; y < int(w.Size.Y()); y++ {

			intY := int(y)
			for z := fromZ; z < toZ; z++ {
				intZ := int(z)

				if w.Blocks[intX][intY][intZ] == nil || w.Blocks[intX][intY][intZ].CountNeighbors() == 6 {
					continue
				}

				w.PopulatedBlocks = append(w.PopulatedBlocks, w.Blocks[intX][intY][intZ])

			}
		}
	}

	w.ShouldUpdatePopulatedBlocks = false
}

func (w *World) Update(roundedPlayerPosition mgl32.Vec3, backOfPlayer, frontOfPlayer mgl32.Vec3) {

	maxDist := float64(32)

	roundedPlayerX, roundedPlayerY, roundedPlayerZ := roundedPlayerPosition.Elem()
	fromX, toX := math.Max(-float64(w.Size.X()), float64(roundedPlayerX)-maxDist), math.Min(float64(w.Size.X()), float64(roundedPlayerX)+maxDist)
	fromY, toY := math.Max(-float64(w.Size.Y()), float64(roundedPlayerY)-maxDist), math.Min(float64(w.Size.Y()), float64(roundedPlayerY)+maxDist)
	fromZ, toZ := math.Max(-float64(w.Size.Z()), float64(roundedPlayerZ)-maxDist), math.Min(float64(w.Size.Z()), float64(roundedPlayerZ)+maxDist)

	// fmt.Println(len(w.PopulatedBlocks))

	sumX := frontOfPlayer.X() - backOfPlayer.X()
	sumZ := frontOfPlayer.Z() - backOfPlayer.Z()

	for _, populatedBlock := range w.PopulatedBlocks {
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
