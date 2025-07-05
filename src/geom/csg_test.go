package geom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCuboid(t *testing.T) {
	assertThat := assert.New(t)
	cuboid := NewCuboid(V(0, 0, 0), V(10, 20, 30))

	// front
	assertThat.Equal(cuboid.polygons[0].Vertices[0].Position, V(5, 10, 15))
	assertThat.Equal(cuboid.polygons[0].Vertices[1].Position, V(-5, 10, 15))
	assertThat.Equal(cuboid.polygons[0].Vertices[2].Position, V(-5, -10, 15))
	assertThat.Equal(cuboid.polygons[0].Vertices[3].Position, V(5, -10, 15))

	assertThat.True(cuboid.polygons[0].Plane.Normal.EqualDelta(V(0, 0, 1), TEST_DELTA))

	// back
	assertThat.Equal(cuboid.polygons[1].Vertices[0].Position, V(5, 10, -15))
	assertThat.Equal(cuboid.polygons[1].Vertices[1].Position, V(5, -10, -15))
	assertThat.Equal(cuboid.polygons[1].Vertices[2].Position, V(-5, -10, -15))
	assertThat.Equal(cuboid.polygons[1].Vertices[3].Position, V(-5, 10, -15))

	assertThat.True(cuboid.polygons[1].Plane.Normal.EqualDelta(V(0, 0, -1), TEST_DELTA))
}
