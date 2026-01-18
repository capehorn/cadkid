package geom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPolygon_Transform(t *testing.T) {
	assertThat := assert.New(t)
	p := NewPolygon([]Vertex{
		VertexOf(V(10, 10, 0)),
		VertexOf(V(10, -10, 0)),
		VertexOf(V(-10, -10, 0)),
		VertexOf(V(-10, 10, 0)),
	}, nil)

	assertThat.True(p.Plane.Normal.EqualDelta(V(0, 0, -1), TEST_DELTA))
}
