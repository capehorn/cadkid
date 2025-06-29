package geom

import "slices"

type PolyMesh struct {
	vertices    []Vector
	faces       []int
	faceOffsets []int
	edges       []int // length is multiple of 4, repetition of: face offset idx, face offset idx, vertex idx, vertex idx
}

func NewEmptyPolyMesh() *PolyMesh {
	return &PolyMesh{vertices: []Vector{}, faces: []int{}, faceOffsets: []int{}, edges: []int{}}
}

func NewPolyMesh(vertices []Vector, faces [][]int) *PolyMesh {
	m := NewEmptyPolyMesh()
	m.vertices = append([]Vector(nil), vertices...)

	// key: ordered edge, value: faceIndices
	edgeHits := make(map[OrderedKey][]int)

	for fIdx, face := range faces {
		m.faceOffsets = append(m.faceOffsets, len(m.faces))
		m.faces = append(m.faces, face...)

		for _, pair := range SlidingWindow2(face, true) {
			edgeKey := CreateOrderedKey(pair[0], pair[1])
			if _, ok := edgeHits[edgeKey]; !ok {
				edgeHits[edgeKey] = []int{fIdx}
			} else {
				edgeHits[edgeKey] = append(edgeHits[edgeKey], fIdx)
			}
		}
	}

	for k, v := range edgeHits {
		switch len(v) {
		case 1:
			m.edges = append(m.edges, v[0], v[0], k[0], k[1]) // egde is bounding one face
		case 2:
			m.edges = append(m.edges, v[0], v[1], k[0], k[1]) // edge is bounding two faces
		default:
			panic("Edge bounds more than 2 faces")
		}
	}
	return m
}

func (m *PolyMesh) GetFace(fOffsetIdx int) []int {
	if 0 <= fOffsetIdx {
		if fOffsetIdx == len(m.faceOffsets)-1 {
			return m.faces[m.faceOffsets[fOffsetIdx]:]
		} else if fOffsetIdx < len(m.faceOffsets)-1 {
			return m.faces[m.faceOffsets[fOffsetIdx]:m.faceOffsets[fOffsetIdx+1]]
		}
	}
	return nil
}

func (m *PolyMesh) GetFaces() [][]int {
	faces := [][]int{}
	for fOffsetIdx := 0; fOffsetIdx < len(m.faceOffsets); fOffsetIdx++ {
		faces = append(faces, m.GetFace(fOffsetIdx))
	}
	return faces
}

func (m *PolyMesh) GetCommonEdgesIndices(fOffsetIdx int, vIdx int) []int {
	commonEdgesIndices := []int{}
	for i := 0; i < len(m.edges)-4; i += 4 {
		if (m.edges[i] == fOffsetIdx || m.edges[i+1] == fOffsetIdx) && (m.edges[i+2] == vIdx || m.edges[i+3] == vIdx) {
			commonEdgesIndices = append(commonEdgesIndices, i/4)
		}
	}
	return commonEdgesIndices
}

func (m *PolyMesh) Triangulate() []*Triangle {
	triangles := []*Triangle{}
	for _, face := range m.GetFaces() {
		switch len(face) {

		case 3:
			triangles = append(triangles, NewTriangleForPoints(m.vertices[face[0]], m.vertices[face[1]], m.vertices[face[2]]))
		case 4:
			t0, t1 := QuadToTriangles(m.vertices[face[0]], m.vertices[face[1]], m.vertices[face[2]], m.vertices[face[3]])
			triangles = append(triangles, t0, t1)
		default:
			avgPoint := m.faceAveragePoint(face)
			for _, pair := range SlidingWindow2(face, true) {
				triangles = append(triangles, NewTriangleForPoints(avgPoint, m.vertices[pair[0]], m.vertices[pair[1]]))
			}
		}
	}
	return triangles
}

