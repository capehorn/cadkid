package io

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMFWriter_withIndent4(t *testing.T) {
	assertThat := assert.New(t)
	result := bytes.NewBuffer([]byte{})

	writer := NewMFWriter(result, 4)
	model := writer.Model(MFModelAttr{Unit: Millimeter, Lang: "en-us"})
	model.Metadata("John Doe", MFMetadataAttr{Name: "author"})
	model.Metadata("cadkid", MFMetadataAttr{Name: "application", Preserve: True})
	model.Metadata("2025-07-18", MFMetadataAttr{Name: "date", Type: "xs:date"})
	resources := model.Resources()
	baseMat := resources.BaseMaterials(1)
	baseMat.Base("green", "#21BB4CFF")
	baseMat.Base("red", "#FF0000FF")

	obj1 := resources.Object(MFObjectAttr{Id: 1, Name: "anObject"})
	mesh := obj1.Mesh()
	vs := mesh.Vertices()
	vs.Vertex(0, 0, 0)
	vs.Vertex(10, 0, 0)
	vs.Vertex(10, 10, 0)

	ts := mesh.Triangles()
	ts.Triangle(0, 1, 2, nil)

	obj2 := resources.Object(MFObjectAttr{Id: 2, Name: "otherObject"})
	components := obj2.Components()
	components.Component(3, nil)
	components.Component(4, nil)
	writer.Done()

	assertThat.True(result.Len() > 0)
}
