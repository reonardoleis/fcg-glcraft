// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.
package main // import "github.com/go-gl/example/gl41core-cube"

import (
	"fmt"
	"go/build"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"math"
	"os"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/camera"
	"github.com/reonardoleis/fcg-glcraft/engine/controls"
	rendererPkg "github.com/reonardoleis/fcg-glcraft/engine/renderer"
	"github.com/reonardoleis/fcg-glcraft/engine/scene"
	"github.com/reonardoleis/fcg-glcraft/engine/shaders"
	"github.com/reonardoleis/fcg-glcraft/engine/window"
	"github.com/reonardoleis/fcg-glcraft/game_objects"
	math2 "github.com/reonardoleis/fcg-glcraft/math"
	"github.com/reonardoleis/fcg-glcraft/player"
	"github.com/reonardoleis/fcg-glcraft/world"
)

const windowWidth = 1280
const windowHeight = 720

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

func main() {

	window, err := window.NewWindow("fcg-glcraft", windowWidth, windowHeight)
	if err != nil {
		log.Fatal(err)
	}

	renderer := rendererPkg.NewRenderer()
	err = renderer.Init()
	if err != nil {
		log.Fatal(err)
	}

	//mainScene := scene.NewScene()

	_, err = shaders.InitShaderProgram("standard")
	if err != nil {
		panic(err)
	}

	_, err = shaders.InitShaderProgram2("crosshair")
	if err != nil {
		panic(err)
	}
	game_objects.InitBlock(1, 0.0, 1.0, 0.0)
	world := world.NewWorld("", mgl32.Vec3{100, 16, 100}, 2300932812397)
	world.GenerateWorld()

	dayTimeDirection := 1
	dayTimeColor := 1.0

	// model_uniform := gl.GetUniformLocation(program, gl.Str("model\000"))                     // Variável da matriz "model"
	// view_uniform := gl.GetUniformLocation(program, gl.Str("view\000"))             // Variável da matriz "view" em shader_vertex.glsl
	// projection_uniform := gl.GetUniformLocation(program, gl.Str("projection\000")) // Variável da matriz "projection" em shader_vertex.glsl
	// render_as_black_uniform := gl.GetUniformLocation(program, gl.Str("render_as_black\000")) // Variável booleana em shader_vertex.glsl

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)
	gl.FrontFace(gl.CCW)

	controlHandler := controls.NewControls(window)
	controlHandler.StartKeyHandlers()

	camera1 := camera.NewCamera(mgl32.Vec4{0.0, 0.0, 0.0, 1.0}, controlHandler, math.Pi/3, camera.FirstPersonCamera)
	player1 := player.NewPlayer(mgl32.Vec4{-1.0, 32, -6.0, 1.0}, controlHandler, 10, 2.0, 4, 10, 2)
	player1.BeFollowedByCamera(camera1)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	camera2 := camera.NewCamera(mgl32.Vec4{0.0, 0.0, 0.0, 1.0}, controlHandler, math.Pi/3, camera.FirstPersonCamera)
	player2 := player.NewPlayer(mgl32.Vec4{15.0, 30.0, -6.0, 1.0}, controlHandler, 10, 2.0, 4, 10, 2)
	player2.BeFollowedByCamera(camera2)

	sceneManager := scene.NewSceneManager()
	scene1 := scene.NewScene(world, camera1, &player1, controlHandler, scene.GameScene)
	scene2 := scene.NewScene(world, camera2, &player2, controlHandler, scene.GameScene)

	sceneManager.AddScene(scene1)
	sceneManager.AddScene(scene2)
	sceneManager.SetActiveScene(0)

	start := float64(0.0)
	end := float64(0.0)
	for !window.ShouldClose() {
		dayTimeColor += float64(dayTimeDirection) * 0.1 * math2.DeltaTime
		if dayTimeColor > 1.0 {
			dayTimeColor = 1.0
			dayTimeDirection = -1
		}
		if dayTimeColor < 0.0 {
			dayTimeColor = 0.0
			dayTimeDirection = 1
		}
		start = glfw.GetTime()

		thunder := math2.RandInt(0, 1000)
		if thunder >= 999 {
			gl.ClearColor(1.0, 1.0, 1.0, 1.0)
		} else {
			gl.ClearColor(0, float32(dayTimeColor), float32(dayTimeColor), 1.0)
		}

		if controlHandler.IsDown(int(glfw.KeyEnter)) {
			if sceneManager.ActiveScene == 0 {
				sceneManager.SetActiveScene(1)
			} else {
				sceneManager.SetActiveScene(0)
			}
		}

		sceneManager.HandleActiveScene(*window)

		end = glfw.GetTime()

		math2.DeltaTime = end - start
	}
}

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}

var vertexShader = `
#version 330
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
in vec3 vert;
in vec2 vertTexCoord;
out vec2 fragTexCoord;
void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
}
` + "\x00"

var fragmentShader = `
#version 330
uniform sampler2D tex;
in vec2 fragTexCoord;
out vec4 outputColor;
void main() {
    outputColor = texture(tex, fragTexCoord);
}
` + "\x00"

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	-1.0, -1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,

	// Top
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 1.0,

	// Left
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,

	// Right
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
}

// Set the working directory to the root of Go package, so that its assets can be accessed.
func init() {
	dir, err := importPathToDir("github.com/go-gl/example/gl41core-cube")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}

// importPathToDir resolves the absolute path from importPath.
// There doesn't need to be a valid Go package inside that import path,
// but the directory must exist.
func importPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}
