package window

import (
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	ScreenRatio float32 = 1.0
)

func NewWindow(title string, width, height int) (*glfw.Window, error) {
	err := glfw.Init()
	if err != nil {
		log.Println(StrNewWindowFail, err)
		return nil, err
	}

	ScreenRatio = float32(width) / float32(height)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Samples, 32)

	window, _ := glfw.CreateWindow(width, height, title, nil, nil)
	window.MakeContextCurrent()

	return window, nil
}
