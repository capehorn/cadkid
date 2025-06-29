package geom

type PolyMeshBuilder struct {
	*PolyMesh
	sel *Selection
}

type Selection struct {
	vertices    []int
	faceOffsets []int
	edges       [][2]int
}

func NewEmptyPolyMeshBuilder() *PolyMeshBuilder {
	return &PolyMeshBuilder{PolyMesh: NewEmptyPolyMesh(), sel: &Selection{}}
}

func NewPolyMeshBuilder(pm *PolyMesh) *PolyMeshBuilder {
	return &PolyMeshBuilder{PolyMesh: pm, sel: &Selection{vertices: []int{}, faceOffsets: []int{}, edges: [][2]int{}}}
}

func (b *PolyMeshBuilder) SelectFace(fOffsetIdx int) {
	b.sel.selectFace(fOffsetIdx)
}

func (b *PolyMeshBuilder) ClearSelection() {
	b.sel.clear()
}

// func (b *PolyMeshBuilder) Pull(d float64) {
// 	if 0 < len(b.sel.faceOffsets) {
// 		for _, fOffsetIdx := range b.sel.faceOffsets {
// 			face := b.GetFace(fOffsetIdx)
// 			byVector := b.ComputeFaceNormal(face).MulScalar(d)
// 			for _, vIdx := range face {
// 				newVertex := Translate(byVector).MulPosition(b.vertices[vIdx])

// 			}
// 		}
// 	}
// }

func (b *PolyMeshBuilder) FaceOffset(d float64) {

}

func (b *PolyMeshBuilder) Rotate(deg float64) {

}

// SELECTION

func (sel *Selection) selectFace(fOffsetIdx int) {
	for _, v := range sel.faceOffsets {
		if v == fOffsetIdx {
			return
		}
	}
	sel.faceOffsets = append(sel.faceOffsets, fOffsetIdx)
}

func (sel *Selection) clear() {
	sel.vertices = []int{}
	sel.faceOffsets = []int{}
	sel.edges = [][2]int{}
}
