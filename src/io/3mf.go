package io

import (
	"capehorn/cadkid/lang"
	"fmt"
	"io"
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

type elementStringer interface {
	toString() string
}

type attrsStringer interface {
	toString() string
}

func NewMFWriter(writer io.Writer, indent uint8) *MFWriter {
	nfWriter := &MFWriter{writer: writer, elements: lang.Stack[elementStringer]{}, useIndent: 0 < indent, indent: strings.Repeat(" ", int(indent))}
	nfWriter.write("<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	return nfWriter
}

func (mw *MFWriter) writeElement(elStringer elementStringer, attrsStringer attrsStringer) *MFWriter {
	return mw.writeElementWithText(elStringer, attrsStringer, "")
}

func (mw *MFWriter) writeElementWithText(elStringer elementStringer, attrsStringer attrsStringer, text string) *MFWriter {
	// don't use indent for the first element
	if mw.useIndent {
		mw.write("\n")
		switch elStringer.(type) {
		case MFModel:
			break
		default:
			mw.write(mw.indent)
		}
	}
	if attrsStringer == nil {
		mw.write("<" + elStringer.toString() + ">" + text)
	} else {
		mw.write("<" + elStringer.toString() + " " + attrsStringer.toString() + ">" + text)
	}
	return mw
}

func (mw *MFWriter) writeCloseElement(elStringer elementStringer) *MFWriter {
	mw.write("</" + elStringer.toString() + ">")
	return mw
}

func attr(name, value string) string {
	return name + "=\"" + value + "\""
}

func (mw *MFWriter) Done() {
	if !mw.elements.IsEmpty() {
		for i := mw.elements.Length() - 1; 0 <= i; i-- {
			mw.write("\n")
			elStringer, _ := mw.elements.Pop()
			mw.write(strings.Repeat(mw.indent, i))
			mw.writeCloseElement(elStringer)
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

type MFUnit int8

const (
	Micron MFUnit = iota
	Millimeter
	Centimeter
	Inch
	Foot
	Meter
)

type MFModel struct {
	writer *MFWriter
}

func (m MFModel) toString() string {
	return "model"
}

type MFModelAttr struct {
	Unit                  MFUnit
	RequiredExtensions    string
	RecommendedExtensions string
	Lang                  string
}

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
	str.WriteString(attr("name", m.Name))
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

type MFResources struct {
	writer *MFWriter
}

type MFBaseMaterials struct {
	writer *MFWriter
}

type MFObject struct {
	writer *MFWriter
}

type MFObjectAttr struct {
	Id uint32
}

type MFMesh struct {
	writer *MFWriter
}

func (mw *MFWriter) Model(attr MFModelAttr) MFModel {
	model := MFModel{writer: mw}
	mw.elements.Push(model)
	mw.writeElement(model, nil)
	//mw.write("<model>")
	return model
}

func (m MFModel) Metadata(text string, attr MFMetadataAttr) MFMetadata {
	metadata := MFMetadata{writer: m.writer}
	m.writer.writeElementWithText(metadata, attr, text).writeCloseElement(metadata)
	return metadata
}

func (m MFModel) Resources() MFResources {
	resources := MFResources{writer: m.writer}
	//m.writer.elements.Push(resources)
	//m.writer.write("<resources>")
	return resources
}

func (m MFResources) BaseMaterials(id uint32) MFBaseMaterials {
	baseMaterials := MFBaseMaterials{writer: m.writer}
	//m.writer.elements.Push(baseMaterials)
	//m.writer.write("<basematerials id=\"" + fmt.Sprint(id) + "\">")
	return baseMaterials
}

type SRGB = string

func (m MFBaseMaterials) Base(name string, displayColor SRGB) MFBaseMaterials {
	//m.writer.write("<base name=\"" + name + "\" displaycolor=\"" + displayColor + "\"/>")
	return m
}

func (m MFResources) Object(attr MFObjectAttr) MFObject {
	object := MFObject{writer: m.writer}
	//m.writer.elements.Push(object)
	//m.writer.write("<object>")
	return object
}

func (m MFObject) Mesh() MFMesh {
	mesh := MFMesh{writer: m.writer}
	//m.writer.elements.Push(mesh)
	//m.writer.write("<mesh>")
	return mesh
}

func (m MFMesh) Vertex(x, y, z float64) MFMesh {
	// TODO write x, y, z
	return m
}

func (m MFMesh) Triangle(v1, v2, v3 uint32) MFMesh {
	// TODO write x, y, z
	return m
}

type MFReader struct {
}

func kvAttributesToString(kvAttributes ...string) string {
	if len(kvAttributes) == 0 {
		return ""
	}
	if len(kvAttributes)%2 != 0 {
		err := fmt.Errorf("odd number of kbAttributes")
		fmt.Println(err)
		return ""
	}
	attrs := make([]string, len(kvAttributes)/2)
	for i := 0; i < len(kvAttributes); i += 2 {
		attrs = append(attrs, "\""+kvAttributes[i]+"\"=\""+kvAttributes[i+1]+"\"")
	}
	return strings.Join(attrs, " ")
}
