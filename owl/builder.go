package owl

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"todos-cnvt/todos"
)

const ()

type owlOnto struct {
	Class      map[string]ClassTag
	Individual map[string]NamedIndividualTag
	Collection map[string]ClassCollTag

}

func TodosFileBuilder(name string, ontology *todos.Ontology, w io.Writer) (err error) {

	ont := owlOnto{}
	ont.Class = make(map[string]ClassTag, 0)
	ont.Individual = make(map[string]NamedIndividualTag, 0)
	ont.Collection = make(map[string]ClassCollTag, 0)

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

	//ont.FillNodeList(ontology)
	ont.EdgeNodeAnalyze(ontology)

	enc := xml.NewEncoder(w)
	enc.Indent("  ", "    ")
	if err = enc.Encode(h); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	//if err = enc.Encode(root); err != nil {
	//	fmt.Printf("error: %v\n", err)
	//	return
	//}

	for _, n := range ont.Class {
		if err = enc.Encode(n); err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
	}
	for _, n := range ont.Individual {
		if err = enc.Encode(n); err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
	}
	for _, n := range ont.Collection {
		if err = enc.Encode(n); err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
	}

	//head, err := xmltree.Parse(todos)
	//if err != nil {
	//	return
	//}
	//head.SetAttr("", "xmlns:"+name, XMLNS+name+"#")
	//header := xmltree.MarshalIndent(head, "", "   ")

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

func (o *owlOnto) EdgeNodeAnalyze(ontology *todos.Ontology) {
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
		case todos.SubClassGroup:
			node1 := n.GetAtrValue("node1")
			o.AddClassNode(node1, ontology)
			node2 := n.GetAtrValue("node2")
			o.AddClassNode(node2, ontology)
 			o.AddSubClassNode(node1, node2, ontology)
		case todos.DefaultGroup:
		case todos.CollectionGroup:
			node1 := n.GetAtrValue("node1")
			o.AddClassNode(node1, ontology)
			node2 := n.GetAtrValue("node2")
			o.AddClassCollNode(node2, ontology)
			o.AddSubClassNode(node1, node2, ontology)
			o.AddCollectionNode(node1,node2,ontology)

		case todos.DomainGroup:
			node1 := n.GetAtrValue("node1")
			o.AddClassNode(node1, ontology)
			node2 := n.GetAtrValue("node2")
			o.AddClassNode(node2, ontology)
			o.AddSubClassNode(node1, node2, ontology)
		case todos.IndividualGroup:
			node1 := n.GetAtrValue("node1")
			node2 := n.GetAtrValue("node2")
			o.AddIndividualNode(node1,node2, ontology)
		case todos.PropertyGroup:

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

	//guid :=n.GetAtrValue("guid")

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
	c, ok := o.Class[class]
	if !ok {
		fmt.Println("Class node not found", class)
		return
	}
	sb := SubClassOfTag{}
	if !strings.HasPrefix(subclass, "#") {
		subclass = "#" + subclass
	}
	sb.Resource = subclass
	c.SubClass = append(c.SubClass, sb)
	o.Class[class]=c
}

func (o *owlOnto) AddIndividualNode(name string, class string, ontology *todos.Ontology) {
	if !strings.HasPrefix(name, "#") {
		name = "#" + name
	}
	_, ok := o.Individual[name]
	if ok {
		return
	}
	i:= NamedIndividualTag{}
	i.About=name
	i.Label.Datatype = XMLNS_XSD + "string"
	if !strings.HasPrefix(class, "#") {
		class = "#" + class
	}
i.Type.Resource=class

	o.Individual[name]=i
}

func (o *owlOnto) AddClassCollNode(name string, ontology *todos.Ontology) {

	_, n := ontology.FindByAttr(todos.NodeTag, "nodeName", name, 0)
	if n == nil {
		fmt.Println("Node not found", name)
		return
	}

	//guid :=n.GetAtrValue("guid")

	if !strings.HasPrefix(name, "#") {
		name = "#" + name
	}
	if _, ok := o.Class[name]; ok {
		return
	}
	c := ClassCollTag{
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

	coll:= c.Collection
	coll.ParseType = "Collection"
	d:=DescriptionTag{}
	d.About=name
	coll.Description = append(coll.Description,d)

	c.Collection = coll
	o.Collection[class] = c
}
