package geom

//import (
//	"log"
//	"testing"
//
//	//"github.com/stretchr/testify/assert"
//)

//func Test_SubD_aQuad(t *testing.T) {
//	assert := assert.New(t)
//	vertices := []Vector{V(0, 0, 0), V(10, 0, 0), V(10, 10, 0), V(0, 10, 0)}
//	faces := [][]int{{0, 1, 2, 3}}
//	m := NewPolyMesh(vertices, faces)
//
//	s := m.SubD()
//
//	assert.Equal(9, len(s.vertices))
//	assert.Equal(4*4, len(s.faces))
//	subDFaces := s.GetFaces()
//	assert.Contains(subDFaces, []int{1, 5, 4, 6})
//}
//
//func Test_SubD_twoQuads(t *testing.T) {
//	assert := assert.New(t)
//	vertices := []Vector{V(0, 0, 0), V(10, 0, 0), V(10, 10, 0), V(0, 10, 0), V(10, 0, 10), V(0, 0, 10)}
//	faces := [][]int{{0, 1, 2, 3}, {0, 1, 4, 5}}
//	m := NewPolyMesh(vertices, faces)
//
//	s := m.SubD()
//
//	assert.Equal(15, len(s.vertices))
//	assert.Equal(8*4, len(s.faces))
//}
//
//func Test_SubD_Triangulate_cube(t *testing.T) {
//	vertices := []Vector{V(0, 0, 0), V(10, 0, 0), V(10, 10, 0), V(0, 10, 0), V(10, 0, 10), V(0, 0, 10), V(0, 10, 10), V(10, 10, 10)}
//	faces := [][]int{{0, 1, 2, 3}, {0, 1, 4, 5}, {2, 3, 6, 7}, {6, 7, 4, 5}, {1, 2, 7, 4}, {0, 3, 6, 5}}
//	s := NewPolyMesh(vertices, faces)
//	s = s.SubD()
//	s = s.SubD()
//	s = s.SubD()
//
//	triangles := s.Triangulate()
//
//	mesh := NewTriangleMesh(triangles)
//
//	// fit mesh in a bi-unit cube centered at the origin
//	mesh.BiUnitCube()
//
//	// smooth the normals
//	mesh.SmoothNormalsThreshold(Radians(30))
//
//	err := SaveSTL("../test/subd-cube.stl", mesh)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
//func Test_SubD_Triangulate_twoQuads(t *testing.T) {
//	vertices := []Vector{V(0, 0, 0), V(10, 0, 0), V(10, 10, 0), V(0, 10, 0), V(10, 0, 10), V(0, 0, 10)}
//	faces := [][]int{{0, 1, 2, 3}, {0, 1, 4, 5}}
//	s := NewPolyMesh(vertices, faces)
//	s = s.SubD()
//	s = s.SubD()
//	s = s.SubD()
//
//	triangles := s.Triangulate()
//
//	mesh := NewTriangleMesh(triangles)
//
//	// fit mesh in a bi-unit cube centered at the origin
//	mesh.BiUnitCube()
//
//	// smooth the normals
//	mesh.SmoothNormalsThreshold(Radians(30))
//
//	err := SaveSTL("../test/subd-two-quads.stl", mesh)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//}
