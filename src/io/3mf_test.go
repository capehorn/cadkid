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
	model := writer.Model()
	model.Metadata("John Doe", "name", "author")
	model.Metadata("cadkid", "name", "application")
	model.Metadata("2025-07-18", "name", "date")
	resources := model.Resources()
	mesh := resources.Mesh()
	mesh.Vertex(0, 0, 0)
	mesh.Vertex(10, 0, 0)
	mesh.Vertex(10, 10, 0)

	mesh.Triangle(0, 1, 2)

	writer.Done()

	assertThat.True(result.Len() > 0)
}
