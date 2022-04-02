package renderer

import (
	"log"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Renderer struct {
	available bool
}

func NewRenderer() *Renderer {
	return &Renderer{
		available: false,
	}
}

func (r *Renderer) Init() error {
	if err := gl.Init(); err != nil {
		log.Println(StrInitFail, err)
		return err
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("GL initialized with OpenGL version", version)

	r.available = true
	return nil
}
