package owl

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"todos-cnvt/todos"
)

type owlOnto struct {
	Annotation map[string]AnnotationTag
	Class      map[string]ClassTag
	Individual map[string]NamedIndividualTag
	Collection map[string]ClassCollTag
}

func TodosFileBuilder(cfg string, name string, ontology *todos.Ontology, w io.Writer) (err error) {
	ont := owlOnto{}
	ont.Class = make(map[string]ClassTag, 0)
	ont.Individual = make(map[string]NamedIndividualTag, 0)
	ont.Collection = make(map[string]ClassCollTag, 0)
	ont.Annotation = make(map[string]AnnotationTag, 0)

	h := BuildHeader(name)

	//find root node
	n := ""
	_, r := ontology.FindByAttr(todos.NodeTag, "nodeName", "Thing", 0)
	if r != nil {

	} else {
		_, r = ontology.FindByAttr(todos.NodeTag, "nclass", "Thing", 0)
		n = "#" + r.GetAtrValue("nodeName")

	}
	root := ClassTag{
		About: n,
	}
	root.Label.InnerXML = root.About
	root.Label.Datatype = XMLNS_XSD + "string"

	guid := AnnotationTag{}
	guid.About = h.Base + "guid"
	ont.Annotation["guid"] = guid

	edge := ReadConfigEdge(cfg)

	//ont.FillNodeList(ontology)
	ont.EdgeNodeAnalyze(edge, ontology)

	enc := xml.NewEncoder(w)
	enc.Indent("  ", "    ")

	for _, n := range ont.Annotation {
		h.Body = append(h.Body, n)
	}

	for _, n := range ont.Class {
		h.Body = append(h.Body, n)
	}

	for _, n := range ont.Individual {
		h.Body = append(h.Body, n)
	}

	for _, n := range ont.Collection {
		h.Body = append(h.Body, n)
	}
	if err = enc.Encode(h); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	return
}

func (o *owlOnto) FillNodeList(ontology *todos.Ontology) {
	skip := 0
	for {

		n := ontology.FindByTag(todos.NodeTag, skip)
		if n == nil {
			return
		}

		name := n.GetAtrValue("nodeName")

		if !strings.HasPrefix(name, "#") {
			name = "#" + name
		}
		c := ClassTag{
			About: name,
		}
		c.Label.InnerXML = c.About
		c.Label.Datatype = XMLNS_XSD + "string"

		if len(n.Children) > 0 {
			for _, v := range n.Children {
				if v.Tag == todos.DataTag {
					t := v.GetAtrValue("type")
					if len(t) > 0 && strings.EqualFold(t, "text") {
						c.Comment.CharData = v.Text
						c.Comment.Lang = "ua"
					} else {
						fmt.Println("[data] - value attr type not supported", t)
					}
				} else {
					fmt.Println("[Node] - Nested tag not supported", v.Tag)
				}
			}
		}
		o.Class[c.About] = c
		skip++
	}
}

func (o *owlOnto) EdgeNodeAnalyze(edgeName todos.EdgeTagName, ontology *todos.Ontology) {
	skip := 0

	for {
		n := ontology.FindByTag(todos.EdgeTag, skip)
		if n == nil {
			fmt.Println("Edge tags not found")
			return
		}

		group := n.GetAtrValue("group")
		if len(group) == 0 {
			guid := n.GetAtrValue("guid")
			fmt.Println("Attribute group not found", guid)
			continue
		}
		switch group {
		case edgeName.SubClassGroup:
			node1 := n.GetAtrValue("node1")
			o.AddClassNode(node1, ontology)
			node2 := n.GetAtrValue("node2")
			o.AddClassNode(node2, ontology)
			o.AddSubClassNode(node1, node2, ontology)
		case edgeName.DefaultGroup:
		case edgeName.CollectionGroup:
			node1 := n.GetAtrValue("node1")
			o.AddClassNode(node1, ontology)
			node2 := n.GetAtrValue("node2")
			o.AddClassCollNode(node2, ontology)
			o.AddSubClassNode(node1, node2, ontology)
			o.AddCollectionNode(node1, node2, ontology)

		case edgeName.DomainGroup:
			node1 := n.GetAtrValue("node1")
			o.AddClassNode(node1, ontology)
			node2 := n.GetAtrValue("node2")
			o.AddClassNode(node2, ontology)
			o.AddSubClassNode(node1, node2, ontology)
		case edgeName.IndividualGroup:
			node1 := n.GetAtrValue("node1")
			node2 := n.GetAtrValue("node2")
			o.AddIndividualNode(node1, node2, ontology)
		case edgeName.PropertyGroup:

		}
		skip++
	}
}

