package io

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMFWriter(t *testing.T) {
	assertThat := assert.New(t)
	result := bytes.NewBuffer([]byte{})

	writer := NewMFWriter(result)
	model := writer.Model(MFModelAttr{Unit: Millimeter, Lang: "en"})
	model.Metadata("John Doe", "name", "author")
	model.Metadata("cadkid", "name", "application")
	model.Metadata("2025-07-18", "name", "date")
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
