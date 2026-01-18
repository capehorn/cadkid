package geom

type edgeMidAndNewVertex struct {
	edgeMidPoint Vector
	newVertex    Vector
	vIdx         int // vertex index in the new subd mesh
}

func (m *PolyMesh) SubD() *PolyMesh {
	// create subd mesh, original verticies are also verticies in this new mesh
	subD := &PolyMesh{vertices: make([]Vector, len(m.vertices)), faces: []int{}, faceOffsets: []int{}, edges: []int{}}
	copy(subD.vertices, m.vertices)

	// tracking where the face avg vertices starts in subd
	faceAvgOffset := len(subD.vertices)
	// for every original face stores the avg point
	faceAvgs := []Vector{}
	faces := m.GetFaces()
	// compute each face's average point
	for _, face := range faces {
		faceAvgs = append(faceAvgs, m.faceAveragePoint(face))
	}
	// add face avg points to the subd vertices
	subD.vertices = append(subD.vertices, faceAvgs...)

	// for every original edge stores its midpoint and the new edge point
	edgeMidAndNewVertices := make(map[OrderedKey]edgeMidAndNewVertex)

	// compute each edge midpoint and the new edge point
	for _, face := range faces {
		for _, pair := range SlidingWindow2(face, true) {
			m.computeMidEdgeAndNewEdgeVertex(pair[0], pair[1], faceAvgs, subD, edgeMidAndNewVertices)
		}
	}

	// track computing the new vertex positions
	computedVertices := make(map[int]Vector)

	// build new faces and compute new vertices position
	for fOffsetIdx := 0; fOffsetIdx < len(m.faceOffsets); fOffsetIdx++ {
		for _, t := range SlidingWindow3(m.GetFace(fOffsetIdx), true) {
			edgeKey01 := CreateOrderedKey(t[0], t[1])
			edgeKey12 := CreateOrderedKey(t[1], t[2])
			// add face offset and face
			subD.faceOffsets = append(subD.faceOffsets, len(subD.faces))
			subD.faces = append(subD.faces, t[1], edgeMidAndNewVertices[edgeKey01].vIdx, faceAvgOffset+fOffsetIdx, edgeMidAndNewVertices[edgeKey12].vIdx)

			// only compute new vertices position if not computed yet
			if _, ok := computedVertices[t[1]]; !ok {
				computedVertices[t[1]] = m.computeVertex(t[1], edgeKey01, edgeKey12, fOffsetIdx, edgeMidAndNewVertices, faceAvgs)
			}
		}
	}

	for k := range computedVertices {
		subD.vertices[k] = computedVertices[k]
	}

	return subD
}

func (m *PolyMesh) computeVertex(vIdx int, edgeKey01, edgeKey02 OrderedKey, fOffsetIdx int, edgeMidAndNewVertices map[OrderedKey]edgeMidAndNewVertex, faceAvgs []Vector) Vector {
	faceOffsetIndices := make(map[int]bool)
	faceOffsetIndices[fOffsetIdx] = true
	edgeKeys := make(map[OrderedKey]bool)
	edgeKeys[edgeKey01] = true
	edgeKeys[edgeKey02] = true

	// find all other edges from the vertex
	for fIdx, face := range m.GetFaces() {
	FaceVertices:
		for _, pair := range SlidingWindow2(face, true) {
			if pair[0] == vIdx || pair[1] == vIdx {
				ek := CreateOrderedKey(pair[0], pair[1])
				if _, ok := edgeKeys[ek]; !ok {
					edgeKeys[ek] = true
					if _, ok := faceOffsetIndices[fIdx]; !ok {
						faceOffsetIndices[fIdx] = true
					}
					break FaceVertices
				}
			}
		}
	}

	newVertex := V(0, 0, 0)

	for k := range faceOffsetIndices {
		newVertex = newVertex.Add(faceAvgs[k])
	}

	newVertex = newVertex.DivScalar(float64(len(faceOffsetIndices)))

	midEdgeSum := V(0, 0, 0)

	for ek := range edgeKeys {
		midEdgeSum = midEdgeSum.Add(edgeMidAndNewVertices[ek].edgeMidPoint)
	}

	midEdgeSum = midEdgeSum.DivScalar(float64(len(edgeKeys)))
	newVertex = newVertex.Add(midEdgeSum.MulScalar(2))
	newVertex = newVertex.Add(m.vertices[vIdx].MulScalar(float64(len(edgeKeys) - 3)))
	return newVertex.DivScalar(float64(len(edgeKeys)))
}

func (m *PolyMesh) findCommonEdgeIdx(vIdxA int, vIdxB int) int {
	for i := 0; i < len(m.edges)/4; i++ {
		if (m.edges[4*i+2] == vIdxA && m.edges[4*i+3] == vIdxB) || (m.edges[4*i+3] == vIdxA && m.edges[4*i+2] == vIdxB) {
			return i
		}
	}
	return -1
}

func (m *PolyMesh) computeMidEdgeAndNewEdgeVertex(vIdxA int, vIdxB int, faceAvgs []Vector, subD *PolyMesh, edgeMidAndNewPoints map[OrderedKey]edgeMidAndNewVertex) {
	edgeKey := CreateOrderedKey(vIdxA, vIdxB)

	// only add if absent
	if _, ok := edgeMidAndNewPoints[edgeKey]; !ok {
		edgeMidPoint := m.vertices[vIdxA].Lerp(m.vertices[vIdxB], 0.5)
		var newEdgePoint Vector
		commonEdgeIdx := m.findCommonEdgeIdx(vIdxA, vIdxB)
		if commonEdgeIdx == -1 {
			newEdgePoint = edgeMidPoint
		} else {
			f := faceAvgs[m.edges[commonEdgeIdx*4]]
			g := faceAvgs[m.edges[commonEdgeIdx*4+1]]
			faceMidPoint := f.Lerp(g, 0.5)
			newEdgePoint = edgeMidPoint.Lerp(faceMidPoint, 0.5)
		}
		subD.vertices = append(subD.vertices, newEdgePoint)
		edgeMidAndNewPoints[edgeKey] = edgeMidAndNewVertex{edgeMidPoint: edgeMidPoint, newVertex: newEdgePoint, vIdx: len(subD.vertices) - 1}
	}
}
