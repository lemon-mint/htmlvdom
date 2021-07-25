package htmlvdom

import (
	"encoding/binary"
	"fmt"
	"hash"
	"html"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/cespare/xxhash"
)

var globalIDCounter uint64

var elPool = sync.Pool{New: func() interface{} {
	return &Element{
		TagName: "__textnode__",
		Attrs:   make(map[string]string),
	}
}}

var hashPool = sync.Pool{New: func() interface{} {
	return xxhash.New()
}}

type Element struct {
	TagName  string
	Attrs    map[string]string
	Children []*Element
	Value    string
	ID       uint64
	XXHash   uint64
	parent   *Element
}

func (el *Element) InitHash() {
	h := hashPool.Get().(hash.Hash64)
	h.Write(StringToBytes(el.TagName))
	keys := make([]string, 0, len(el.Attrs))
	for k := range el.Attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		h.Write(StringToBytes(k))
		h.Write(StringToBytes(el.Attrs[k]))
	}
	h.Write(StringToBytes(el.Value))
	for _, child := range el.Children {
		var ch [8]byte
		binary.BigEndian.PutUint64(ch[:], child.XXHash)
		h.Write(ch[:])
	}
	el.XXHash = h.Sum64()
	hashPool.Put(h)
}

func (el *Element) Hash() {
	el.InitHash()
	if el.parent != nil {
		el.parent.Hash()
	}
}

func CreateElement(tagName string) *Element {
	el := elPool.Get().(*Element)
	el.GetNewID()
	el.TagName = tagName
	el.Hash()
	return el
}

func (el *Element) GetNewID() uint64 {
	id := atomic.AddUint64(&globalIDCounter, 1)
	el.ID = id
	return id
}

func (el *Element) AppendChild(child *Element) {
	defer el.Hash()
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
	defer el.Hash()
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
	defer el.Hash()
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
			XXHash:  el.XXHash,
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
		XXHash:  el.XXHash,
	}
}

func (el *Element) RemoveChild(child *Element) {
	defer el.Hash()
	for i, c := range el.Children {
		if c == child {
			el.Children = append(el.Children[:i], el.Children[i+1:]...)
			child.parent = nil
			return
		}
	}
}

func (el *Element) ReplaceChild(oldChild, newChild *Element) {
	for i, c := range el.Children {
		if c == oldChild {
			el.Children[i] = newChild
			newChild.parent = el
			oldChild.parent = nil
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
		child.parent = nil
	}
	el.Children = el.Children[:0]
	el.AppendChild(textNode)
}

func CreateTextNode(text string) *Element {
	textNode := CreateElement("__textnode__")
	textNode.Value = html.EscapeString(text)
	textNode.Hash()
	return textNode
}

func (el *Element) Destory() {
	el.Remove()
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
	el.Hash()
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

func (el *Element) Init() {
	el.GetNewID()
	for _, child := range el.Children {
		child.parent = el
		child.Init()
		child.InitHash()
	}
	el.InitHash()
}
