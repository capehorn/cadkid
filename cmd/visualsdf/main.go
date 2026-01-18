package main

import (
	g "capehorn/cadkid/pkg/geom"
	sdf "capehorn/cadkid/pkg/sdf"
	"fmt"
	"image/color"
	"math"

	"github.com/fogleman/gg"
)

var width = 600
var height = 600
var isoStep = 10
var isoThickness = 2

var colorRed = color.NRGBA{255, 0, 0, 255}
var colorGreen = color.NRGBA{0, 255, 0, 255}
var colorBlue = color.NRGBA{0, 0, 255, 255}
var colorHalfRed = color.NRGBA{127, 0, 0, 255}
var colorHalfGreen = color.NRGBA{0, 127, 0, 255}
var colorHalfBlue = color.NRGBA{0, 0, 127, 255}

// UNION
// var derived = sdf.FRepDerived{
// 	Items: []sdf.FRep{
// 		sdf.Sphere{Center: sdf.V(200, 200, 0), Radius: 80},
// 		sdf.Sphere{Center: sdf.V(260, 200, 0), Radius: 80},
// 		sdf.Sphere{Center: sdf.V(280, 240, 0), Radius: 60},
// 	},
// 	Operator: "union",
// }
// var freps = []sdf.FRep{derived}

// BOX
//var freps = []sdf.FRep{sdf.Box{sdf.V(400, 400, 0), sdf.V(40, 60, 0)}}

// ONION
// var derived = sdf.FRepDerived{
// 	Items: []sdf.FRep{
// 		sdf.Sphere{Center: sdf.V(200, 200, 0), Radius: 80},
// 	},
// 	Operator: "onion",
// 	OpProps: map[string]any{
// 		"r": 10.0,
// 	},
// }
// var freps = []sdf.FRep{derived}

// MIX
var derived = sdf.FRepDerived{
	Items: []sdf.FRep{
		sdf.SdfSphere{Center: g.V(200, 200, 0), Radius: 80},
		sdf.SdfSphere{Center: g.V(350, 200, 0), Radius: 80},
		sdf.SdfSphere{Center: g.V(280, 240, 0), Radius: 60},
	},
	Operator: sdf.OpUnion{},
}
var onion = sdf.FRepDerived{
	Items: []sdf.FRep{
		derived,
	},
	Operator: sdf.OpOnion{Dist: 20.0},
}
var freps = []sdf.FRep{onion}

func main() {
	dc := gg.NewContext(600, 600)

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			dists := make([]float64, len(freps))
			if i == 310 && j == 220 {
				fmt.Print("hello")
			}
			for n, frep := range freps {
				dists[n] = frep.Eval(g.V(float64(i), float64(j), 0))
			}
			color := getColor(dists)
			dc.SetColor(color)
			dc.SetPixel(i, j)
		}
	}
	dc.SavePNG("out.png")
}

// 	//drawCircle(dc, color.NRGBA{128, 0, 0, 255}, 200, 200, 40)
// 	//drawCircle(dc, color.NRGBA{0, 128, 128, 255}, 300, 200, 40)

// 	dc.SavePNG("out.png")
// }

//	func drawCircle(dc *gg.Context, color color.Color, cx float64, cy float64, r float64) {
//		dc.DrawCircle(cx, cy, r)
//		dc.SetColor(color)
//		dc.Fill()
//	}
func getColor(dists []float64) color.Color {
	color := colorHalfGreen

	dist := dists[0]
	q := float64(dist) / float64(isoStep)
	if math.Abs(dist) < float64(isoThickness) {
		return colorHalfRed
	} else if math.Abs(float64(dist)-math.Ceil(q)*float64(isoStep)) < float64(isoThickness) {
		return colorHalfBlue
	}
	return color
}

// 	// for _, dist := range dists {
// 	// 	q := float64(dist) / float64(isoStep)
// 	// 	if math.Abs(dist) < float64(isoThickness) {
// 	// 		return colorHalfRed
// 	// 	} else if math.Abs(float64(dist)-math.Ceil(q)*float64(isoStep)) < float64(isoThickness) {
// 	// 		return colorHalfBlue
// 	// 	}
// 	// }
// 	return color
// }
