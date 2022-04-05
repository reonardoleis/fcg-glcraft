package geometry

import "unsafe"

type GeometryInformation struct {
	FirstIndex    unsafe.Pointer
	NumIndices    int
	RenderingMode int
	VaoID         uint32
	Vertexes      []float32
}
