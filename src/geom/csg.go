package geom

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

type polygonType int8
const (
	COPLANAR polygonType = iota
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
	var polygonType = 0;
	var types = [];
	for (var i = 0; i < polygon.vertices.length; i++) {
		var t = this.normal.dot(polygon.vertices[i].pos) - this.w;
		var type = (t < -CSG.Plane.EPSILON) ? BACK : (t > CSG.Plane.EPSILON) ? FRONT : COPLANAR;
		polygonType |= type;
		types.push(type);
	}

	// Put the polygon in the correct list, splitting it when necessary.
	switch (polygonType) {
		case COPLANAR:
			(this.normal.dot(polygon.plane.normal) > 0 ? coplanarFront : coplanarBack).push(polygon);
			break;
		case FRONT:
			front.push(polygon);
			break;
		case BACK:
			back.push(polygon);
			break;
		case SPANNING:
			var f = [], b = [];
			for (var i = 0; i < polygon.vertices.length; i++) {
				var j = (i + 1) % polygon.vertices.length;
				var ti = types[i], tj = types[j];
				var vi = polygon.vertices[i], vj = polygon.vertices[j];
				if (ti != BACK) f.push(vi);
				if (ti != FRONT) b.push(ti != BACK ? vi.clone() : vi);
				if ((ti | tj) == SPANNING) {
					var t = (this.w - this.normal.dot(vi.pos)) / this.normal.dot(vj.pos.minus(vi.pos));
					var v = vi.interpolate(vj, t);
					f.push(v);
					b.push(v.clone());
				}
			}
			if (f.length >= 3) front.push(new CSG.Polygon(f, polygon.shared));
			if (b.length >= 3) back.push(new CSG.Polygon(b, polygon.shared));
			break;
	}
}