package controls

import "github.com/go-gl/glfw/v3.3/glfw"

const (
	LMB = iota
	RMB
)

var (
	keys         = make(map[int]bool)
	keysToggling = make(map[int]int)
)

var (
	gLastCursorPosX  float64 = 0
	gLastCursorPosY  float64 = 0
	gCursorDeltaX    float64 = 0
	gCursorDeltaY    float64 = 0
	gMousePosChanged bool    = false
)

type Controls struct {
	window *glfw.Window
}

func boolFromInt(n int) bool {
	if n == 0 {
		return false
	}

	return true
}

func NewControls(window *glfw.Window) Controls {
	return Controls{
		window: window,
	}
}

func (c Controls) IsDown(key int) bool {
	return keys[key]
}

func (c Controls) IsToggled(key int) bool {
	return boolFromInt(keysToggling[key])
}

func (c Controls) SetKeyStatus(key int, status bool) {
	keys[key] = status
}

func (c Controls) StartKeyHandlers() {
	c.window.SetMouseButtonCallback(MouseButtonCallback)
	c.window.SetCursorPosCallback(CursorPosCallback)
	c.window.SetKeyCallback(KeyCallback)
}

func (c Controls) GetMouseDeltas() (float64, float64) {
	return gCursorDeltaX, gCursorDeltaY
}

func (c Controls) MousePositionChanged() bool {
	defer func() {
		gMousePosChanged = false
	}()
	return gMousePosChanged
}
