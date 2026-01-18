package geom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNGon(t *testing.T) {
	assertThat := assert.New(t)
	points := NGon(4, 10)

	assertThat.True(points[0].EqualDelta(V(10, 0, 0), TEST_DELTA))
	assertThat.True(points[1].EqualDelta(V(0, -10, 0), TEST_DELTA))
	assertThat.True(points[2].EqualDelta(V(-10, 0, 0), TEST_DELTA))
	assertThat.True(points[3].EqualDelta(V(0, 10, 0), TEST_DELTA))
}

func TestCircle(t *testing.T) {
	assertThat := assert.New(t)
	curve := Circle(10)

	assertThat.True(curve.PointAt(0).EqualDelta(V(10, 0, 0), TEST_DELTA))
	assertThat.True(curve.PointAt(0.25).EqualDelta(V(0, 10, 0), TEST_DELTA))
	assertThat.True(curve.PointAt(0.5).EqualDelta(V(-10, 0, 0), TEST_DELTA))
	assertThat.True(curve.PointAt(0.75).EqualDelta(V(0, -10, 0), TEST_DELTA))
	assertThat.True(curve.PointAt(1).EqualDelta(V(10, 0, 0), TEST_DELTA))
}
