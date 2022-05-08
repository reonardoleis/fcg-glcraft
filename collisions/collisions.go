package collisions

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/reonardoleis/fcg-glcraft/camera"
	math2 "github.com/reonardoleis/fcg-glcraft/math"
)

type CubeCollider struct{}

func NewCubeCollider() *CubeCollider {
	return &CubeCollider{}
}

func (cc CubeCollider) CollidesABBABB(a, b CubeBoundingBox, av, bv [8]mgl32.Vec3) bool {
	// ABB-ABB Collision
	if av[0].X() >= b.Mins.X() && av[0].X() <= b.Maxes.X() && av[0].Y() >= b.Mins.Y() && av[0].Y() <= b.Maxes.Y() && av[0].Z() >= b.Mins.Z() && av[0].Z() <= b.Maxes.Z() ||
		av[1].X() >= b.Mins.X() && av[1].X() <= b.Maxes.X() && av[1].Y() >= b.Mins.Y() && av[1].Y() <= b.Maxes.Y() && av[1].Z() >= b.Mins.Z() && av[1].Z() <= b.Maxes.Z() ||
		av[2].X() >= b.Mins.X() && av[2].X() <= b.Maxes.X() && av[2].Y() >= b.Mins.Y() && av[2].Y() <= b.Maxes.Y() && av[2].Z() >= b.Mins.Z() && av[2].Z() <= b.Maxes.Z() ||
		av[3].X() >= b.Mins.X() && av[3].X() <= b.Maxes.X() && av[3].Y() >= b.Mins.Y() && av[3].Y() <= b.Maxes.Y() && av[3].Z() >= b.Mins.Z() && av[3].Z() <= b.Maxes.Z() ||
		av[4].X() >= b.Mins.X() && av[4].X() <= b.Maxes.X() && av[4].Y() >= b.Mins.Y() && av[4].Y() <= b.Maxes.Y() && av[4].Z() >= b.Mins.Z() && av[4].Z() <= b.Maxes.Z() ||
		av[5].X() >= b.Mins.X() && av[5].X() <= b.Maxes.X() && av[5].Y() >= b.Mins.Y() && av[5].Y() <= b.Maxes.Y() && av[5].Z() >= b.Mins.Z() && av[5].Z() <= b.Maxes.Z() ||
		av[6].X() >= b.Mins.X() && av[6].X() <= b.Maxes.X() && av[6].Y() >= b.Mins.Y() && av[6].Y() <= b.Maxes.Y() && av[6].Z() >= b.Mins.Z() && av[6].Z() <= b.Maxes.Z() ||
		av[7].X() >= b.Mins.X() && av[7].X() <= b.Maxes.X() && av[7].Y() >= b.Mins.Y() && av[7].Y() <= b.Maxes.Y() && av[7].Z() >= b.Mins.Z() && av[7].Z() <= b.Maxes.Z() ||
		bv[0].X() >= a.Mins.X() && bv[0].X() <= a.Maxes.X() && bv[0].Y() >= a.Mins.Y() && bv[0].Y() <= a.Maxes.Y() && bv[0].Z() >= a.Mins.Z() && bv[0].Z() <= a.Maxes.Z() ||
		bv[1].X() >= a.Mins.X() && bv[1].X() <= a.Maxes.X() && bv[1].Y() >= a.Mins.Y() && bv[1].Y() <= a.Maxes.Y() && bv[1].Z() >= a.Mins.Z() && bv[1].Z() <= a.Maxes.Z() ||
		bv[2].X() >= a.Mins.X() && bv[2].X() <= a.Maxes.X() && bv[2].Y() >= a.Mins.Y() && bv[2].Y() <= a.Maxes.Y() && bv[2].Z() >= a.Mins.Z() && bv[2].Z() <= a.Maxes.Z() ||
		bv[3].X() >= a.Mins.X() && bv[3].X() <= a.Maxes.X() && bv[3].Y() >= a.Mins.Y() && bv[3].Y() <= a.Maxes.Y() && bv[3].Z() >= a.Mins.Z() && bv[3].Z() <= a.Maxes.Z() ||
		bv[4].X() >= a.Mins.X() && bv[4].X() <= a.Maxes.X() && bv[4].Y() >= a.Mins.Y() && bv[4].Y() <= a.Maxes.Y() && bv[4].Z() >= a.Mins.Z() && bv[4].Z() <= a.Maxes.Z() ||
		bv[5].X() >= a.Mins.X() && bv[5].X() <= a.Maxes.X() && bv[5].Y() >= a.Mins.Y() && bv[5].Y() <= a.Maxes.Y() && bv[5].Z() >= a.Mins.Z() && bv[5].Z() <= a.Maxes.Z() ||
		bv[6].X() >= a.Mins.X() && bv[6].X() <= a.Maxes.X() && bv[6].Y() >= a.Mins.Y() && bv[6].Y() <= a.Maxes.Y() && bv[6].Z() >= a.Mins.Z() && bv[6].Z() <= a.Maxes.Z() ||
		bv[7].X() >= a.Mins.X() && bv[7].X() <= a.Maxes.X() && bv[7].Y() >= a.Mins.Y() && bv[7].Y() <= a.Maxes.Y() && bv[7].Z() >= a.Mins.Z() && bv[7].Z() <= a.Maxes.Z() ||

		av[0].X() >= b.Mins.X() && av[0].X() <= b.Maxes.X() && av[0].Y()+0.5 >= b.Mins.Y() && av[0].Y()+0.5 <= b.Maxes.Y() && av[0].Z() >= b.Mins.Z() && av[0].Z() <= b.Maxes.Z() ||
		av[1].X() >= b.Mins.X() && av[1].X() <= b.Maxes.X() && av[1].Y()+0.5 >= b.Mins.Y() && av[1].Y()+0.5 <= b.Maxes.Y() && av[1].Z() >= b.Mins.Z() && av[1].Z() <= b.Maxes.Z() ||
		av[2].X() >= b.Mins.X() && av[2].X() <= b.Maxes.X() && av[2].Y()+0.5 >= b.Mins.Y() && av[2].Y()+0.5 <= b.Maxes.Y() && av[2].Z() >= b.Mins.Z() && av[2].Z() <= b.Maxes.Z() ||
		av[3].X() >= b.Mins.X() && av[3].X() <= b.Maxes.X() && av[3].Y()+0.5 >= b.Mins.Y() && av[3].Y()+0.5 <= b.Maxes.Y() && av[3].Z() >= b.Mins.Z() && av[3].Z() <= b.Maxes.Z() ||
		av[0].X() >= b.Mins.X() && av[0].X() <= b.Maxes.X() && av[0].Y()+1.0 >= b.Mins.Y() && av[0].Y()+1.0 <= b.Maxes.Y() && av[0].Z() >= b.Mins.Z() && av[0].Z() <= b.Maxes.Z() ||
		av[1].X() >= b.Mins.X() && av[1].X() <= b.Maxes.X() && av[1].Y()+1.0 >= b.Mins.Y() && av[1].Y()+1.0 <= b.Maxes.Y() && av[1].Z() >= b.Mins.Z() && av[1].Z() <= b.Maxes.Z() ||
		av[2].X() >= b.Mins.X() && av[2].X() <= b.Maxes.X() && av[2].Y()+1.0 >= b.Mins.Y() && av[2].Y()+1.0 <= b.Maxes.Y() && av[2].Z() >= b.Mins.Z() && av[2].Z() <= b.Maxes.Z() ||
		av[3].X() >= b.Mins.X() && av[3].X() <= b.Maxes.X() && av[3].Y()+1.0 >= b.Mins.Y() && av[3].Y()+1.0 <= b.Maxes.Y() && av[3].Z() >= b.Mins.Z() && av[3].Z() <= b.Maxes.Z() {
		return true
	}

	return false

}

