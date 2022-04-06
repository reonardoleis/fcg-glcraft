package collisions

type CubeCollider struct{}

func NewCubeCollider() *CubeCollider {
	return &CubeCollider{}
}

func (cc CubeCollider) Collides(a, b CubeBoundingBox) bool {
	return (a.Maxes.X() >= b.Mins.X() && a.Mins.X() <= b.Maxes.X()) &&
		(a.Maxes.Y() >= b.Mins.Y() && a.Mins.Y() <= b.Maxes.Y()) &&
		(a.Maxes.Z() >= b.Mins.Z() && a.Mins.Z() <= b.Maxes.Z())
}
