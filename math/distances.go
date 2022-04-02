package math2

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func Distance(u, v mgl32.Vec4) float64 {
	return math.Sqrt(math.Pow((float64(u.X())-float64(v.X())), 2) + math.Pow((float64(u.Y())-float64(v.Y())), 2) + math.Pow((float64(u.Z())-float64(v.Z())), 2))
}