func (o *owlOnto) AddClassNode(name string, ontology *todos.Ontology) {

	_, n := ontology.FindByAttr(todos.NodeTag, "nodeName", name, 0)
	if n == nil {
		fmt.Println("Node not found", name)
		return
	}

	if !strings.HasPrefix(name, "#") {
		name = "#" + name
	}
	if _, ok := o.Class[name]; ok {
		return
	}
	c := ClassTag{
		About: name,
	}
	c.Label.InnerXML = c.About
	c.Label.Datatype = XMLNS_XSD + "string"
	guid := n.GetAtrValue("guid")
	c.Guid.DataType = GuidDataType
	c.Guid.CharData = guid

	if len(n.Children) > 0 {
		for _, v := range n.Children {
			if v.Tag == todos.DataTag {
				t := v.GetAtrValue("type")
				if len(t) > 0 && strings.EqualFold(t, "text") {
					c.Comment.CharData = v.Text
					c.Comment.Lang = "ua"
				} else {
					fmt.Println("[data] - value attr type not supported", t)
				}
			} else {
				fmt.Println("[Node] - Nested tag not supported", v.Tag)
			}
		}
	}
	o.Class[c.About] = c
}

func (o *owlOnto) AddSubClassNode(class string, subclass string, ontology *todos.Ontology) {
	if !strings.HasPrefix(class, "#") {
		class = "#" + class
	}

	sb := SubClassOfTag{}
	if !strings.HasPrefix(subclass, "#") {
		subclass = "#" + subclass
	}
	sb.Resource = subclass
	found := false
	if c, ok := o.Class[class]; ok {
		c.SubClass = append(c.SubClass, sb)
		o.Class[class] = c
		found = true
	}
	if c, ok := o.Collection[class]; ok {
		c.SubClass = append(c.SubClass, sb)
		o.Collection[class] = c
		found = true
	}
	if !found {
		fmt.Println("Subclass not found class", class, subclass)
	}
}

func (o *owlOnto) AddIndividualNode(name string, class string, ontology *todos.Ontology) {
	_, n := ontology.FindByAttr(todos.NodeTag, "nodeName", name, 0)
	if n == nil {
		fmt.Println("Node not found", name)
		return
	}
	if !strings.HasPrefix(name, "#") {
		name = "#" + name
	}
	_, ok := o.Individual[name]
	if ok {
		return
	}
	i := NamedIndividualTag{}
	i.About = name
	i.Label.Datatype = XMLNS_XSD + "string"
	if !strings.HasPrefix(class, "#") {
		class = "#" + class
	}
	i.Type.Resource = class

	guid := n.GetAtrValue("guid")
	i.Guid.DataType = GuidDataType
	i.Guid.CharData = guid

	o.Individual[name] = i
}

func (o *owlOnto) AddClassCollNode(name string, ontology *todos.Ontology) {

	_, n := ontology.FindByAttr(todos.NodeTag, "nodeName", name, 0)
	if n == nil {
		fmt.Println("Node not found", name)
		return
	}

	if !strings.HasPrefix(name, "#") {
		name = "#" + name
	}

	if _, ok := o.Class[name]; ok {
		delete(o.Class, name)
	}

	if _, ok := o.Collection[name]; ok {
		return
	}

	c := ClassCollTag{
		About: name,
	}

	c.Label.InnerXML = c.About
	c.Label.Datatype = XMLNS_XSD + "string"

	guid := n.GetAtrValue("guid")
	c.Guid.DataType = GuidDataType
	c.Guid.CharData = guid

	if len(n.Children) > 0 {
		for _, v := range n.Children {
			if v.Tag == todos.DataTag {
				t := v.GetAtrValue("type")
				if len(t) > 0 && strings.EqualFold(t, "text") {
					c.Comment.CharData = v.Text
					c.Comment.Lang = "ua"
				} else {
					fmt.Println("[data] - value attr type not supported", t)
				}
			} else {
				fmt.Println("[Node] - Nested tag not supported", v.Tag)
			}
		}
	}

	o.Collection[c.About] = c
}

func (o *owlOnto) AddCollectionNode(name string, class string, ontology *todos.Ontology) {
	if !strings.HasPrefix(class, "#") {
		class = "#" + class
	}
	c, ok := o.Collection[class]
	if !ok {
		fmt.Println("After added class node not found", class)
		return
	}

	if !strings.HasPrefix(name, "#") {
		name = "#" + name
	}

	c.Collection.ParseType = "Collection"
	d := DescriptionTag{}
	d.About = name
	c.Collection.Description = append(c.Collection.Description, d)

	o.Collection[class] = c
}
