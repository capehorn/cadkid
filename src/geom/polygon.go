package geom

import "capehorn/cadkid/lang"

type Polygon struct {
	Vertices []Vertex
	Shared   any
	Plane    Plane
}

func NewPolygon(vertices []Vertex, shared any) Polygon {
	if len(vertices) < 3 {
		panic("polygon: at least 3 vertices needed")
	}
	return Polygon{
		Vertices: vertices,
		Shared:   shared,
		Plane:    NewPlane(vertices[0].Position, vertices[1].Position, vertices[2].Position)}
}

func (p Polygon) Clone() Polygon {
	vs := make([]Vertex, len(p.Vertices))
	for i := 0; i < len(p.Vertices); i++ {
		vs[i] = p.Vertices[i].Clone()
	}
	return Polygon{
		Vertices: vs,
		Shared:   p.Shared,
		Plane:    NewPlane(vs[0].Position, vs[1].Position, vs[2].Position)}
}

func (p Polygon) Flip() Polygon {
	vs := make([]Vertex, len(p.Vertices))
	for i := len(p.Vertices) - 1; 0 < i; i-- {
		vs[i] = p.Vertices[i].Flip()
	}
	return Polygon{
		Vertices: vs,
		Shared:   p.Shared,
		Plane:    NewPlane(vs[0].Position, vs[1].Position, vs[2].Position),
	}
}

func (p Polygon) Transform(m Matrix) Polygon {
	vertices := lang.Map(p.Vertices, func(v Vertex) Vertex {
		return Vertex{
			Position: m.MulPosition(v.Position),
			Normal:   m.MulDirection(v.Normal)}
	})
	return Polygon{Vertices: vertices}
}
