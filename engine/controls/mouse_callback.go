package controls

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Mouse button callback
func MouseButtonCallback(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		gLastCursorPosX, gLastCursorPosY = window.GetCursorPos()
		keys[int(glfw.MouseButtonLeft)] = true
	}
	if button == glfw.MouseButtonLeft && action == glfw.Release {
		gLastCursorPosX, gLastCursorPosY = window.GetCursorPos()
		keys[int(glfw.MouseButtonLeft)] = false
	}

	if button == glfw.MouseButtonRight && action == glfw.Press {
		keys[int(glfw.MouseButtonRight)] = true
	}
	if button == glfw.MouseButtonRight && action == glfw.Release {
		keys[int(glfw.MouseButtonRight)] = false
	}
}

// Cursor pos callback as seen on classes
func CursorPosCallback(window *glfw.Window, x float64, y float64) {
	gCursorDeltaX = x - gLastCursorPosX
	gCursorDeltaY = y - gLastCursorPosY

	gLastCursorPosX = x
	gLastCursorPosY = y

	gMousePosChanged = true
}
