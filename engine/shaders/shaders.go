package shaders

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

var (
	ShaderProgramDefault   uint32
	ShaderProgramCrosshair uint32
	FragmentShader         string
	VertexShader           string
)

// Loads a fragment shader
func LoadFragmentShader(name string) (string, error) {
	fragmentShaderFile := fmt.Sprintf(fmt.Sprintf("./%v_shader_fragment.glsl", name))

	fragmentShader, err := os.ReadFile(fragmentShaderFile)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(fragmentShader) + "\x00", nil
}

// Loads a vertex shader
func LoadVertexShader(name string) (string, error) {
	vertexShaderFile := fmt.Sprintf(fmt.Sprintf("./%v_shader_vertex.glsl", name))

	vertexShader, err := os.ReadFile(vertexShaderFile)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(vertexShader) + "\x00", nil
}

// compile standard shaders
func InitStandardShaderPrograms(vertex, frag string) (uint32, error) {
	vertexShaderSource, err := LoadVertexShader(vertex)
	if err != nil {
		return 0, err
	}

	fragmentShaderSource, err := LoadFragmentShader(frag)
	if err != nil {
		return 0, err
	}

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	ShaderProgramDefault = program
	return program, nil
}

// compile shader
func InitShaderProgram(name string) (uint32, error) {
	vertexShaderSource, err := LoadVertexShader(name)
	if err != nil {
		return 0, err
	}

	fragmentShaderSource, err := LoadFragmentShader(name)
	if err != nil {
		return 0, err
	}

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	ShaderProgramDefault = program
	return program, nil
}

// compile alternative shader program
func InitShaderProgram2(name string) (uint32, error) {
	vertexShaderSource, err := LoadVertexShader(name)
	if err != nil {
		return 0, err
	}

	fragmentShaderSource, err := LoadFragmentShader(name)
	if err != nil {
		return 0, err
	}

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	ShaderProgramCrosshair = program
	return program, nil
}

// compiles the shader with gl calls
func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
