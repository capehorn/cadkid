package cad

import (
	. "capehorn/cadkid/geom"
)

type DatumId int32

type Datum struct {
	id    DatumId
	Label string
	Frame Frame
}
