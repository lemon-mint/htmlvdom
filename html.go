package htmlvdom

import (
	"fmt"
	"html"
	"strings"
	"sync"
)

var elPool = sync.Pool{New: func() interface{} {
	return &Element{
		TagName: "__textnode__",
		Attrs:   make(map[string]string),
	}
}}

type Element struct {
	TagName  string
	Attrs    map[string]string
	Children []*Element
	Value    string
	parent   *Element
}

func CreateElement(tagName string) *Element {
	el := elPool.Get().(*Element)
	el.TagName = tagName
	return el
}

func (el *Element) AppendChild(child *Element) {
	if el.TagName == "__textnode__" {
		return
	}
	if child.parent != nil {
		child.parent.RemoveChild(child)
	}
	el.Children = append(el.Children, child)
	child.parent = el
}

func (el *Element) SetAttribute(key, value string) {
	el.Attrs[key] = value
}

var ErrAttributeNotFound = fmt.Errorf("attribute not found")

func (el *Element) GetAttribute(key string) (string, error) {
	value, ok := el.Attrs[key]
	if !ok {
		return "", ErrAttributeNotFound
	}
	return value, nil
}

func (el *Element) RemoveAttribute(key string) {
	delete(el.Attrs, key)
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

func (el *Element) Remove() {
	if el.parent != nil {
		el.parent.RemoveChild(el)
	}
}

func (el *Element) SetInnerText(text string) {
	textNode := CreateTextNode(text)
	for _, child := range el.Children {
		el.RemoveChild(child)
	}
	el.Children = el.Children[:0]
	el.AppendChild(textNode)
}

func CreateTextNode(text string) *Element {
	textNode := CreateElement("__textnode__")
	textNode.Value = html.EscapeString(text)
	return textNode
}

func (el *Element) Destory() {
	if el.parent != nil {
		el.parent.RemoveChild(el)
	}
	for _, child := range el.Children {
		child.Destory()
	}
	for key := range el.Attrs {
		delete(el.Attrs, key)
	}
	el.Children = el.Children[:0]
	el.parent = nil
	el.Value = ""
	el.TagName = ""
	elPool.Put(el)
}

func (el *Element) String() string {
	if el.TagName == "__textnode__" {
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
