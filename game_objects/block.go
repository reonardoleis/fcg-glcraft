package game_objects

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"os"
	"path"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/engine/shaders"
	"github.com/reonardoleis/fcg-glcraft/geometry"
	math2 "github.com/reonardoleis/fcg-glcraft/math"
)

var (
	BlockEdgesOnly bool = false
)

type BlockType = uint

const (
	BlockGrass = iota
	BlockDirt
	BlockWood
	BlockLeaves
	BlockSand
	BlockStone
	BlockWater
)

var (
	cubeVaoID         geometry.GeometryInformation
	cubeModelMatrix          = math2.Matrix_Identity()
	dirtTexture       uint32 = 0
	grassSideTexture  uint32 = 0
	grassTopTexture   uint32 = 0
	woodSideTexture   uint32 = 0
	leavesTexture     uint32 = 0
	stoneTexture      uint32 = 0
	waterTexture      uint32 = 0
	sandTexture       uint32 = 0
	numTexturesLoaded        = 0
	model_uniform     int32
)

type Block struct {
	Position mgl32.Vec4
	Size     float32
	Model    mgl32.Mat4
	// ModelGeometry geometry.GeometryInformation
	WithEdges bool
	// EdgesGeometry geometry.GeometryInformation
	Ephemeral bool

	BlockType

	Neighbors [6]bool
}

func GetBlockTypes() []BlockType {
	return []BlockType{
		BlockGrass,
		BlockDirt,
		BlockWood,
		BlockLeaves,
		BlockSand,
		BlockStone,
		BlockWater,
	}
}

func getBlockTexture(blockType BlockType) []uint32 {
	switch blockType {
	case BlockDirt:
		return []uint32{dirtTexture, dirtTexture, dirtTexture, dirtTexture, dirtTexture, dirtTexture}
	case BlockGrass:
		return []uint32{grassSideTexture, grassSideTexture, grassSideTexture, grassSideTexture, grassTopTexture, dirtTexture}
	case BlockWood:
		return []uint32{woodSideTexture, woodSideTexture, woodSideTexture, woodSideTexture, woodSideTexture, woodSideTexture}
	case BlockLeaves:
		return []uint32{leavesTexture, leavesTexture, leavesTexture, leavesTexture, leavesTexture, leavesTexture}
	case BlockStone:
		return []uint32{stoneTexture, stoneTexture, stoneTexture, stoneTexture, stoneTexture, stoneTexture}
	case BlockWater:
		return []uint32{waterTexture, waterTexture, waterTexture, waterTexture, waterTexture, waterTexture}
	case BlockSand:
		return []uint32{sandTexture, sandTexture, sandTexture, sandTexture, sandTexture, sandTexture}
	}

	return []uint32{grassSideTexture, grassSideTexture, grassSideTexture, grassSideTexture, grassTopTexture, dirtTexture}
}

type BlockTexture = mgl32.Vec3

func getTexture(blockType BlockType) BlockTexture {
	switch blockType {
	case BlockGrass:
		return BlockTexture{0.0, 0.7, 0.0}
	case BlockDirt:
		return BlockTexture{0.7, 0.5, 0.3}
	case BlockWood:
		return BlockTexture{0.3, 0.25, 0.15}
	case BlockLeaves:
		return BlockTexture{0.0, 1.0, 0.0}
	case BlockSand:
		return BlockTexture{1.0, 0.9, 0.5}
	}

	return BlockTexture{0.0, 0.0, 0.0}
}

func InitBlock() {
	dirtTexture = newTexture("dirt_0.png")
	grassTopTexture = newTexture("grass_0.png")
	grassSideTexture = newTexture("grass_1.png")
	woodSideTexture = newTexture("wood_0.png")
	leavesTexture = newTexture("leaves_0.png")
	stoneTexture = newTexture("stone_0.png")
	waterTexture = newTexture("water_0.png")
	sandTexture = newTexture("sand_0.png")
}

func NewBlock(x, y, z, size float32, withEdges, ephemeral bool, blockType BlockType) Block {
	// modelGeometry := cubeVaoID
	// edgesGeometry := geometry.GeometryInformation{}

	model := math2.Matrix_Identity().Mul4(math2.Matrix_Translate(x, y, z))

	if withEdges {
		// edgesGeometry = geometry.BuildCubeEdges(0, 0, 0, size)
	}
	return Block{
		Position: mgl32.Vec4{x, y, z, 0.0},
		Size:     size,
		Model:    model,
		// ModelGeometry: modelGeometry,
		WithEdges: false,
		// EdgesGeometry: edgesGeometry,
		Ephemeral: ephemeral,
		BlockType: blockType,
	}
}

func (c *Block) Translate(x, y, z float32) {
	c.Model = math2.Matrix_Identity().Mul4(math2.Matrix_Translate(x, y, z))
}

func newTexture(file string) uint32 {
	_, filename, _, _ := runtime.Caller(0)
	textureFile := fmt.Sprintf("%v/%v", path.Dir(filename), file)
	imgFile, err := os.Open(textureFile)
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		panic(err)
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0 + uint32(numTexturesLoaded))
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	numTexturesLoaded++
	return texture
}

func (b Block) CountNeighbors() int {
	count := 0
	for _, neighbor := range b.Neighbors {
		if neighbor {
			count++
		}
	}

	return count
}

