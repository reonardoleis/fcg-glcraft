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

	mainScene := scene.NewScene()

	program, err := shaders.InitShaderProgram("standard")
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	cubeInformation := make(map[int]map[int]map[int]*game_objects.GameObject)

	for x := -world.WorldSizeX; x < world.WorldSizeX; x++ {
		cubeInformation[x] = make(map[int]map[int]*game_objects.GameObject)
		for y := -world.WorldSizeY; y < world.WorldSizeY; y++ {
			cubeInformation[x][y] = make(map[int]*game_objects.GameObject)
			for z := -world.WorldSizeZ; z < world.WorldSizeZ; z++ {
				if z == 2 {
					continue
				}
				newCube := game_objects.NewCube(float32(x), float32(y), float32(z), 1, true)
				mainScene.Add(newCube)
				cubeInformation[x][y][z] = &newCube
			}
		}
	}

	// model_uniform := gl.GetUniformLocation(program, gl.Str("model\000"))                     // Variável da matriz "model"
	// view_uniform := gl.GetUniformLocation(program, gl.Str("view\000"))             // Variável da matriz "view" em shader_vertex.glsl
	// projection_uniform := gl.GetUniformLocation(program, gl.Str("projection\000")) // Variável da matriz "projection" em shader_vertex.glsl
	// render_as_black_uniform := gl.GetUniformLocation(program, gl.Str("render_as_black\000")) // Variável booleana em shader_vertex.glsl

	gl.Enable(gl.DEPTH_TEST)

	controlHandler := controls.NewControls(window)
	controlHandler.StartKeyHandlers()

	camera := camera.NewCamera(mgl32.Vec4{0.0, 0.0, 0.0, 1.0}, controlHandler, math.Pi/3, camera.FirstPersonCamera)
	player := player.NewPlayer(mgl32.Vec4{-1.0, 30.0, -6.0, 1.0}, controlHandler, 10, 2.0, 4, 10, 2)
	player.SetCamera(camera)

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	start := float64(0.0)
	end := float64(0.0)
	for !window.ShouldClose() {
		start = glfw.GetTime()
		gl.ClearColor(0, 1, 1, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		gl.UseProgram(program)

		camera.Update()
		player.Update()

		if controlHandler.IsToggled(int(glfw.KeyZ)) {
			game_objects.CubeEdgesOnly = true
		} else {
			game_objects.CubeEdgesOnly = false
		}

		// gl.UniformMatrix4fv(view_uniform, 1, false, &view[0])
		// gl.UniformMatrix4fv(projection_uniform, 1, false, &projection[0])

		maxDist := float64(50)

		roundedPlayerX, _, roundedPlayerZ := player.GetRoundedPosition()
		playerY := float64(player.Position.Y())

		for x := math.Max(-float64(world.WorldSizeX), float64(roundedPlayerX)-maxDist); x < math.Min(float64(world.WorldSizeX), float64(roundedPlayerX)+maxDist); x++ {
			for y := -world.WorldSizeY; y < world.WorldSizeY; y++ {
				for z := math.Max(-float64(world.WorldSizeZ), float64(roundedPlayerZ)-maxDist); z < math.Min(float64(world.WorldSizeZ), float64(roundedPlayerZ)+maxDist); z++ {
					if cubeInformation[int(x)][y][int(z)] == nil {
						continue
					}
					cubeInformation[int(x)][y][int(z)].Draw()
				}
			}
		}

		window.SetTitle(fmt.Sprintf("X: %v - Y: %v - Z: %v - wsX: %v - wsZ: %v", roundedPlayerX, playerY, roundedPlayerZ, world.WorldSizeX, world.WorldSizeZ))

		//	fmt.Println(roundedPlayerX, worldSizeX)
		//	fmt.Println(roundedPlayerZ, worldSizeZ)

		firstTime := true
		highest := 0
		for k := range cubeInformation[roundedPlayerX] {
			if firstTime {
				highest = k
				firstTime = false
			} else if k > highest {
				highest = k
			}
		}

		blockBelow := cubeInformation[roundedPlayerX][highest][roundedPlayerZ]

		player.Fall(blockBelow)

		controlHandler.FinishMousePositionChanged()
		window.SwapBuffers()
		glfw.PollEvents()
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
