package io

type MFWriter struct {
	elements []any
}

func NewMFWriter() *MFWriter {
	return &MFWriter{}
}

func (mw *MFWriter) Done() {

}

type MFModel struct {
	writer *MFWriter
}

func (mw *MFWriter) Model() MFModel {
	return MFModel{writer: w}
}

func (m MFModel) Metadata(name, text string) MFMetadata {
	// TODO write name and text
	return MFMetadata{writer: m.writer}
}

type MFMetadata struct {
	writer *MFWriter
}

func (m MFMetadata) PreserveAttr(preserve bool) MFMetadata {
	return m
}

func (m MFMetadata) TypeAttr(t string) MFMetadata {
	return m
}

func (m MFModel) Resources() MFResources {
	return MFResources{writer: m.writer}
}

func (r MFResources) Mesh() MFMesh {
	return MFMesh{writer: r.writer}
}

type MFMesh struct {
	writer *MFWriter
}

func (m MFMesh) Vertex(x, y, z float64) MFMesh {
	// TODO write x, y, z
	return m
}

func (m MFMesh) Triangle(v1, v2, v3 uint32) MFMesh {
	// TODO write x, y, z
	return m
}

type MFResources struct {
	writer *MFWriter
}

type MFReader struct {
}
