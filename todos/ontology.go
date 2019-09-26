package todos

import "strings"

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
//
//const(
//	SubClassGroup="SubClass"
//	DomainGroup="Domain"
//	IndividualGroup="Individual"
//	PropertyGroup="Property"
//	DefaultGroup="Default"
//	CollectionGroup="Collection"
//)
type EdgeTagName struct {
	SubClassGroup   string
	DomainGroup     string
	IndividualGroup string
	PropertyGroup   string
	DefaultGroup    string
	CollectionGroup string
}

var SupportTag = map[string]NodeTodos{
	"graph":     GraphTag,
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

type Ontology struct {
	cnt int
	tag []Node
}

type stk struct {
	Id  int
	Tag NodeTodos
}

func (t *Ontology) findById(id int) *Node {
	for i, _ := range t.tag {
		if t.tag[i].Id == id {
			return &t.tag[i]
		}
	}
	return nil
}

func (t *Ontology) FindByTag(tag NodeTodos, skip int) *Node {
	for i, _ := range t.tag {
		if t.tag[i].Tag == tag {
			if skip == 0 {
				return &t.tag[i]
			}
			skip--
		}
	}
	return nil
}

func (t *Ontology) FindByAttr(tag NodeTodos, name string, value string, skip int) (int, *Node) {
	for{
	n := t.FindByTag(tag, skip)
	if n == nil {
		break
	}
	for idx, a := range n.AtrList {
		if strings.EqualFold(a.Name, name) && strings.EqualFold(a.Value, value) {
			return idx, n
		}
	}
	skip++
	}
	return 0, nil
}

func (n *Node) GetAtrValue(name string) string {
	for _, a := range n.AtrList {
		if strings.EqualFold(a.Name, name) {
			return a.Value
		}
	}
	return ""
}
