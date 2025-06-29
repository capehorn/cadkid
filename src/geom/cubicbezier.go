package geom

type CubicBezier struct {
	P0, P1, P2, P3 Vector
	tcps           [4]Vector
}

func NewCubicBezier(p0, p1, p2, p3 Vector) CubicBezier {
	cb := CubicBezier{P0: p0, P1: p1, P2: p2, P3: p3}
	cb.tcps = cb.precomputeTransformedControlPoints()
	return cb
}

func (cb CubicBezier) PointAt(t float64) Vector {
	tExp2 := t * t
	tExp3 := tExp2 * t
	return cb.tcps[0].Add(cb.tcps[1].MulScalar(t)).Add(cb.tcps[2].MulScalar(tExp2)).Add(cb.tcps[3].MulScalar(tExp3))
}

func (cb CubicBezier) precomputeTransformedControlPoints() [4]Vector {
	return [4]Vector{
		cb.P0,
		cb.P0.MulScalar(-3).Add(cb.P1.MulScalar(3)),
		cb.P0.MulScalar(3).Add(cb.P1.MulScalar(-6)).Add(cb.P2.MulScalar(3)),
		cb.P0.MulScalar(-1).Add(cb.P1.MulScalar(3)).Add(cb.P2.MulScalar(-3)).Add(cb.P3),
	}
}

func (cb CubicBezier) Transform(matrix Matrix) CubicBezier {
	cb.P0 = matrix.MulPosition(cb.P0)
	cb.P1 = matrix.MulPosition(cb.P1)
	cb.P2 = matrix.MulPosition(cb.P2)
	cb.P3 = matrix.MulPosition(cb.P3)
	cb.tcps = cb.precomputeTransformedControlPoints()
	return cb
}

func (cb CubicBezier) SplitAt(t float64) (CubicBezier, CubicBezier) {
	p := cb.PointAt(t)
	m := cb.P1.Lerp(cb.P2, t)

	q1 := cb.P0.Lerp(cb.P1, t)
	q2 := q1.Lerp(m, t)

	cb0 := NewCubicBezier(cb.P0, q1, q2, p)

	r2 := cb.P3.Lerp(cb.P2, 1-t)
	r1 := r2.Lerp(m, 1-t)
	cb1 := NewCubicBezier(p, r1, r2, cb.P3)
	return cb0, cb1
}

func (cb CubicBezier) Approximate(tolerance float64, points []Vector) []Vector {
	sDist1 := cb.P1.SegmentDistance(cb.P0, cb.P3)
	sDist2 := cb.P2.SegmentDistance(cb.P0, cb.P3)
	if sDist1 <= tolerance && sDist2 <= tolerance {
		points = append(points, cb.P0, cb.P3)
	} else {
		cb0, cb1 := cb.SplitAt(0.5)
		points = cb0.Approximate(tolerance, points)
		points = cb1.Approximate(tolerance, points)
	}
	return points
}
