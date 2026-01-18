package cad

import (
	g "capehorn/cadkid/pkg/geom"
)

type PartId = int32

type Part struct {
	Id     PartId
	Label  string
	Datums map[DatumId]g.Matrix
}
