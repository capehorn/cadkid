package io

import (
	"capehorn/cadkid/lang"
	"fmt"
	"io"
	"strings"
)

type MFWriter struct {
	elements lang.Stack
	writer   io.Writer
}

func NewMFWriter(writer io.Writer) *MFWriter {
	return &MFWriter{writer: writer}
}

func (mw *MFWriter) Done() {

}

func (mw *MFWriter) write(s string) {
	_, err := mw.writer.Write([]byte(s))
	if err != nil {
		contextErr := fmt.Errorf("failed to write data in 3mf format %w", err)
		fmt.Println(contextErr)
		return
	}
}

type MFModel struct {
	writer *MFWriter
}

func (mw *MFWriter) Model() MFModel {
	model := MFModel{writer: mw}
	mw.elements.Push(model)
	mw.write("<model>")
	return model
}

func (m MFModel) Metadata(text string, kvAttributes ...string) MFMetadata {
	if len(kvAttributes) == 0 {
		m.writer.write("<metadata>")
	} else {
		m.writer.write("<metadata")
		m.writer.write(kvAttributesToString(kvAttributes...))
		m.writer.write(">")
	}
	m.writer.write(text)
	m.writer.write("</metadata>")
	return MFMetadata{writer: m.writer}
}

type MFMetadata struct {
	writer *MFWriter
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

func kvAttributesToString(kvAttributes ...string) string {
	if len(kvAttributes) == 0 {
		return ""
	}
	if len(kvAttributes)%2 != 0 {
		err := fmt.Errorf("odd number of kbAttributes")
		fmt.Println(err)
		return ""
	}
	attrs := make([]string, len(kvAttributes)/2)
	for i := 0; i < len(kvAttributes); i += 2 {
		attrs = append(attrs, "\""+kvAttributes[i]+"\"=\""+kvAttributes[i+1]+"\"")
	}
	return strings.Join(attrs, " ")
}
