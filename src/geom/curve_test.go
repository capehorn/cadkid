package geom

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var DELTA = 0.1

func TestNGon(t *testing.T) {
	assertThat := assert.New(t)
	points := NGon(4, 10)

	assertThat.True(points[0].EqualDelta(V(10, 0, 0), DELTA))
	assertThat.True(points[1].EqualDelta(V(0, -10, 0), DELTA))
	assertThat.True(points[2].EqualDelta(V(-10, 0, 0), DELTA))
	assertThat.True(points[3].EqualDelta(V(0, 10, 0), DELTA))
}

func TestCircle(t *testing.T) {
	assertThat := assert.New(t)
	curve := Circle(10)

	assertThat.True(curve.PointAt(0).EqualDelta(V(10, 0, 0), DELTA))
	assertThat.True(curve.PointAt(0.25).EqualDelta(V(0, 10, 0), DELTA))
	assertThat.True(curve.PointAt(0.5).EqualDelta(V(-10, 0, 0), DELTA))
	assertThat.True(curve.PointAt(0.75).EqualDelta(V(0, -10, 0), DELTA))
	assertThat.True(curve.PointAt(1).EqualDelta(V(10, 0, 0), DELTA))
}
