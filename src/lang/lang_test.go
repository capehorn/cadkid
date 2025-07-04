package lang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var Tolerance = 0.00000001

func TestFn_MapOverSlice(t *testing.T) {
	assertThat := assert.New(t)
	intValues := []int{1, 2, 3}
	intResult := Map(intValues, func(v int) int { return 100 + v })
	assertThat.Equal(intResult, []int{101, 102, 103})

	dummyValues := []dummy{dummy{1}, dummy{2}, dummy{3}, dummy{4}}
	dummyResult := Map(dummyValues, dummy.dummyAdder)
	assertThat.Equal(dummyResult[0], dummy{11})
	assertThat.Equal(dummyResult[1], dummy{12})
	assertThat.Equal(dummyResult[2], dummy{13})
	assertThat.Equal(dummyResult[3], dummy{14})
}

type dummy struct {
	v int
}

func (d dummy) dummyAdder() dummy {
	return dummy{v: d.v + 10}
}
