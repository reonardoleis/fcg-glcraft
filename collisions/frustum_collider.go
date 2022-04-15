package collisions

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/camera"
)

type FrustumCollider struct {
	Frustum camera.Frustum
}

func NewFrustumCollider(frustum camera.Frustum) *FrustumCollider {
	return &FrustumCollider{
		Frustum: frustum,
	}
}

func (f *FrustumCollider) UpdateFrustum(frustum camera.Frustum) {
	f.Frustum = frustum
}

func getSmallest(values []float32) float32 {
	smallest := values[0]
	for _, item := range values {
		if item < smallest {
			smallest = item
		}
	}

	return smallest
}

func getBiggest(values []float32) float32 {
	biggest := values[0]
	for _, item := range values {
		if item > biggest {
			biggest = item
		}
	}

	return biggest
}

func (f FrustumCollider) CollidesWithBlock(p mgl32.Vec3) bool {
	ft := f.Frustum
	xs := []float32{ft.Ntl.X(), ft.Ntr.X(), ft.Ftl.X(), ft.Ftr.X(), ft.Nbl.X(), ft.Nbr.X(), ft.Fbr.X(), ft.Fbl.X()}
	ys := []float32{ft.Ntl.Y(), ft.Ntr.Y(), ft.Ftl.Y(), ft.Ftr.Y(), ft.Nbl.Y(), ft.Nbr.Y(), ft.Fbr.Y(), ft.Fbl.Y()}
	zs := []float32{ft.Ntl.Z(), ft.Ntr.Z(), ft.Ftl.Z(), ft.Ftr.Z(), ft.Nbl.Z(), ft.Nbr.Z(), ft.Fbr.Z(), ft.Fbl.Z()}

	minX := getSmallest(xs)
	maxX := getBiggest(xs)

	minY := getSmallest(ys)
	maxY := getBiggest(ys)

	minZ := getSmallest(zs)
	maxZ := getBiggest(zs)

	return (p.X() >= float32(minX) && p.X() <= float32(maxX)) &&
		(p.Y() >= float32(minY) && p.Y() <= float32(maxY)) &&
		(p.Z() >= float32(minZ) && p.Z() <= float32(maxZ))

}
