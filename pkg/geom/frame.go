package geom

import "log"

var StandardFrame = Frame{
	ZeroVector(),
	UnitX(),
	UnitY(),
	UnitZ(),
}

type Frame struct {
	Orig, E1, E2, E3 Vector
}

func RotationMinimizingFrames(points []Vector, tangents []Vector, initialFrame Frame) []Frame {
	if !initialFrame.E3.Equal(tangents[0]) {
		log.Print("For calculating RMF, the initial frame's third component must be equal to the first tangent vector")
		return nil
	}
	frames := make([]Frame, len(points))
	frames[0] = initialFrame

	for i := 0; i < len(tangents)-1; i++ {
		tg := tangents[i]
		tgNext := tangents[i+1]
		r := frames[i].E1
		v1 := points[i+1].Sub(points[i])
		c1 := v1.Dot(v1)
		rli := r.Sub(v1.MulScalar(v1.Dot(r)).MulScalar(2 / c1))
		tli := tg.Sub(v1.MulScalar(v1.Dot(tg)).MulScalar(2 / c1))
		v2 := tgNext.Sub(tli)
		c2 := v2.Dot(v2)
		rNext := rli.Sub(v2.MulScalar(v2.Dot(rli)).MulScalar(2 / c2))
		sNext := tgNext.Cross(rNext)
		frames[i+1] = Frame{points[i+1], rNext, sNext, tgNext}
	}
	return frames
}
