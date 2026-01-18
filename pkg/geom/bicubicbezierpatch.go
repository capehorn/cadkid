package geom

// https://www.educative.io/answers/what-is-a-bzier-patch
// p(u,v) = u × B × P × transpose(B) × v
// u = [ u * u * u, u * u, u, 1 ]
// v = [ v * v * v, v * v, v, 1 ]
//
// B = [
//	     -1  3 -3  1
//        3 -6  3  0
//       -3  3  0  0
//        1  0  0  0
//     ]
// The control points are calculated using the Bernstein polynomial
// P = [
//	     p(1,1) p(1,2) p(1,3) p(1,4)
//       p(2,1) p(2,2)...
//       ...
//       ...           p(4,3) p(4,4)
//     ]

type BicubicBezierPatch struct {
	Points  []Vector
	tPoints []Vector
}

func NewBicubicBezierPatch(points [16]Vector) BicubicBezierPatch {
	bbp := BicubicBezierPatch{Points: points[:]}
	bbp.tPoints = bbp.computeSurfaceControlPoints()
	return bbp
}

func (bbp BicubicBezierPatch) computeSurfaceControlPoints() []Vector {
	ps := bbp.Points
	cb0 := NewCubicBezier(ps[0], ps[4], ps[8], ps[12])
	cb1 := NewCubicBezier(ps[1], ps[5], ps[9], ps[13])
	cb2 := NewCubicBezier(ps[2], ps[6], ps[10], ps[14])
	cb3 := NewCubicBezier(ps[3], ps[7], ps[11], ps[15])

	m := []Vector{
		cb0.tcps[0],
		cb0.tcps[1],
		cb0.tcps[2],
		cb0.tcps[3],

		cb1.tcps[0],
		cb1.tcps[1],
		cb1.tcps[2],
		cb1.tcps[3],

		cb2.tcps[0],
		cb2.tcps[1],
		cb2.tcps[2],
		cb2.tcps[3],

		cb3.tcps[0],
		cb3.tcps[1],
		cb3.tcps[2],
		cb3.tcps[3],
	}

	cb_0 := NewCubicBezier(m[0], m[4], m[8], m[12])
	cb_1 := NewCubicBezier(m[1], m[5], m[9], m[13])
	cb_2 := NewCubicBezier(m[2], m[6], m[10], m[14])
	cb_3 := NewCubicBezier(m[3], m[7], m[11], m[15])

	return []Vector{

		cb_0.tcps[0],
		cb_0.tcps[1],
		cb_0.tcps[2],
		cb_0.tcps[3],

		cb_1.tcps[0],
		cb_1.tcps[1],
		cb_1.tcps[2],
		cb_1.tcps[3],

		cb_2.tcps[0],
		cb_2.tcps[1],
		cb_2.tcps[2],
		cb_2.tcps[3],

		cb_3.tcps[0],
		cb_3.tcps[1],
		cb_3.tcps[2],
		cb_3.tcps[3],
	}
}

func (bbp BicubicBezierPatch) PointAt(u, v float64) Vector {
	uExp2 := u * u
	uExp3 := uExp2 * u
	vExp2 := v * v
	vExp3 := vExp2 * v

	tcps := bbp.tPoints

	tps := []Vector{
		tcps[0].Add(tcps[4].MulScalar(u)).Add(tcps[8].MulScalar(uExp2)).Add(tcps[12].MulScalar(uExp3)),
		tcps[1].Add(tcps[5].MulScalar(u)).Add(tcps[9].MulScalar(uExp2)).Add(tcps[13].MulScalar(uExp3)),
		tcps[2].Add(tcps[6].MulScalar(u)).Add(tcps[10].MulScalar(uExp2)).Add(tcps[14].MulScalar(uExp3)),
		tcps[3].Add(tcps[7].MulScalar(u)).Add(tcps[11].MulScalar(uExp2)).Add(tcps[15].MulScalar(uExp3)),
	}

	return tps[0].Add(tps[1].MulScalar(v)).Add(tps[2].MulScalar(vExp2)).Add(tps[3].MulScalar(vExp3))
}

func (bbp BicubicBezierPatch) Triangulate(m, n int) []*Triangle {
	triangles := make([]*Triangle, 0, m*n*2)
	du := 1 / float64(m)
	dv := 1 / float64(n)
	u := 0.0
	v := 0.0

	for i := 0; i < m; i++ {
		uNext := u + du
		if i == m-1 {
			uNext = 1.0
		}
		v = 0
		for j := 0; j < n; j++ {
			vNext := v + dv
			if j == n-1 {
				vNext = 1.0
			}
			tBottom, tTop := QuadToTriangles(bbp.PointAt(u, v), bbp.PointAt(u, vNext), bbp.PointAt(uNext, vNext), bbp.PointAt(uNext, v))
			triangles = append(triangles, tBottom, tTop)
			v = vNext
		}
		u = uNext
	}
	return triangles
}

func (bbp BicubicBezierPatch) Transform(matrix Matrix) BicubicBezierPatch {
	for i, p := range bbp.Points {
		bbp.Points[i] = matrix.MulPosition(p)
	}
	bbp.tPoints = bbp.computeSurfaceControlPoints()
	return bbp
}

