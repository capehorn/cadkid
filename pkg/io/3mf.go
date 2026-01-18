package io

import (
	"capehorn/cadkid/lang"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type MFOptionalBool int8

const (
	Undefined MFOptionalBool = iota
	True
	False
)

type MFWriter struct {
	writer    io.Writer
	elements  lang.Stack[elementStringer]
	useIndent bool
	indent    string
}

type element interface {
	getParent() *element
	getWriter() *MFWriter
	getName() string
}

type elementStringer interface {
	toString() string
}

func NewMFWriter(writer io.Writer, indent uint8) *MFWriter {
	nfWriter := &MFWriter{writer: writer, elements: lang.Stack[elementStringer]{}, useIndent: 0 < indent, indent: strings.Repeat(" ", int(indent))}
	nfWriter.write("<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	return nfWriter
}

func (mw *MFWriter) writeElement(element string, attrs string, selfClosing bool) *MFWriter {
	return mw.writeElementWithText(element, attrs, "", selfClosing)
}

func (mw *MFWriter) writeElementWithText(element string, attrs string, text string, selfClosing bool) *MFWriter {
	// don't use indent for the first element
	if mw.useIndent {
		mw.write("\n")
		if element != "model" {
			mw.write(strings.Repeat(mw.indent, mw.elements.Length()))
		}
	}
	mw.write("<" + element)
	if attrs == "" {
		if text == "" {
			if selfClosing {
				mw.write("/>")
			} else {
				mw.write(">")
			}
		} else {
			mw.write(">" + text)
		}
	} else {
		if text == "" {
			if selfClosing {
				mw.write(attrs + "/>")
			} else {
				mw.write(attrs + ">")
			}
		} else {
			mw.write(attrs + ">" + text)
		}
	}
	return mw
}

func (mw *MFWriter) writeCloseElement(element string) *MFWriter {
	mw.write("</" + element + ">")
	return mw
}

func attr(name, value string) string {
	return name + "=\"" + value + "\""
}

func attrWithDefault(name, value string, defaultValue string) string {
	if value == "" {
		return name + "=\"" + defaultValue + "\""
	}
	return name + "=\"" + value + "\""
}

func (mw *MFWriter) Done() {
	if !mw.elements.IsEmpty() {
		for i := mw.elements.Length() - 1; 0 <= i; i-- {
			mw.write("\n")
			elStringer, _ := mw.elements.Pop()
			mw.write(strings.Repeat(mw.indent, i))
			mw.writeCloseElement(elStringer.toString())
		}
	}
}

func (mw *MFWriter) write(s string) *MFWriter {
	return mw.writeBytes([]byte(s))
}

func (mw *MFWriter) writeBytes(b []byte) *MFWriter {
	_, err := mw.writer.Write(b)
	if err != nil {
		contextErr := fmt.Errorf("failed to write data in 3mf format %w", err)
		fmt.Println(contextErr)
		return mw
	}
	return mw
}

type MFResourceId = uint32
type MFUnit = string

const (
	Micron     MFUnit = "micron"
	Millimeter MFUnit = "millimeter"
	Centimeter MFUnit = "centimeter"
	Inch       MFUnit = "inch"
	Foot       MFUnit = "foot"
	Meter      MFUnit = "meter"
)

// MODEL

type MFModel struct {
	writer *MFWriter
	parent *element
}

func (m MFModel) toString() string {
	return "model"
}

type MFModelAttr struct {
	Unit                  MFUnit
	RequiredExtensions    string
	RecommendedExtensions string
	Lang                  string
	Namespaces            map[string]string
}

func (m MFModelAttr) toString() string {
	str := strings.Builder{}
	str.WriteString(attrWithDefault(" unit", m.Unit, Millimeter))
	str.WriteString(attrWithDefault(" xml:lang", m.Lang, "en-us"))

	if m.RecommendedExtensions != "" {
		str.WriteString(attr(" requiredextensions", m.RequiredExtensions))
	}
	if m.RecommendedExtensions != "" {
		str.WriteString(attr(" recommendedextensions", m.RecommendedExtensions))
	}
	if 0 < len(m.Namespaces) {
		for k, v := range m.Namespaces {
			str.WriteString(attr(" xmlns:"+k, v))
		}
	}

	return str.String()
}

// METADATA

type MFMetadata struct {
	writer *MFWriter
}

func (m MFMetadata) toString() string {
	return "metadata"
}

type MFMetadataAttr struct {
	Name     string
	Preserve MFOptionalBool
	Type     string
}

func (m MFMetadataAttr) toString() string {
	str := strings.Builder{}
	str.WriteString(attr(" name", m.Name))
	switch m.Preserve {
	case Undefined:
		break
	case True:
		str.WriteString(attr(" preserve", "true"))
	case False:
		str.WriteString(attr(" preserve", "false"))
	}
	if m.Type != "" {
		str.WriteString(attr(" type", m.Type))
	}
	return str.String()
}

// RESOURCES

type MFResources struct {
	writer *MFWriter
	parent any
}

func (m MFResources) toString() string {
	return "resources"
}

type MFBaseMaterials struct {
	writer *MFWriter
	parent any
}

func (m MFBaseMaterials) toString() string {
	return "basematerials"
}

type MFBaseMaterialsAttr struct {
	Id uint32
}

func (m MFBaseMaterialsAttr) toString() string {
	return attr(" id", strconv.Itoa(int(m.Id)))
}

// OBJECT

type MFObject struct {
	writer *MFWriter
	parent any
}

func (m MFObject) toString() string {
	return "object"
}

type MFObjectType = string

const (
	Model        MFObjectType = "model"
	SolidSupport MFObjectType = "solidsupport"
	Support      MFObjectType = "support"
	Surface      MFObjectType = "surface"
	Other        MFObjectType = "other"
)

type MFObjectAttr struct {
	Id         uint32 // required
	ObjectType MFObjectType
	Thumbnail  string
	PartNumber string
	Name       string
	Pid        string
	PIndex     string
}

func (m MFObjectAttr) toString() string {
	str := strings.Builder{}
	str.WriteString(attr(" id", strconv.Itoa(int(m.Id))))
	str.WriteString(attrWithDefault(" type", m.ObjectType, Model))
	if m.Thumbnail != "" {
		str.WriteString(attr(" thumbnail", m.Thumbnail))
	}
	if m.PartNumber != "" {
		str.WriteString(attr(" partNumber", m.PartNumber))
	}
	if m.Name != "" {
		str.WriteString(attr(" name", m.Name))
	}
	if m.Pid != "" {
		str.WriteString(attr(" pid", m.Pid))
	}
	if m.PIndex != "" {
		str.WriteString(attr(" pindex", m.PIndex))
	}
	return str.String()
}

// MESH

type MFMesh struct {
	writer *MFWriter
	parent any
}

func (m MFMesh) toString() string {
	return "mesh"
}

// COMPONENTS

type MFComponents struct {
	writer *MFWriter
	parent any
}

func (m MFComponents) toString() string {
	return "components"
}

//type MFComponentAttr struct {
//	ObjectId  MFResourceId
//	Transform [12]float64
//}

// VERTICES

type MFVertices struct {
	writer *MFWriter
	parent any
}

func (m MFVertices) toString() string {
	return "vertices"
}

// TRIANGLES

type MFTriangles struct {
	writer *MFWriter
	parent any
}

func (m MFTriangles) toString() string {
	return "triangles"
}

// TRIANGLE

type MFTriangleAttr struct {
	P1, P2, P3 uint32
	Pid        string
}

func (m MFTriangleAttr) toString() string {
	str := strings.Builder{}
	if m.P1 != 0 {
		str.WriteString(attr(" p1", strconv.Itoa(int(m.P1))))
	}
	if m.P2 != 0 {
		str.WriteString(attr(" p2", strconv.Itoa(int(m.P2))))
	}
	if m.P3 != 0 {
		str.WriteString(attr(" p3", strconv.Itoa(int(m.P3))))
	}
	if m.Pid != "" {
		str.WriteString(attr(" pid", m.Pid))
	}
	return str.String()
}

// BUILD

type MFBuild struct {
	writer *MFWriter
	parent any
}

func (m MFBuild) toString() string {
	return "build"
}

// ITEM

type MFItem struct {
	writer *MFWriter
	parent any
}

func (m MFItem) toString() string {
	return "item"
}

type MFItemAttr struct {
	Transform  [12]float64 // Row-major affine transformation matrix, column four [ 0, 0, 0, 1 ] is not present
	PartNumber string
}

func (m MFItemAttr) toString() string {
	return "itemAttr, not implemented yet"
}

// Construction of a 3mf

func (mw *MFWriter) Model(attr MFModelAttr) MFModel {
	model := MFModel{writer: mw, parent: nil}
	mw.elements.Push(model)
	mw.writeElement(model.toString(), attr.toString(), false)
	return model
}

func (m MFModel) Metadata(text string, attr MFMetadataAttr) MFMetadata {
	metadata := MFMetadata{writer: m.writer}
	m.writer.writeElementWithText("metadata", attr.toString(), text, false).writeCloseElement("metadata")
	return metadata
}

func (m MFModel) Resources() MFResources {
	resources := MFResources{writer: m.writer, parent: m}
	m.writer.elements.Push(resources)
	m.writer.writeElement(resources.toString(), "", false)
	return resources
}

func (m MFResources) BaseMaterials(id uint32) MFBaseMaterials {
	baseMaterials := MFBaseMaterials{writer: m.writer, parent: m}
	m.writer.elements.Push(baseMaterials)
	m.writer.writeElement(baseMaterials.toString(), MFBaseMaterialsAttr{Id: id}.toString(), false)
	return baseMaterials
}

type SRGB = string

func (m MFBaseMaterials) Base(name string, displayColor SRGB) MFBaseMaterials {
	m.writer.writeElement("base", strings.Join([]string{
		attr(" name", name),
		attr(" displaycolor", displayColor)}, ""),
		true)
	return m
}

func (m MFResources) Object(attr MFObjectAttr) MFObject {
	object := MFObject{writer: m.writer, parent: m}
	m.writer.elements.Push(object)
	m.writer.writeElement(object.toString(), attr.toString(), false)
	return object
}

func (m MFObject) Mesh() MFMesh {
	mesh := MFMesh{writer: m.writer, parent: m}
	m.writer.elements.Push(mesh)
	m.writer.writeElement(mesh.toString(), "", false)
	return mesh
}

func (m MFMesh) Vertices() MFVertices {
	vertices := MFVertices{writer: m.writer, parent: m}
	m.writer.elements.Push(vertices)
	m.writer.writeElement(vertices.toString(), "", false)
	return vertices
}

func (m MFObject) Components() MFComponents {
	components := MFComponents{writer: m.writer, parent: m}
	m.writer.elements.Push(components)
	m.writer.writeElement(components.toString(), "", false)
	return components
}

func (m MFComponents) Component(objectId MFResourceId, transform []float64) MFComponents {
	builder := strings.Builder{}
	builder.WriteString(attr(" objectid", strconv.Itoa(int(objectId))))
	if transform != nil {
		builder.WriteString(attr("transform",
			strings.Join(lang.Map(transform,
				func(x float64) string {
					return strconv.FormatFloat(x, 'f', -1, 64)
				}), ", ")))
	}
	m.writer.writeElement("component", builder.String(), true)
	return m
}

func (m MFVertices) Vertex(x, y, z float64) MFVertices {
	m.writer.writeElement("vertex", strings.Join([]string{
		attr(" x", strconv.FormatFloat(x, 'f', -1, 64)),
		attr(" y", strconv.FormatFloat(y, 'f', -1, 64)),
		attr(" z", strconv.FormatFloat(z, 'f', -1, 64)),
	}, ""), true)
	return m
}

func (m MFMesh) Triangles() MFTriangles {
	triangles := MFTriangles{writer: m.writer, parent: m}
	m.writer.elements.Push(triangles)
	m.writer.writeElement(triangles.toString(), "", false)
	return triangles
}

func (m MFTriangles) Triangle(v1, v2, v3 uint32, attrs *MFTriangleAttr) MFTriangles {
	attrStr := strings.Join([]string{
		attr(" v1", strconv.Itoa(int(v1))),
		attr(" v2", strconv.Itoa(int(v2))),
		attr(" v3", strconv.Itoa(int(v3))),
	}, "")

	if attrs != nil {
		attrStr += attrs.toString()
	}

	m.writer.writeElement("triangle", attrStr, true)
	return m
}

func (m MFModel) Build() MFBuild {
	build := MFBuild{writer: m.writer, parent: m}
	m.writer.elements.Push(build)
	m.writer.writeElement(build.toString(), "", false)
	return build
}

func (m MFBuild) Item(objectId uint32, attr MFItemAttr) MFItem {
	item := MFItem{writer: m.writer, parent: m}
	m.writer.elements.Push(item)
	m.writer.writeElement(item.toString(), attr.toString(), false)
	return item
}

type MFReader struct {
}
