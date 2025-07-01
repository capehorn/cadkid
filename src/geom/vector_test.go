package geom

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var Tolerance = 0.00000001

func TestVector_SegmentDistance(t *testing.T) {
	p := V(40, 10, 0)

	segmentStart := V(0, 0, 0)
	segmentEnd := V(100, 0, 0)

	distance := p.SegmentDistance(segmentStart, segmentEnd)
	if math.Abs(distance-10) > Tolerance {
		t.Errorf("Output %f is not equal with expected %f", distance, 10.0)
	}
}

func TestVector_EqualDelta(t *testing.T) {
	assertThat := assert.New(t)
	a := V(40, 10, 0)
	b := V(40.1, 9.9, 0)

	assertThat.True(a.EqualDelta(b, 0.11))
}

//func StringFormat(t *testing.T) {
//	//p := V(40, 10, 0)
//	//print(p.String())
//}
