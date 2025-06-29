package cad

import (
	. "capehorn/cadkid/geom"
)

type PartId = int32

type Part struct {
	Id     PartId
	Label  string
	Datums map[DatumId]Matrix
}
