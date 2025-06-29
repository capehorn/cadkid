package cad

type Link struct {
	Head      Part
	HeadDatum Datum
	Tail      Part
	TailDatum Datum
}

type Construct struct {
	Parts []Part
	Links []Link
}
