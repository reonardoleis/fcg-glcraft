package world

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/game_objects"
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

	newCube := game_objects.NewBlock(float32(x), float32(y), float32(z), 1, true, ephemeral, game_objects.BlockGrass)
	newCube.WithEdges = false
	w.Blocks[int(x)][int(y)][int(z)] = &newCube

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
		fmt.Println(bb)
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

	generatedWorld := make(WorldBlocks)
	for x := int(-w.Size.X()); x < int(w.Size.X()); x++ {
		generatedWorld[x] = make(map[int]map[int]*game_objects.Block)
		for z := int(-w.Size.Z()); z < int(w.Size.Z()); z++ {

			// fmt.Println(y)

			blockTypeNoise := ValueNoise_2D(float64(x), float64(z))
			blockType := game_objects.BlockDirt
			if blockTypeNoise >= 0 && blockTypeNoise <= 0.06 {
				blockType = game_objects.BlockDirt
			} else if blockTypeNoise > 0.06 && blockTypeNoise <= 0.1 {
				blockType = game_objects.BlockGrass
			} else {
				blockType = game_objects.BlockSand
			}
			y := 4 + math.Round(float64(w.Size.Y()-4)*ValueNoise_2D(float64(x), float64(z)))

			for i := int(y); i >= int(math.Max(0, float64(int(y)-4))); i-- {
				if generatedWorld[x][i] == nil {
					generatedWorld[x][i] = make(map[int]*game_objects.Block)
				}
				newCube := game_objects.NewBlock(float32(x), float32(i), float32(z), 1, true, false, uint(blockType))
				newCube.WithEdges = false
				fmt.Println(unsafe.Sizeof(newCube))
				generatedWorld[x][i][z] = &newCube
			}

		}

	}

	w.Blocks = generatedWorld
	w.Size = mgl32.Vec3{w.Size.X(), 128, w.Size.Z()}
}

func (w *World) Update(roundedPlayerPosition mgl32.Vec3) {
	maxDist := float64(25)

	roundedPlayerX, roundedPlayerY, roundedPlayerZ := roundedPlayerPosition.Elem()

	for x := math.Max(-float64(w.Size.X()), float64(roundedPlayerX)-maxDist); x < math.Min(float64(w.Size.X()), float64(roundedPlayerX)+maxDist); x++ {
		for y := math.Max(-float64(w.Size.Y()), float64(roundedPlayerY)-maxDist); y < math.Min(float64(w.Size.Y()), float64(roundedPlayerY)+maxDist); y++ {
			for z := math.Max(-float64(w.Size.Z()), float64(roundedPlayerZ)-maxDist); z < math.Min(float64(w.Size.Z()), float64(roundedPlayerZ)+maxDist); z++ {
				if w.Blocks[int(x)][int(y)][int(z)] == nil {
					continue
				}

				w.Blocks[int(x)][int(y)][int(z)].Draw()
				w.Blocks[int(x)][int(y)][int(z)].WithEdges = false
				if w.Blocks[int(x)][int(y)][int(z)].Ephemeral {
					w.Blocks[int(x)][int(y)][int(z)] = nil
				}

			}
		}
	}
}
