package geom

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlidingWindow2_open(t *testing.T) {
	assert := assert.New(t)
	s := []int{0, 1, 2}

	assert.Equal([][2]int{{0, 1}, {1, 2}}, SlidingWindow2(s, false))
}

func TestSlidingWindow2_closed(t *testing.T) {
	assert := assert.New(t)
	s := []int{0, 1, 2}

	assert.Equal([][2]int{{0, 1}, {1, 2}, {2, 0}}, SlidingWindow2(s, true))
}

func TestSlidingWindow3_open(t *testing.T) {
	assert := assert.New(t)
	s := []int{0, 1, 2, 3}

	assert.Equal([][3]int{{0, 1, 2}, {1, 2, 3}}, SlidingWindow3(s, false))
}

func TestSlidingWindow3_closed(t *testing.T) {
	assert := assert.New(t)
	s := []int{0, 1, 2, 3}

	assert.Equal([][3]int{{0, 1, 2}, {1, 2, 3}, {2, 3, 0}, {3, 0, 1}}, SlidingWindow3(s, true))
}

func TestMakeSlice(t *testing.T) {
	assert := assert.New(t)
	s := make([]int, 3)
	s[0] = 10

	assert.Equal([]int{10, 0, 0}, s)
}

func TestCopySliceOfVectors(t *testing.T) {
	assert := assert.New(t)
	s := []Vector{V(0, 0, 0), V(1, 1, 1)}
	s2 := make([]Vector, len(s))

	copy(s2, s)

	assert.Equal([]Vector{V(0, 0, 0), V(1, 1, 1)}, s2)

	s[0] = V(10, 10, 10)
	s[1].X = 100
	assert.Equal([]Vector{V(0, 0, 0), V(1, 1, 1)}, s2)
}

func Test_deleteFromSlice(t *testing.T) {
	assert := assert.New(t)

	ints := []int{1, 2, 3, 4, 5, 6}
	afterDelete := slices.Delete(ints, 2, 4)
	assert.Equal([]int{1, 2, 5, 6}, afterDelete)
	assert.Equal(4, len(afterDelete))
	assert.Equal(6, cap(afterDelete))
}

func Test_deleteInPlaceFromSlice(t *testing.T) {
	assert := assert.New(t)

	ints := []int{1, 2, 3, 4, 5, 6}

	j := 0
	for _, v := range ints {
		if v < 4 {
			ints[j] = v
			j++
		}
	}

	ints = ints[:j]
	assert.Equal([]int{1, 2, 3}, ints)
}
