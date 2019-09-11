package todos

import (
	"encoding/xml"
	"fmt"
	"github.com/golang-collections/collections/stack"
	"io"
)

func OntologyParser(r io.Reader) (err error) {
	decoder := xml.NewDecoder(r)
	//	total := 0
	stckTag := stack.New()
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		//fmt.Println(token)
		switch currentTag := token.(type) {
		case xml.StartElement:
			stckTag.Push(currentTag.Name.Local)
			if currentTag.Name.Local == "entry" {
				// do what you need to do for each entry below
			}
			fmt.Println("StartElement ", currentTag)
		case xml.EndElement:
			lastTag := stckTag.Peek()
			if lastTag == currentTag.Name.Local {
				//todo save
			}
			stckTag.Pop()
			fmt.Println("EndElement ", currentTag)
		case xml.Attr:
			fmt.Println("Attr ", currentTag)
		case xml.Comment:
			fmt.Println("Comment ", currentTag)
		case xml.Name:
			fmt.Println("Name ", currentTag)

		case xml.CharData:
			str := string([]byte(currentTag))
			fmt.Println("CharData ", str)
		}
	}
	return
}
