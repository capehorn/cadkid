package geom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PolyMesh_twoQuads(t *testing.T) {
	assert := assert.New(t)

	m := newTwoQuadsMesh()

	assert.Equal(6, len(m.vertices))
	assert.Equal([]int{0, 1, 2, 3, 0, 1, 4, 5}, m.faces)
	assert.Equal([]int{0, 4}, m.faceOffsets)

	edges := getEdges(m.edges)
	assert.Equal(7, len(edges))
	assert.Contains(edges, [4]int{0, 1, 0, 1}) // common
	assert.Contains(edges, [4]int{0, 0, 1, 2}) // first face
	assert.Contains(edges, [4]int{0, 0, 2, 3}) // first face
	assert.Contains(edges, [4]int{0, 0, 0, 3}) // first face
	assert.Contains(edges, [4]int{1, 1, 1, 4}) // second face
	assert.Contains(edges, [4]int{1, 1, 4, 5}) // second face
	assert.Contains(edges, [4]int{1, 1, 0, 5}) // second face

	assert.Equal([]int{0, 1, 2, 3}, m.GetFace(0))
	assert.Equal([]int{0, 1, 4, 5}, m.GetFace(1))
	assert.Equal([]int(nil), m.GetFace(2))
}

func Test_PolyMesh_Cube(t *testing.T) {
	assert := assert.New(t)
	m := newCubeMesh()

	assert.Equal(8, len(m.vertices))

	assert.Equal([]int{0, 1, 2, 3}, m.GetFace(0))
	assert.Equal([]int{0, 1, 4, 5}, m.GetFace(1))
	assert.Equal([]int{2, 3, 6, 7}, m.GetFace(2))
	assert.Equal([]int{6, 7, 4, 5}, m.GetFace(3))
	assert.Equal([]int{1, 2, 7, 4}, m.GetFace(4))
	assert.Equal([]int{0, 3, 6, 5}, m.GetFace(5))
	assert.Equal([]int(nil), m.GetFace(7))

	assert.Equal(6, len(m.GetFaces()))

	assert.Equal([]int{0, 4, 8, 12, 16, 20}, m.faceOffsets)

	edges := getEdges(m.edges)
	assert.Equal(12, len(edges))
	assert.Contains(edges, [4]int{0, 1, 0, 1})
	assert.Contains(edges, [4]int{0, 2, 2, 3})
	assert.Contains(edges, [4]int{0, 4, 1, 2})
	assert.Contains(edges, [4]int{0, 5, 0, 3})

	assert.Equal(V(5, 5, 0), m.faceAveragePoint([]int{0, 1, 2, 3}))

	// faceAvgs := []Vector{}
	// for _, face := range m.GetFaces() {
	// 	faceAvgs = append(faceAvgs, m.faceAveragePoint(face))
	// }

	// edgeMidAndNewVertices := make(map[OrderedKey]edgeMidAndNewVertex)

	// m.computeMidEdgeAndNewEdgeVertex(0, 1, faceAvgs, m, edgeMidAndNewVertices)
	// m.computeMidEdgeAndNewEdgeVertex(1, 2, faceAvgs, m, edgeMidAndNewVertices)

	// assert.Equal(edgeMidAndNewVertex{edgeMidPoint: V(5, 0, 0), newVertex: V(5, 1.25, 1.25), vIdx: 8}, edgeMidAndNewVertices[CreateOrderedKey(0, 1)])
	// assert.Equal(edgeMidAndNewVertex{edgeMidPoint: V(10, 5, 0), newVertex: V(8.75, 5, 1.25), vIdx: 9}, edgeMidAndNewVertices[CreateOrderedKey(1, 2)])
}

func Test_PolyMesh_Cube_deleteFace(t *testing.T) {
	assert := assert.New(t)
	// GIVEN - a cube
	m := newCubeMesh()

	assert.Equal(8, len(m.vertices))

	assert.Equal([]int{0, 1, 2, 3}, m.GetFace(0))
	assert.Equal([]int{0, 1, 4, 5}, m.GetFace(1))
	assert.Equal([]int{2, 3, 6, 7}, m.GetFace(2))
	assert.Equal([]int{6, 7, 4, 5}, m.GetFace(3))
	assert.Equal([]int{1, 2, 7, 4}, m.GetFace(4))
	assert.Equal([]int{0, 3, 6, 5}, m.GetFace(5))
	assert.Equal([]int(nil), m.GetFace(7))

	assert.Equal(6, len(m.GetFaces()))

	// WHEN - delete face 1
	m.DeleteFace(1)

	// THEN
	assert.Equal([]int{0, 1, 2, 3}, m.GetFace(0))
	assert.Equal([]int{2, 3, 6, 7}, m.GetFace(1))
	assert.Equal([]int{6, 7, 4, 5}, m.GetFace(2))
	assert.Equal([]int{1, 2, 7, 4}, m.GetFace(3))
	assert.Equal([]int{0, 3, 6, 5}, m.GetFace(4))
	assert.Equal([]int(nil), m.GetFace(6))

	edges := getEdges(m.edges)
	assert.Equal(12, len(edges))
}

func Test_PolyMesh_quads_deleteFace(t *testing.T) {
	assert := assert.New(t)
	// GIVEN - connected quads
	m := newTwoQuadsMesh()
	// WHEN - delete first face
	m.DeleteFace(0)
	// THEN
	assert.Equal([]int{0, 1, 4, 5}, m.GetFace(0))
	edges := getEdges(m.edges)
	assert.Equal(4, len(edges))
	assert.Contains(edges, [4]int{0, 0, 0, 1})
	// TODO complete edge assertions
}

func Test_PolyMesh_emptyMesh_addFace(t *testing.T) {
	assert := assert.New(t)
	// GIVEN - an empty mesh
	m := NewEmptyPolyMesh()
	// WHEN - add a face
	vIndices := m.AddVertices(V(0, 0, 0), V(10, 0, 0), V(10, 10, 0), V(0, 10, 0))
	m.AddFace(vIndices)
	// THEN
	assert.Equal([]int{0, 1, 2, 3}, m.GetFace(0))
}

func newTwoQuadsMesh() *PolyMesh {
	vertices := []Vector{V(0, 0, 0), V(10, 0, 0), V(10, 10, 0), V(0, 10, 0), V(10, 0, 10), V(0, 10, 10)}
	faces := [][]int{{0, 1, 2, 3}, {0, 1, 4, 5}}
	return NewPolyMesh(vertices, faces)
}

func newCubeMesh() *PolyMesh {
	vertices := []Vector{V(0, 0, 0), V(10, 0, 0), V(10, 10, 0), V(0, 10, 0), V(10, 0, 10), V(0, 0, 10), V(0, 10, 10), V(10, 10, 10)}
	faces := [][]int{{0, 1, 2, 3}, {0, 1, 4, 5}, {2, 3, 6, 7}, {6, 7, 4, 5}, {1, 2, 7, 4}, {0, 3, 6, 5}}
	return NewPolyMesh(vertices, faces)
}

func getEdges(flatEdges []int) [][4]int {
	edges := [][4]int{}
	for i := 0; i < len(flatEdges)/4; i++ {
		edges = append(edges, [4]int{flatEdges[4*i], flatEdges[4*i+1], flatEdges[4*i+2], flatEdges[4*i+3]})
	}
	return edges
}
