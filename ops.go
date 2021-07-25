package htmlvdom

type Operation struct {
	Type   OperationType
	Target uint64

	CreateInfo          *CreateInfo          `json:"CreateInfo,omitempty"`
	SetAttributeInfo    *SetAttributeInfo    `json:"SetAttributeInfo,omitempty"`
	DeleteInfo          *DeleteInfo          `json:"DeleteInfo,omitempty"`
	RemoveAttributeInfo *RemoveAttributeInfo `json:"RemoveAttributeInfo,omitempty"`
	SetValueInfo        *SetValueInfo        `json:"SetValueInfo,omitempty"`
}

type Difference struct {
	Diff []*Operation
}

type OperationType byte

const (
	OP_CREATE_ELEMENT OperationType = iota
	OP_CREATE_TEXT
	OP_SET_VALUE
	OP_SET_ATTRIBUTE
	OP_REMOVE_ATTRIBUTE
	OP_APPEND_CHILD
	OP_REMOVE_CHILD
	OP_REPLACE_CHILD
)

type CreateInfo struct {
	TagName string
	NewID   uint64
}

type SetAttributeInfo struct {
	Key string
	Val string
}

type DeleteInfo struct {
	ID uint64
}

type RemoveAttributeInfo struct {
	Key string
}

type SetValueInfo struct {
	Val string
}
