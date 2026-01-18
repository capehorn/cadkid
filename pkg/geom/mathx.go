package geom

import (
	"math"
	"strconv"
)

type OrderedKey = [2]int

func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func Degrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

func LatLngToXYZ(lat, lng float64) Vector {
	lat, lng = Radians(lat), Radians(lng)
	x := math.Cos(lat) * math.Cos(lng)
	y := math.Cos(lat) * math.Sin(lng)
	z := math.Sin(lat)
	return Vector{x, y, z}
}

func Round(a float64) int {
	if a < 0 {
		return int(math.Ceil(a - 0.5))
	} else {
		return int(math.Floor(a + 0.5))
	}
}

func RoundPlaces(a float64, places int) float64 {
	shift := powersOfTen[places]
	return float64(Round(a*shift)) / shift
}

func ParseFloats(items []string) []float64 {
	result := make([]float64, len(items))
	for i, item := range items {
		f, _ := strconv.ParseFloat(item, 64)
		result[i] = f
	}
	return result
}

func SlidingWindow2[T any](s []T, closed bool) [][2]T {
	var ret [][2]T
	for i := 0; i <= len(s)-2; i++ {
		ret = append(ret, [2]T{s[i], s[i+1]})
	}
	if closed {
		ret = append(ret, [2]T{s[len(s)-1], s[0]})
	}
	return ret
}

func SlidingWindow3[T any](s []T, closed bool) [][3]T {
	var ret [][3]T
	l := len(s)
	for i := 0; i <= l-3; i++ {
		ret = append(ret, [3]T{s[i], s[i+1], s[i+2]})
	}
	if closed {
		ret = append(ret, [3]T{s[l-2], s[l-1], s[0]})
		ret = append(ret, [3]T{s[l-1], s[0], s[1]})
	}
	return ret
}

func CreateOrderedKey(p int, q int) OrderedKey {
	if p <= q {
		return [2]int{p, q}
	}
	return [2]int{q, p}
}

func CreateOrderedKeyFromPair(pair [2]int) OrderedKey {
	return CreateOrderedKey(pair[0], pair[1])
}

var powersOfTen = []float64{1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16}

// type Ord int
// const (
// 	LT Ord = iota
// 	EQ
// 	GT
// )

// type Comparator[T any] interface {
// 	ThenCompare(slice []T, fn func (i, j int) Ord) Comparator
// }

// func Compare[T any](slice []T, fn func (i, j int) Ord) Comparator {

// }
