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
	Class map[string]ClassTag
}

func TodosFileBuilder(name string, ontology *todos.Ontology, w io.Writer) (err error) {

	ont := owlOnto{}
	ont.Class = make(map[string]ClassTag, 0)
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


	ont.FillNodeList(ontology)

	enc := xml.NewEncoder(w)
	enc.Indent("  ", "    ")
	if err = enc.Encode(h); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	if err = enc.Encode(root); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	for _,n:=range ont.Class{
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
					}else{
						fmt.Println("[data] - value attr type not supported",t)
					}
				}else {
					fmt.Println("[Node] - Nested tag not supported", v.Tag)
				}
			}
		}
		o.Class[c.About] = c
		skip++
	}
}
