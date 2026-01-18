package geom

import "capehorn/cadkid/pkg/lang"

var EPSILON_PLANE = 0.001

type CSG struct {
	polygons []Polygon
}

func NewCuboid(c, s Vector) CSG {
	xHalf := s.X / 2
	yHalf := s.Y / 2
	zHalf := s.Z / 2
	vs := []Vertex{
		// front face
		VertexOf(V(c.X+xHalf, yHalf, zHalf)),
		VertexOf(V(c.X-xHalf, yHalf, zHalf)),
		VertexOf(V(c.X-xHalf, -yHalf, zHalf)),
		VertexOf(V(c.X+xHalf, -yHalf, zHalf)),

		// back face
		VertexOf(V(c.X+xHalf, yHalf, -zHalf)),
		VertexOf(V(c.X+xHalf, -yHalf, -zHalf)),
		VertexOf(V(c.X-xHalf, -yHalf, -zHalf)),
		VertexOf(V(c.X-xHalf, +yHalf, -zHalf)),
	}

	return NewCsg([]Polygon{
		NewPolygon([]Vertex{vs[0], vs[1], vs[2], vs[3]}, nil), // front
		NewPolygon([]Vertex{vs[4], vs[5], vs[6], vs[7]}, nil), // back
		NewPolygon([]Vertex{vs[0], vs[4], vs[7], vs[1]}, nil), // top
		NewPolygon([]Vertex{vs[3], vs[2], vs[6], vs[5]}, nil), // bottom
		NewPolygon([]Vertex{vs[1], vs[7], vs[6], vs[2]}, nil), // left
		NewPolygon([]Vertex{vs[0], vs[3], vs[5], vs[4]}, nil), // right
	})
}

func NewSphere(r float64) CSG {
	return CSG{polygons: nil}
}

func NewCylinder(r, y float64) CSG {
	return CSG{polygons: nil}
}

func NewCsg(polygons []Polygon) CSG {
	return CSG{polygons: polygons}
}

func (csg CSG) clone() CSG {
	clone := CSG{}
	clone.polygons = lang.Map(csg.polygons, Polygon.Clone)
	return clone
}

func (csg CSG) Union(that CSG) CSG {
	var a = newNode(lang.Map(csg.polygons, Polygon.Clone))
	var b = newNode(lang.Map(that.polygons, Polygon.Clone))
	a.clipTo(b)
	b.clipTo(a)
	b.invert()
	b.clipTo(a)
	b.invert()
	a.build(b.allPolygons())
	return NewCsg(a.allPolygons())
}

func (csg CSG) Subtract(that CSG) CSG {
	var a = newNode(lang.Map(csg.polygons, Polygon.Clone))
	var b = newNode(lang.Map(that.polygons, Polygon.Clone))
	a.invert()
	a.clipTo(b)
	b.clipTo(a)
	b.invert()
	b.clipTo(a)
	b.invert()
	a.build(b.allPolygons())
	a.invert()
	return NewCsg(a.allPolygons())
}

func (csg CSG) Intersect(that CSG) CSG {
	var a = newNode(lang.Map(csg.polygons, Polygon.Clone))
	var b = newNode(lang.Map(that.polygons, Polygon.Clone))
	a.invert()
	b.clipTo(a)
	b.invert()
	a.clipTo(b)
	b.clipTo(a)
	a.build(b.allPolygons())
	a.invert()
	return NewCsg(a.allPolygons())
}

func (csg CSG) Inverse() CSG {
	var clone = csg.clone()
	clone.polygons = lang.Map(clone.polygons, Polygon.Flip)
	return clone
}

type node struct {
	plane    Plane
	front    *node
	back     *node
	polygons []Polygon
}

func newNode(polygons []Polygon) *node {
	node := &node{}
	if len(polygons) == 0 {
		return node
	}
	node.build(polygons)
	return node
}

func (n *node) clone() *node {
	node := &node{}
	node.plane = n.plane.Clone()
	if n.front != nil {
		node.front = n.front.clone()
	}
	if n.back != nil {
		node.back = n.back.clone()
	}
	node.polygons = lang.Map(n.polygons, Polygon.Clone)
	return node
}

// invert converts solid space to empty space and empty space to solid space.
func (n *node) invert() {
	for i, p := range n.polygons {
		n.polygons[i] = p.Flip()
	}
	n.plane = n.plane.Flip()
	if n.front != nil {
		n.front.invert()
	}
	if n.back != nil {
		n.back.invert()
	}
	temp := n.front
	n.front = n.back
	n.back = temp
}

// Recursively remove all polygons in `polygons` that are inside this BSP
// tree.
func (n *node) clipPolygons(polygons []Polygon) []Polygon {

	if n.plane.Normal.Equal(ZeroVector()) {
		return lang.Copy(polygons)
	}
	var coplanarFront, coplanarBack, front, back []Polygon
	for _, p := range polygons {
		splitPolygon(p, n.plane, coplanarFront, coplanarBack, front, back)
	}
	if n.front != nil {
		front = n.front.clipPolygons(front)
	}
	if n.back != nil {
		back = n.back.clipPolygons(back)
	}
	return append(front, back...)
}

// clipTo removes all polygons in this BSP tree that are inside the `other` BSP tree
func (n *node) clipTo(other *node) {
	n.polygons = other.clipPolygons(n.polygons)
	if n.front != nil {
		n.front.clipTo(other)
	}
	if n.back != nil {
		n.back.clipTo(other)
	}
}

// Return a list of all polygons in this BSP tree.
func (n *node) allPolygons() []Polygon {
	allPolygons := append(make([]Polygon, 0, len(n.polygons)), n.polygons...)
	if n.front != nil {
		allPolygons = append(allPolygons, n.front.allPolygons()...)
	}
	if n.back != nil {
		allPolygons = append(allPolygons, n.back.allPolygons()...)
	}
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
		splitPolygon(p, n.plane, n.polygons, n.polygons, front, back)
	}
	if 0 < len(front) {
		if n.front != nil {
			n.front = &node{plane: Plane{}, front: nil, back: nil}
		}
		n.front.build(front)
	}
	if 0 < len(back) {
		if n.back != nil {
			n.back = &node{plane: Plane{}, front: nil, back: nil}
		}
		n.back.build(back)
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
