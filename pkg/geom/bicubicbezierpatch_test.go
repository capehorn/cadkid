package geom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBicubicBezier_SplitAt(t *testing.T) {
	assert := assert.New(t)
	bb := NewBicubicBezierPatch([16]Vector{
		V(0, 0, 0), V(100, 0, 100), V(200, 0, 100), V(300, 0, 100),
		V(0, 100, 0), V(100, 100, 100), V(200, 100, 100), V(300, 100, 100),
		V(0, 200, 0), V(100, 200, 100), V(200, 200, 100), V(300, 200, 100),
		V(0, 300, 0), V(100, 300, 100), V(200, 300, 100), V(300, 300, 100),
	})

	bb0, bb1, bb2, bb3 := bb.SplitAt(0.5, 0.5)
	assert.Len(bb0.tPoints, 16)
	assert.Len(bb1.tPoints, 16)
	assert.Len(bb2.tPoints, 16)
	assert.Len(bb3.tPoints, 16)
}
