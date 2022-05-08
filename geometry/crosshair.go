package geometry

import "github.com/go-gl/gl/v3.3-core/gl"

// Draws crosshair on the screen
func DrawCrosshair() {
	crosshair := []float32{
		-0.02, 0,
		0.02, 0,
		0, -0.02,
		0, 0.02,
	}

	var vbo_2d, vao_2d uint32
	gl.GenBuffers(1, &vbo_2d)
	gl.GenVertexArrays(1, &vao_2d)

	gl.BindVertexArray(vao_2d)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo_2d)
	gl.BufferData(gl.ARRAY_BUFFER, len(crosshair)*4, gl.Ptr(crosshair), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 8, nil)
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(vao_2d)
	gl.DrawArrays(gl.LINES, 0, 4)
}
