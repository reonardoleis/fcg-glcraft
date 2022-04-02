package renderer

import "unsafe"

type SceneObject struct {
	Name          string
	FirstIndex    unsafe.Pointer
	NumIndices    int
	RenderingMode int
	VaoID         uint32
}
