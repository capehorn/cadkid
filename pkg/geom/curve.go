package geom

import "math"

type ParametricCurve struct {
	Name            string
	PointAt         func(t float64) Vector
	IsClosed        bool
	IsPlaneCurve    bool
	IsPolygonalPath bool
	Components      []ParametricCurve
}

//func (p ParametricPlaneCurve) Remap(tStart, tEnd float64) {
//	p.PointAt = func(t float64) Vector { return V(0, 0, 0) }
//}

func Circle(radius float64) ParametricCurve {
	p := 2.0 * math.Pi
	return ParametricCurve{
		"circle",
		func(t float64) Vector { return V(radius*math.Cos(p*t), radius*math.Sin(p*t), 0) },
		true,
		true,
		false,
		nil,
	}
}

func NGon(numOfSegments int, radius float64) []Vector {
	result := make([]Vector, numOfSegments)
	result[0] = Vector{radius, 0, 0}
	rot := 2 * math.Pi / float64(numOfSegments)
	var mat Matrix
	for i := 1; i < numOfSegments; i++ {
		mat = Rotate(UnitZ(), rot*float64(i))
		result[i] = mat.MulPosition(result[0])
	}
	return result
}

func Square(sideLength float64) []Vector {
	result := make([]Vector, 4)
	x := sideLength / 2
	result[0] = V(x, x, 0)
	result[1] = V(-x, x, 0)
	result[2] = V(-x, -x, 0)
	result[3] = V(x, -x, 0)
	return result
}

func (c ParametricCurve) TangentAt(t, delta float64) Vector {
	var p1, p2 Vector
	halfDelta := delta / 2
	if c.IsClosed {
		p1 = c.PointAt(t - halfDelta)
		p2 = c.PointAt(t + halfDelta)
	} else {
		p1 = c.PointAt(math.Max(0, t-halfDelta))
		p2 = c.PointAt(math.Min(1, t+halfDelta))
	}
	return p2.Sub(p1).Normalize()
}

func (c ParametricCurve) NormalAt(t, delta float64) Vector {
	var p1 = c.PointAt(math.Max(0, t-delta))
	var p2 = c.PointAt(math.Min(1, t+delta))
	return p2.Sub(p1).Normalize()
}
