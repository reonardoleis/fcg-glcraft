package entrypoint

import (
	"log"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type EntryPoint struct {
	available bool
}

func NewEntryPoint() *EntryPoint {
	return &EntryPoint{
		available: false,
	}
}

// starts opengl
func (r *EntryPoint) Init() error {
	if err := gl.Init(); err != nil {
		log.Println(StrInitFail, err)
		return err
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("GL initialized with OpenGL version", version)

	r.available = true
	return nil
}
