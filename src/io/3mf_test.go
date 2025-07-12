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
	model := writer.Model(MFModelAttr{Unit: Millimeter, Lang: "en"})
	model.Metadata("John Doe", MFMetadataAttr{Name: "author"})
	model.Metadata("cadkid", MFMetadataAttr{Name: "application", Preserve: True})
	model.Metadata("2025-07-18", MFMetadataAttr{Name: "date", Type: "xs:date"})
	resources := model.Resources()
	baseMat := resources.BaseMaterials(1)
	baseMat.Base("green", "#21BB4CFF")
	baseMat.Base("red", "#FF0000FF")
	obj1 := resources.Object(MFObjectAttr{Id: 1})
	mesh := obj1.Mesh()
	mesh.Vertex(0, 0, 0)
	mesh.Vertex(10, 0, 0)
	mesh.Vertex(10, 10, 0)

	mesh.Triangle(0, 1, 2)

	writer.Done()

	assertThat.True(result.Len() > 0)
}
