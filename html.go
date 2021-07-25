package htmlvdom

import (
	"fmt"
	"html"
	"strings"
)

type Element struct {
	TagName  string
	Attrs    map[string]string
	Children []*Element
	Value    string
	parent   *Element
}

func CreateElement(tagName string) *Element {
	return &Element{
		TagName: tagName,
		Attrs:   make(map[string]string),
	}
}

func (el *Element) AppendChild(child *Element) {
	if child.parent != nil {
		child.parent.RemoveChild(child)
	}
	if child.TagName == "_text_node_" &&
		len(el.Children) == 1 &&
		el.Children[0].TagName == "_text_node_" {
		el.Children[0].Value += child.Value
	}
	el.Children = append(el.Children, child)
	child.parent = el
}

func (el *Element) SetAttribute(key, value string) {
	el.Attrs[key] = value
}

var ErrElementAttributeNotFound = fmt.Errorf("element attribute not found")

func (el *Element) GetAttribute(key string) (string, error) {
	value, ok := el.Attrs[key]
	if !ok {
		return "", ErrElementAttributeNotFound
	}
	return value, nil
}

func (el *Element) HasChildNodes() bool {
	return len(el.Children) > 0
}

func (el *Element) Clone(deep bool) *Element {
	if deep {
		clone := &Element{
			TagName: el.TagName,
			Attrs:   make(map[string]string),
		}
		for k, v := range el.Attrs {
			clone.Attrs[k] = v
		}
		for _, child := range el.Children {
			clone.AppendChild(child.Clone(deep))
		}
		return clone
	}
	attrs := make(map[string]string)
	for k, v := range el.Attrs {
		attrs[k] = v
	}
	return &Element{
		TagName: el.TagName,
		Attrs:   attrs,
	}
}

func (el *Element) RemoveChild(child *Element) {
	for i, c := range el.Children {
		if c == child {
			el.Children = append(el.Children[:i], el.Children[i+1:]...)
			child.parent = nil
			return
		}
	}
}

func (el *Element) SetInnerText(text string) {
	textNode := CreateElement("_text_node_")
	textNode.Value = html.EscapeString(text)
	for _, child := range el.Children {
		el.RemoveChild(child)
	}
	el.Children = el.Children[:0]
	el.AppendChild(textNode)
}

func (el *Element) CreateTextNode(text string) *Element {
	textNode := CreateElement("_text_node_")
	textNode.Value = html.EscapeString(text)
	return textNode
}

func (el *Element) String() string {
	if el.TagName == "_text_node_" {
		return el.Value
	}
	var b strings.Builder
	b.WriteRune('<')
	b.WriteString(el.TagName)
	b.WriteString(" ")
	for k, v := range el.Attrs {
		b.WriteString(k)
		b.WriteString("=\"")
		b.WriteString(html.EscapeString(v))
		b.WriteString("\" ")
	}
	b.WriteString(">")
	for _, child := range el.Children {
		b.WriteString(child.String())
	}
	b.WriteString("</")
	b.WriteString(el.TagName)
	b.WriteString(">")
	return b.String()
}
