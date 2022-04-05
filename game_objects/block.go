package game_objects

import (
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
)

var (
	cubeVaoID geometry.GeometryInformation
)

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

type Block struct {
	Position mgl32.Vec4
	Size     float32
	Model    mgl32.Mat4
	// ModelGeometry geometry.GeometryInformation
	WithEdges bool
	// EdgesGeometry geometry.GeometryInformation
	Ephemeral bool

	BlockType
}

func InitBlock(size, r, g, b float32) {
	cubeVaoID = geometry.BuildCube(0, 0, 0, size, r, g, b)
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
	}
}

func (c *Block) Translate(x, y, z float32) {
	c.Model = math2.Matrix_Identity().Mul4(math2.Matrix_Translate(x, y, z))
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
