package htmlvdom

type Element struct {
	TagName  string
	Attrs    map[string]string
	Children []*Element
	Parent   *Element
}