func (bbp BicubicBezierPatch) SplitAt(u, v float64) (BicubicBezierPatch, BicubicBezierPatch, BicubicBezierPatch, BicubicBezierPatch) {
	ps := bbp.Points
	cb0_0, cb0_1 := NewCubicBezier(ps[0], ps[4], ps[8], ps[12]).SplitAt(u)
	cb1_0, cb1_1 := NewCubicBezier(ps[1], ps[5], ps[9], ps[13]).SplitAt(u)
	cb2_0, cb2_1 := NewCubicBezier(ps[2], ps[6], ps[10], ps[14]).SplitAt(u)
	cb3_0, cb3_1 := NewCubicBezier(ps[3], ps[7], ps[11], ps[15]).SplitAt(u)

	m_0 := []Vector{
		cb0_0.P0,
		cb0_0.P1,
		cb0_0.P2,
		cb0_0.P3,

		cb1_0.P0,
		cb1_0.P1,
		cb1_0.P2,
		cb1_0.P3,

		cb2_0.P0,
		cb2_0.P1,
		cb2_0.P2,
		cb2_0.P3,

		cb3_0.P0,
		cb3_0.P1,
		cb3_0.P2,
		cb3_0.P3,
	}

	cb0_0_0, cb0_0_1 := NewCubicBezier(m_0[0], m_0[4], m_0[8], m_0[12]).SplitAt(v)
	cb1_0_0, cb1_0_1 := NewCubicBezier(m_0[1], m_0[5], m_0[9], m_0[13]).SplitAt(v)
	cb2_0_0, cb2_0_1 := NewCubicBezier(m_0[2], m_0[6], m_0[10], m_0[14]).SplitAt(v)
	cb3_0_0, cb3_0_1 := NewCubicBezier(m_0[3], m_0[7], m_0[11], m_0[15]).SplitAt(v)

	m_1 := []Vector{
		cb0_1.P0,
		cb0_1.P1,
		cb0_1.P2,
		cb0_1.P3,

		cb1_1.P0,
		cb1_1.P1,
		cb1_1.P2,
		cb1_1.P3,

		cb2_1.P0,
		cb2_1.P1,
		cb2_1.P2,
		cb2_1.P3,

		cb3_1.P0,
		cb3_1.P1,
		cb3_1.P2,
		cb3_1.P3,
	}

	cb0_1_0, cb0_1_1 := NewCubicBezier(m_1[0], m_1[4], m_1[8], m_1[12]).SplitAt(v)
	cb1_1_0, cb1_1_1 := NewCubicBezier(m_1[1], m_1[5], m_1[9], m_1[13]).SplitAt(v)
	cb2_1_0, cb2_1_1 := NewCubicBezier(m_1[2], m_1[6], m_1[10], m_1[14]).SplitAt(v)
	cb3_1_0, cb3_1_1 := NewCubicBezier(m_1[3], m_1[7], m_1[11], m_1[15]).SplitAt(v)

	bbp0 := NewBicubicBezierPatch([16]Vector{
		cb0_0_0.P0, cb0_0_0.P1, cb0_0_0.P2, cb0_0_0.P3,
		cb1_0_0.P0, cb1_0_0.P1, cb1_0_0.P2, cb1_0_0.P3,
		cb2_0_0.P0, cb2_0_0.P1, cb2_0_0.P2, cb2_0_0.P3,
		cb3_0_0.P0, cb3_0_0.P1, cb3_0_0.P2, cb3_0_0.P3,
	})

	bbp1 := NewBicubicBezierPatch([16]Vector{
		cb0_0_1.P0, cb0_0_1.P1, cb0_0_1.P2, cb0_0_1.P3,
		cb1_0_1.P0, cb1_0_1.P1, cb1_0_1.P2, cb1_0_1.P3,
		cb2_0_1.P0, cb2_0_1.P1, cb2_0_1.P2, cb2_0_1.P3,
		cb3_0_1.P0, cb3_0_1.P1, cb3_0_1.P2, cb3_0_1.P3,
	})

	bbp2 := NewBicubicBezierPatch([16]Vector{
		cb0_1_0.P0, cb0_1_0.P1, cb0_1_0.P2, cb0_1_0.P3,
		cb1_1_0.P0, cb1_1_0.P1, cb1_1_0.P2, cb1_1_0.P3,
		cb2_1_0.P0, cb2_1_0.P1, cb2_1_0.P2, cb2_1_0.P3,
		cb3_1_0.P0, cb3_1_0.P1, cb3_1_0.P2, cb3_1_0.P3,
	})

	bbp3 := NewBicubicBezierPatch([16]Vector{
		cb0_1_1.P0, cb0_1_1.P1, cb0_1_1.P2, cb0_1_1.P3,
		cb1_1_1.P0, cb1_1_1.P1, cb1_1_1.P2, cb1_1_1.P3,
		cb2_1_1.P0, cb2_1_1.P1, cb2_1_1.P2, cb2_1_1.P3,
		cb3_1_1.P0, cb3_1_1.P1, cb3_1_1.P2, cb3_1_1.P3,
	})

	return bbp0, bbp1, bbp2, bbp3
}
