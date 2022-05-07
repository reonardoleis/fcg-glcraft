package geometry

import (
	"unsafe"

	"github.com/go-gl/mathgl/mgl32"
	math2 "github.com/reonardoleis/fcg-glcraft/math"
)

type GeometryInformation struct {
	FirstIndex    unsafe.Pointer
	NumIndices    int
	RenderingMode int
	VaoID         uint32
	Vertexes      []float32
	Position      mgl32.Vec3        // opcional, vai ser usado nos OBJs e no braço do jogador
	T             float32           // opcional, sera usado para animacoes com curvas de bezier
	Tdir          float32           // opcional sera usado para animacoes com curvas de bezier
	BCurve        math2.BezierCurve // opcional sera usado se tiver curva de bezier para animacao
	Animating     bool              // opcional, se o objeto for animado indicará se deve "tocar" a animação ou nãoäää
}
