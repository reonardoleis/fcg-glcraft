package shaders

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

var (
	ShaderProgramDefault   uint32
	ShaderProgramCrosshair uint32
)

func LoadFragmentShader(name string) (string, error) {
	_, filename, _, _ := runtime.Caller(0)
	fragmentShaderFile := fmt.Sprintf("%v/%v", path.Dir(filename), fmt.Sprintf("%v_shader_fragment.glsl", name))

	fragmentShader, err := os.ReadFile(fragmentShaderFile)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(fragmentShader) + "\x00", nil
}

func LoadVertexShader(name string) (string, error) {
	_, filename, _, _ := runtime.Caller(0)
	vertexShaderFile := fmt.Sprintf("%v/%v", path.Dir(filename), fmt.Sprintf("%v_shader_vertex.glsl", name))

	vertexShader, err := os.ReadFile(vertexShaderFile)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return string(vertexShader) + "\x00", nil
}

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
