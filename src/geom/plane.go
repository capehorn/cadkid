package geom

type Plane struct {
	Normal   Vector
	Distance float64
}

func NewPlane(a, b, c Vector) Plane {
	n := b.Sub(a).Cross(c.Sub(a)).Normalize()
	return Plane{n, n.Dot(a)}
}

func (p Plane) Clone() Plane {
	return Plane{Normal: p.Normal.Clone(), Distance: p.Distance}
}

func (p Plane) Flip() Plane {
	return Plane{Normal: p.Normal.Negate(), Distance: -p.Distance}
}
