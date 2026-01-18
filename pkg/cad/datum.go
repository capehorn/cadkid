package cad

import (
	g "capehorn/cadkid/pkg/geom"
)

type DatumId int32

type Datum struct {
	id    DatumId
	Label string
	Frame g.Frame
}