func (m *PolyMesh) ComputeFaceNormal(face []int) Vector {
	switch len(face) {
	case 3:
		return NewTriangleForPoints(m.vertices[face[0]], m.vertices[face[1]], m.vertices[face[2]]).Normal()
	default:
		avgPoint := m.faceAveragePoint(face)
		v := Vector{}
		for _, pair := range SlidingWindow2(face, true) {
			v = v.Add(NewTriangleForPoints(avgPoint, m.vertices[pair[0]], m.vertices[pair[1]]).Normal())
		}
		return v.Normalize()
	}
}

func (m *PolyMesh) AddVertex(vertex Vector) int {
	m.vertices = append(m.vertices, vertex)
	return len(m.vertices)
}

func (m *PolyMesh) AddVertices(vertices ...Vector) []int {
	firstIdx := len(m.vertices)
	m.vertices = append(m.vertices, vertices...)
	retVal := make([]int, len(m.vertices)-firstIdx)
	j := 0
	for i := firstIdx; i < len(m.vertices); i++ {
		retVal[j] = i
		j++
	}
	return retVal
}

func (m *PolyMesh) AddFace(face []int) {
	sharedEdges := make(map[OrderedKey]bool)
	for i := 0; i < len(m.edges)/4; i++ {
		for _, pair := range SlidingWindow2(face, true) {
			if (m.edges[4*i+2] == pair[0] && m.edges[4*i+3] == pair[1]) || (m.edges[4*i+2] == pair[1] && m.edges[4*i+3] == pair[0]) {
				f0 := m.edges[4*i+0]
				f1 := m.edges[4*i+1]
				if f0 == f1 {
					m.edges[4*i+1] = len(m.faceOffsets) - 1
					sharedEdges[CreateOrderedKeyFromPair(pair)] = true
				}
			}
		}
	}

	for _, pair := range SlidingWindow2(face, true) {
		if _, ok := sharedEdges[CreateOrderedKeyFromPair(pair)]; !ok {
			m.edges = append(m.edges, len(m.faceOffsets), len(m.faceOffsets), pair[0], pair[1])
		}
	}

	m.faceOffsets = append(m.faceOffsets, len(m.faces))
	m.faces = append(m.faces, face...)
}

func (m *PolyMesh) DeleteFace(fOffsetIdx int) {
	face := m.GetFace(fOffsetIdx)
	faceStartIdx := m.faceOffsets[fOffsetIdx]

	m.faces = slices.Delete(m.faces, faceStartIdx, faceStartIdx+len(face))

	//deleteEdge := false

	// copy edges slice in place while removing edges that bounding only the face that we delete
	// and modify other accordingly (shift them)
	j := 0 // output index
	for i, _ := range m.edges {
		if i%4 == 0 {
			if m.edges[i] == fOffsetIdx && m.edges[i+1] == fOffsetIdx { // bounding edge of deleted face (not shared)
				continue
			} else {
				if fOffsetIdx == m.edges[i] {
					m.edges[i] = m.edges[i+1]
				}

				if fOffsetIdx == m.edges[i+1] {
					m.edges[i+1] = m.edges[i]
				}

				if fOffsetIdx < m.edges[i] {
					m.edges[j] = m.edges[i] - 1
				} else {
					m.edges[j] = m.edges[i]
				}

				if fOffsetIdx < m.edges[i+1] {
					m.edges[j+1] = m.edges[i+1] - 1
				} else {
					m.edges[j+1] = m.edges[i+1]
				}

				m.edges[j+2] = m.edges[i+2]
				m.edges[j+3] = m.edges[i+3]
				j += 4
			}
		}
	}
	m.edges = m.edges[:j]
}

func (m *PolyMesh) faceAveragePoint(face []int) Vector {
	x := 0.0
	y := 0.0
	z := 0.0

	l := len(face)
	for i := 0; i < l; i++ {
		vertex := m.vertices[face[i]]
		x += vertex.X
		y += vertex.Y
		z += vertex.Z
	}
	return V(x/float64(l), y/float64(l), z/float64(l))
}
