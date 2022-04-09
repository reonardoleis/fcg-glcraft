package collisions

import "github.com/go-gl/mathgl/mgl32"

type CubeCollider struct{}

func NewCubeCollider() *CubeCollider {
	return &CubeCollider{}
}

func (cc CubeCollider) Collides(a, b CubeBoundingBox, av, bv [8]mgl32.Vec3) bool {
	//fmt.Println(a, b)
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
