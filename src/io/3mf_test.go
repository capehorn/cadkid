package io

import (
	"testing"
)

func TestMFWriter(t *testing.T) {
	//assertThat := assert.New(t)

	writer := NewMFWriter()
	model := writer.Model()
	model.Metadata("author", "John Doe")
	model.Metadata("application", "cadkid")
	model.Metadata("date", "2025-07-18")
	resources := model.Resources()
	mesh := resources.Mesh()
	mesh.Vertex(0, 0, 0)
	mesh.Vertex(10, 0, 0)
	mesh.Vertex(10, 10, 0)

	mesh.Triangle(0, 1, 2)

	writer.Done()

}
