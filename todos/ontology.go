package todos

type NodeTodos int

const (
	NodeTag NodeTodos = iota
	GroupTag
	DataTag
	EdgeTag
	DatagroupTag
)

type Node struct {
	id       int
	tag      NodeTodos
	text     string
	atr      []Atr
	children []*Node
}

type Atr struct {
	name  string
	value string
}
