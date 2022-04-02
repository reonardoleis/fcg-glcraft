package controls

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

func KeyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}

	keys[int(key)] = boolFromInt(int(action))

	if action == glfw.Press {
		current := keysToggling[int(key)]
		if current == 0 {
			keysToggling[int(key)] = 1
		} else {
			keysToggling[int(key)] = 0
		}
	}

}
