package todos

type NodeTodos int

const (
	UnkTag NodeTodos = iota
	GraphTag
	NodeTag
	GroupTag
	DataTag
	EdgeTag
	DatagroupTag
)

var SupportTag = map[string]NodeTodos{
	"graph":GraphTag,
	"node":      NodeTag,
	"group":     GroupTag,
	"datagroup": DatagroupTag,
	"edge":      EdgeTag,
	"data":      DataTag,
}

type Node struct {
	Id       int
	Tag      NodeTodos
	Text     string
	AtrList  []Atr
	Children []*Node
}

type Atr struct {
	Name  string
	Space string
	Value string
}
