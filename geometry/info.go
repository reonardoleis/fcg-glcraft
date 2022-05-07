package geometry

import (
	"unsafe"

	"github.com/go-gl/mathgl/mgl32"
)

type GeometryInformation struct {
	FirstIndex    unsafe.Pointer
	NumIndices    int
	RenderingMode int
	VaoID         uint32
	Vertexes      []float32
	Position      mgl32.Vec3 // opcional, vai ser usado nos OBJs e no braço do jogador
}