func (b Block) Draw2() {

	model_uniform := gl.GetUniformLocation(shaders.ShaderProgramDefault, gl.Str("model\000")) // Variável da matriz "model"
	black := gl.GetUniformLocation(shaders.ShaderProgramDefault, gl.Str("black\000"))         // Variável da matriz "model"
	blockTextures := getBlockTexture(b.BlockType)
	north := math2.North(b.Position, b.Size)
	south := math2.South(b.Position, b.Size)
	east := math2.East(b.Position, b.Size)
	west := math2.West(b.Position, b.Size)
	upper := math2.Upper(b.Position, b.Size)
	lower := math2.Lower(b.Position, b.Size)
	faces := []mgl32.Vec4{north, south, east, west, upper, lower}
	northRotation := math2.Matrix_Rotate_Y((math.Pi / 180) * 90)
	southRotation := math2.Matrix_Rotate_Y((math.Pi / 180) * 90)
	eastRotation := math2.Matrix_Identity()
	westRotation := math2.Matrix_Identity()
	upperRotation := math2.Matrix_Rotate_X((math.Pi / 180) * 90)
	lowerRotation := math2.Matrix_Rotate_X((math.Pi / 180) * 90)

	rotations := []mgl32.Mat4{northRotation, southRotation, eastRotation, westRotation, upperRotation, lowerRotation}

	diff := 0.0
	if b.BlockType == BlockWater {
		diff = 0.2
	}
	gl.BindVertexArray(geometry.Faces[b.BlockType].VaoID)
	for index, face := range faces {
		if b.Neighbors[index] {
			continue
		}
		if !BlockEdgesOnly {
			gl.ActiveTexture(gl.TEXTURE0)
			gl.BindTexture(gl.TEXTURE_2D, blockTextures[index])

			faceMat := math2.Matrix_Identity().Mul4(math2.Matrix_Translate(face.X(), face.Y()-float32(diff), face.Z())).Mul4(rotations[index])

			gl.UniformMatrix4fv(model_uniform, 1, false, &faceMat[0])
			gl.Uniform1i(black, 0)
			gl.DrawElements(
				uint32(geometry.Faces[b.BlockType].RenderingMode), // Veja slides 124-130 do documento Aula_04_Modelagem_Geometrica_3D.pdf
				int32(geometry.Faces[b.BlockType].NumIndices),
				gl.UNSIGNED_INT,
				geometry.Faces[b.BlockType].FirstIndex,
			)
			if b.WithEdges {

				faceMat := math2.Matrix_Identity().Mul4(math2.Matrix_Translate(face.X(), face.Y(), face.Z())).Mul4(rotations[index])
				gl.BindVertexArray(geometry.CommonFaceEdgeGeometry.VaoID)
				gl.UniformMatrix4fv(model_uniform, 1, false, &faceMat[0])
				gl.Uniform1i(black, 1)
				gl.DrawElements(
					uint32(geometry.CommonFaceEdgeGeometry.RenderingMode), // Veja slides 124-130 do documento Aula_04_Modelagem_Geometrica_3D.pdf
					int32(geometry.CommonFaceEdgeGeometry.NumIndices),
					gl.UNSIGNED_INT,
					geometry.CommonFaceEdgeGeometry.FirstIndex,
				)
			}
		} else {
			faceMat := math2.Matrix_Identity().Mul4(math2.Matrix_Translate(face.X(), face.Y(), face.Z())).Mul4(rotations[index])
			gl.BindVertexArray(geometry.CommonFaceEdgeGeometry.VaoID)
			gl.UniformMatrix4fv(model_uniform, 1, false, &faceMat[0])
			// gl.Uniform1i(render_as_black_uniform, 1)
			gl.DrawElements(
				uint32(geometry.CommonFaceEdgeGeometry.RenderingMode), // Veja slides 124-130 do documento Aula_04_Modelagem_Geometrica_3D.pdf
				int32(geometry.CommonFaceEdgeGeometry.NumIndices),
				gl.UNSIGNED_INT,
				geometry.CommonFaceEdgeGeometry.FirstIndex,
			)
		}

	}

}

func (b Block) Draw() {
	model_uniform := gl.GetUniformLocation(shaders.ShaderProgramDefault, gl.Str("model\000"))                     // Variável da matriz "model"
	render_as_black_uniform := gl.GetUniformLocation(shaders.ShaderProgramDefault, gl.Str("render_as_black\000")) // Variável booleana em shader_vertex.glsl

	if !BlockEdgesOnly && !b.Ephemeral {
		gl.BindVertexArray(cubeVaoID.VaoID)
		gl.UniformMatrix4fv(model_uniform, 1, false, &b.Model[0])
		gl.Uniform1i(render_as_black_uniform, 0)
		gl.DrawElements(
			uint32(cubeVaoID.RenderingMode), // Veja slides 124-130 do documento Aula_04_Modelagem_Geometrica_3D.pdf
			int32(cubeVaoID.NumIndices),
			gl.UNSIGNED_INT,
			cubeVaoID.FirstIndex,
		)
	}

	/*if b.WithEdges || b.Ephemeral {
		gl.BindVertexArray(b.EdgesGeometry.VaoID)
		gl.UniformMatrix4fv(model_uniform, 1, false, &b.Model[0])
		gl.Uniform1i(render_as_black_uniform, 1)
		gl.DrawElements(
			uint32(b.EdgesGeometry.RenderingMode), // Veja slides 124-130 do documento Aula_04_Modelagem_Geometrica_3D.pdf
			int32(b.EdgesGeometry.NumIndices),
			gl.UNSIGNED_INT,
			b.EdgesGeometry.FirstIndex,
		)
	}*/
}

func (b Block) GetPosition() mgl32.Vec4 {
	return b.Position
}
