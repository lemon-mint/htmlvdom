package htmlvdom

import "reflect"

func Diff(a, b *Element) *Difference {
	d := &Difference{
		Diff: make([]*Operation, 0, 1024),
	}
	if a == nil && b == nil {
		return d
	}
	diffNodes(a, b, d, b.ID)
	return d
}

func diffNodes(a, b *Element, d *Difference, TargetID uint64) {
	var modified [][2]*Element
	var deleted, added []*Element

	if a != nil {
		TargetID = a.ID
	}

	if a != nil {
		if a.XXHash == b.XXHash {
			return // no change
		}
		for _, c := range a.Children {
			i, ok := b.IndexChild(c)
			if !ok {
				deleted = append(deleted, c)
			} else {
				if c.XXHash != b.Children[i].XXHash {
					modified = append(modified, [2]*Element{c, b.Children[i]})
				}
			}
		}
		for _, c := range b.Children {
			_, ok := a.IndexChild(c)
			if !ok {
				added = append(added, c)
			}
		}
	} else {
		added = b.Children
	}

	for _, el := range deleted {
		d.Diff = append(d.Diff, &Operation{
			Type:   OP_REMOVE_CHILD,
			Target: TargetID,
			DeleteInfo: &DeleteInfo{
				ID: el.ID,
			},
		})
	}

	for _, el := range added {
		d.Diff = append(d.Diff, &Operation{
			Type:   OP_CREATE_ELEMENT,
			Target: TargetID,
			CreateInfo: &CreateInfo{
				TagName: el.TagName,
				NewID:   el.ID,
			},
		})
		diffNodes(nil, el, d, el.ID)
	}

	for _, el := range modified {
		diffNodes(el[0], el[1], d, TargetID)
	}

	var setKV, removeKV [][2]string

	if a != nil {
		if !reflect.DeepEqual(a.Attrs, b.Attrs) {
			for k, v := range a.Attrs {
				_, ok := b.Attrs[k]
				if !ok {
					removeKV = append(removeKV, [2]string{k, v})
				}
			}
			for k, v := range b.Attrs {
				if aval, ok := a.Attrs[k]; !ok || aval != v {
					setKV = append(setKV, [2]string{k, v})
				}
			}
		}
	} else {
		for k, v := range b.Attrs {
			setKV = append(setKV, [2]string{k, v})
		}
	}
	for _, kv := range setKV {
		d.Diff = append(d.Diff, &Operation{
			Type:   OP_SET_ATTRIBUTE,
			Target: TargetID,
			SetAttributeInfo: &SetAttributeInfo{
				Key: kv[0],
				Val: kv[1],
			},
		})
	}
	for _, kv := range removeKV {
		d.Diff = append(d.Diff, &Operation{
			Type:   OP_REMOVE_ATTRIBUTE,
			Target: TargetID,
			RemoveAttributeInfo: &RemoveAttributeInfo{
				Key: kv[0],
			},
		})
	}

	if a != nil {
		if a.Value != b.Value {
			d.Diff = append(d.Diff, &Operation{
				Type:   OP_SET_VALUE,
				Target: TargetID,
				SetValueInfo: &SetValueInfo{
					Val: b.Value,
				},
			})
		}
	} else if b.Value != "" && b.TagName == "__textnode__" {
		d.Diff = append(d.Diff, &Operation{
			Type:   OP_SET_VALUE,
			Target: TargetID,
			SetValueInfo: &SetValueInfo{
				Val: b.Value,
			},
		})
	}
}
