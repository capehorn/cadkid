package geom

type Vertex struct {
	Position Vector
	Normal   Vector
}

func InterpolateVertices(v1, v2, v3 Vertex, b VectorW) Vertex {
	v := Vertex{}
	v.Position = InterpolateVectors(v1.Position, v2.Position, v3.Position, b)
	v.Normal = InterpolateVectors(v1.Normal, v2.Normal, v3.Normal, b).Normalize()
	return v
}

func InterpolateFloats(v1, v2, v3 float64, b VectorW) float64 {
	var n float64
	n += v1 * b.X
	n += v2 * b.Y
	n += v3 * b.Z
	return n * b.W
}

func InterpolateVectors(v1, v2, v3 Vector, b VectorW) Vector {
	n := Vector{}
	n = n.Add(v1.MulScalar(b.X))
	n = n.Add(v2.MulScalar(b.Y))
	n = n.Add(v3.MulScalar(b.Z))
	return n.MulScalar(b.W)
}

func InterpolateVectorWs(v1, v2, v3, b VectorW) VectorW {
	n := VectorW{}
	n = n.Add(v1.MulScalar(b.X))
	n = n.Add(v2.MulScalar(b.Y))
	n = n.Add(v3.MulScalar(b.Z))
	return n.MulScalar(b.W)
}

// Interpolate creates a new vertex between this vertex and `other` by linearly
// interpolating all properties using a parameter of `t`. Subclasses should
// override this to interpolate additional properties.
func (v Vertex) Interpolate(other Vertex, t float64) Vertex {
	return Vertex{
		Position: v.Position.Lerp(other.Position, t),
		Normal:   v.Normal.Lerp(other.Normal, t),
	}
}

func Barycentric(p1, p2, p3, p Vector) VectorW {
	v0 := p2.Sub(p1)
	v1 := p3.Sub(p1)
	v2 := p.Sub(p1)
	d00 := v0.Dot(v0)
	d01 := v0.Dot(v1)
	d11 := v1.Dot(v1)
	d20 := v2.Dot(v0)
	d21 := v2.Dot(v1)
	d := d00*d11 - d01*d01
	v := (d11*d20 - d01*d21) / d
	w := (d00*d21 - d01*d20) / d
	u := 1 - v - w
	return VectorW{u, v, w, 1}
}

func (v Vertex) Clone() Vertex {
	return Vertex{v.Position.Clone(), v.Normal.Clone()}
}

func (v Vertex) Flip() Vertex {
	return Vertex{Position: v.Position, Normal: v.Normal.Negate()}
}
