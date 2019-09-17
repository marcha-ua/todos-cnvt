package todos

import (
	"encoding/xml"
	"fmt"
	"github.com/golang-collections/collections/stack"
	"io"
	"strings"
)

func (t *Ontology) OntologyParser(r io.Reader) (err error) {
	decoder := xml.NewDecoder(r)
	//	total := 0
	stckTag := stack.New()
	currentTag := UnkTag
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		//fmt.Println(token)
		switch tag := token.(type) {
		case xml.StartElement:

			fmt.Println("StartElement ", tag)
			node, ok := SupportTag[strings.ToLower(tag.Name.Local)]
			if !ok {
				continue
			}

			currentTag = node
			id, err := t.NodeParser(token, node)
			if err != nil {
				fmt.Println(err)
			}
			if stckTag.Len() > 0 {
				lastTag := stckTag.Peek().(stk)
				err := t.AddTagLink(lastTag.Id, id)
				if err != nil {

				}
			}
			stckTag.Push(stk{id, currentTag})

		case xml.EndElement:
			node, ok := SupportTag[strings.ToLower(tag.Name.Local)]
			if !ok {
				continue
			}
			lastTag := stckTag.Peek().(stk)
			if lastTag.Tag == node {
				stckTag.Pop()
			}

			fmt.Println("EndElement ", tag)
		case xml.Attr:
			fmt.Println("Attr ", tag)
		case xml.Comment:
			fmt.Println("Comment ", tag)
		case xml.Name:
			fmt.Println("Name ", tag)

		case xml.CharData:
			str := string([]byte(tag))
			fmt.Println("CharData ", str)
			lastTag := stckTag.Peek().(stk)
			err := t.CharDataParser(token, lastTag.Id)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return
}

func (t *Ontology) NodeParser(token xml.Token, node NodeTodos) (int, error) {
	tag := token.(xml.StartElement)
	e := Node{}
	e.Id = t.cnt
	t.cnt++
	e.Tag = node

	for _, atr := range tag.Attr {
		a := Atr{atr.Name.Local, atr.Name.Space, atr.Value}
		e.AtrList = append(e.AtrList, a)
	}

	t.tag = append(t.tag, e)

	return e.Id, nil
}

func (t *Ontology) CharDataParser(token xml.Token, id int) (err error) {
	currentTag := token.(xml.CharData)
	str := cleanCharData(currentTag)
	if len(str) > 0 {
		n := t.findById(id)
		if n == nil {
			return
		}
		n.Text += str
	}
	return
}

func (t *Ontology) AddTagLink(parent int, child int) (err error) {
	p := t.findById(parent)
	if p == nil {
		return
	}
	c := t.findById(child)
	if c == nil {
		return
	}

	p.Children = append(p.Children, c)
	return
}

func StripChars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}

var stripSymbols = "\n\t"

func cleanCharData(data xml.CharData) string {
	str := string([]byte(data))
	str = StripChars(strings.TrimSpace(str), stripSymbols)
	if len(str) > 0 {
		return str
	}
	return ""
}
