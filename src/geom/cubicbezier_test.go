package geom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApproximate(t *testing.T) {
	assert := assert.New(t)
	cb := NewCubicBezier(V(0, 0, 0), V(300, 100, 0), V(250, 380, 0), V(350, 220, 0))

	points := cb.Approximate(5, make([]Vector, 0))
	assert.Len(points, 16)
}

func TestSplitAt(t *testing.T) {
	assert := assert.New(t)
	cb := NewCubicBezier(V(0, 0, 0), V(300, 100, 0), V(250, 380, 0), V(350, 220, 0))

	left, right := cb.SplitAt(0.5)
	assert.Len(left.tcps, 4)
	assert.Len(right.tcps, 4)
}
