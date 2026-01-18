package sdf

import (
	g "capehorn/cadkid/pkg/geom"
	"math"
)

type FRep interface {
	Eval(v g.Vector) float64
	//bbox() BBox
}

type Operator interface {
	Apply(freps []FRep, p g.Vector) float64
}

type SdfSphere struct {
	Center g.Vector
	Radius float64
}

func (f SdfSphere) Eval(p g.Vector) float64 {
	return p.Sub(f.Center).Length() - f.Radius
}

type SdfBox struct {
	Center   g.Vector
	HalfSide g.Vector
}

func (f SdfBox) Eval(p g.Vector) float64 {
	v := p.Sub(f.Center)
	return v.Abs().Sub(f.HalfSide).MaxScalar(0.0).Length()
	//
	//
	// v := p.Sub(f.Center)
	// q := v.Abs().Sub(f.HalfSide)
	// return q.Max(0.0).Length() + math.Min(math.Max(q.X, math.Max(q.Y, q.Z)), 0.0)
}

type FRepDerived struct {
	Items    []FRep
	Operator Operator
}

func (f FRepDerived) Eval(p g.Vector) float64 {
	return f.Operator.Apply(f.Items, p)
}

type OpUnion struct{}

func (op OpUnion) Apply(freps []FRep, p g.Vector) float64 {
	minDist := freps[0].Eval(p)
	for i := 1; i < len(freps); i++ {
		d := freps[i].Eval(p)
		if d < minDist {
			minDist = d
		}
	}
	return minDist
}

type OpOnion struct {
	Dist float64
}

func (op OpOnion) Apply(freps []FRep, p g.Vector) float64 {
	dist := freps[0].Eval(p)
	return math.Abs(dist) - op.Dist
}
