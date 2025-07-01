package geom

var EPSILON_PLANE = 0.001

type CSG struct {
	polygons []Polygon
}

func NewRectCuboid(x, y, z float64) CSG {
	return CSG{polygons: nil}
}

func NewSphere(r float64) CSG {
	return CSG{polygons: nil}
}

func NewCylinder(r, y float64) CSG {
	return CSG{polygons: nil}
}

func (csg CSG) clone(that CSG) CSG {
	return CSG{}
}

func (csg CSG) union(that CSG) CSG {
	return CSG{}
}

func (csg CSG) subtract(that CSG) CSG {
	return CSG{}
}

func (csg CSG) intersect(that CSG) CSG {
	return CSG{}
}

func (csg CSG) inverse(that CSG) CSG {
	return CSG{}
}

type node struct {
	plane    Plane
	front    *node
	back     *node
	polygons []Polygon
}

func (n *node) clone() *node {
	ps := make([]Polygon, len(n.polygons))
	for i := 0; 0 < len(n.polygons); i++ {
		ps[i] = n.polygons[i].Clone()
	}
	node := &node{}
	node.plane = n.plane.Clone()
	node.front = n.front.clone()
	node.back = n.back.clone()
	node.polygons = ps
	return node
}

// invert converts solid space to empty space and empty space to solid space.
func (n *node) invert() {
	for i, p := range n.polygons {
		n.polygons[i] = p.Flip()
	}
	n.plane = n.plane.Flip()
	n.front.invert()
	n.back.invert()
	temp := n.front
	n.front = n.back
	n.back = temp
}

// Recursively remove all polygons in `polygons` that are inside this BSP
// tree.
func (n *node) clipPolygons(polygons []Polygon) []Polygon {
	//if (!this.plane) return polygons.slice();
	var coplanarFront, coplanarBack, front, back []Polygon
	for _, p := range polygons {
		splitPolygon(p, n.plane, coplanarFront, coplanarBack, front, back)
	}
	front = n.front.clipPolygons(front)
	back = n.back.clipPolygons(back)
	return append(front, back...)
}

// clipTo removes all polygons in this BSP tree that are inside the `other` BSP tree
func (n *node) clipTo(other *node) {
	n.polygons = other.clipPolygons(n.polygons)
	n.front.clipTo(other)
	n.back.clipTo(other)
}

// Return a list of all polygons in this BSP tree.
func (n *node) allPolygons() []Polygon {
	allPolygons := append(make([]Polygon, 0, len(n.polygons)), n.polygons...)
	allPolygons = append(allPolygons, n.front.allPolygons()...)
	allPolygons = append(allPolygons, n.back.allPolygons()...)
	return allPolygons
}

// build builds a BSP tree out of `polygons`. When called on an existing tree, the
// new polygons are filtered down to the bottom of the tree and become new
// nodes there. Each set of polygons is partitioned using the first polygon
// (no heuristic is used to pick a good split).
func (n *node) build(polygons []Polygon) {
	if n.plane.Normal.Equal(ZeroVector()) {
		n.plane = polygons[0].Plane.Clone()
	}

	front := make([]Polygon, 0)
	back := make([]Polygon, 0)
	for _, p := range polygons {
		splitPolygon(p, n.plane, polygons, polygons, front, back)
	}
	if (front.length) {
	if (!this.front) this.front = new CSG.Node();
	this.front.build(front);
	}
	if (back.length) {
	if (!this.back) this.back = new CSG.Node();
	this.back.build(back);
	}
}

type RelationType int8

const (
	COPLANAR RelationType = iota
	FRONT
	BACK
	SPANNING
)

// Split `polygon` by this plane if needed, then put the polygon or polygon
// fragments in the appropriate slices. Coplanar polygons go into either
// `coplanarFront` or `coplanarBack` depending on their orientation with
// respect to this plane. Polygons in front or in back of this plane go into
// either `front` or `back`.
func splitPolygon(polygon Polygon, plane Plane, coplanarFront, coplanarBack, front, back []Polygon) {

	// Classify each point as well as the entire polygon into one of the above four classes
	var polyType = COPLANAR
	types := make([]RelationType, len(polygon.Vertices))
	for i := 0; i < len(polygon.Vertices); i++ {
		t := plane.Normal.Dot(polygon.Vertices[i].Position) - plane.Distance
		var vertexType = COPLANAR
		if t < -EPSILON_PLANE {
			vertexType = BACK
		} else if t > EPSILON_PLANE {
			vertexType = FRONT
		}

		polyType |= vertexType
		types[i] = vertexType
	}

	// Put the polygon in the correct list, splitting it when necessary.
	switch polyType {
	case COPLANAR:
		if 0 < plane.Normal.Dot(polygon.Plane.Normal) {
			coplanarFront = append(coplanarFront, polygon)
		} else {
			coplanarBack = append(coplanarBack, polygon)
		}
	case FRONT:
		front = append(front, polygon)
	case BACK:
		back = append(back, polygon)
	case SPANNING:
		f := make([]Vertex, 0)
		b := make([]Vertex, 0)
		numOfVertices := len(polygon.Vertices)
		for i, vi := range polygon.Vertices {
			j := (i + 1) % numOfVertices
			ti := types[i]
			tj := types[j]
			vj := polygon.Vertices[j]
			if ti != BACK {
				f = append(f, vi)
			}
			if ti != FRONT {
				if ti != BACK {
					b = append(b, vi.Clone())
				} else {
					b = append(b, vi)
				}
			}
			if (ti | tj) == SPANNING {
				t := (plane.Distance - plane.Normal.Dot(vi.Position)) / plane.Normal.Dot(vj.Position.Sub(vi.Position))
				v := vi.Interpolate(vj, t)
				f = append(f, v)
				b = append(b, v.Clone())
			}
		}
		if len(f) >= 3 {
			front = append(front, NewPolygon(f, polygon.Shared))
		}
		if len(b) >= 3 {
			back = append(back, NewPolygon(b, polygon.Shared))
		}

	}
}
