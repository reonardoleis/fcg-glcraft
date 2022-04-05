package math2

import "github.com/go-gl/mathgl/mgl32"

func North(position mgl32.Vec4, size float32) mgl32.Vec4 {
	return mgl32.Vec4{position.X() + size/2, position.Y(), position.Z(), position.W()}
}

func South(position mgl32.Vec4, size float32) mgl32.Vec4 {
	return mgl32.Vec4{position.X() - size/2, position.Y(), position.Z(), position.W()}
}

func East(position mgl32.Vec4, size float32) mgl32.Vec4 {
	return mgl32.Vec4{position.X(), position.Y(), position.Z() + size/2, position.W()}
}

func West(position mgl32.Vec4, size float32) mgl32.Vec4 {
	return mgl32.Vec4{position.X(), position.Y(), position.Z() - size/2, position.W()}
}

func Upper(position mgl32.Vec4, size float32) mgl32.Vec4 {
	return mgl32.Vec4{position.X(), position.Y() + size/2, position.Z(), position.W()}
}

func Lower(position mgl32.Vec4, size float32) mgl32.Vec4 {
	return mgl32.Vec4{position.X(), position.Y() - size/2, position.Z(), position.W()}
}