type CubeBoundingBox struct {
	Maxes  mgl32.Vec3
	Mins   mgl32.Vec3
	width  float32
	height float32
}

func NewCubeBoundingBox(position mgl32.Vec3, width, height float32) *CubeBoundingBox {
	maxes := mgl32.Vec3{position.X() + (width / 2), position.Y() + (height / 2), position.Z() + (width / 2)}
	mins := mgl32.Vec3{position.X() - (width / 2), position.Y() - (height / 2), position.Z() - (width / 2)}
	return &CubeBoundingBox{
		Maxes:  maxes,
		Mins:   mins,
		width:  width,
		height: height,
	}
}

func (cbb *CubeBoundingBox) UpdateBounds(position mgl32.Vec3) {
	cbb.Maxes = mgl32.Vec3{position.X() + cbb.width/2, position.Y() + cbb.height/2, position.Z() + cbb.width/2}
	cbb.Mins = mgl32.Vec3{position.X() - cbb.width/2, position.Y() - cbb.height/2, position.Z() - cbb.width/2}
}

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

// OBB-Point collision
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

type SphereCollider struct {
	Center mgl32.Vec3
	Radius float32
}

// Sphere-Point collision
func (sc SphereCollider) CollidesWith(point mgl32.Vec3) bool {
	dist := math2.Distance(mgl32.Vec4{sc.Center[0], sc.Center[1], sc.Center[2], 1.0}, mgl32.Vec4{point[0], point[1], point[2], 1.0})
	return dist <= float64(sc.Radius)
}
