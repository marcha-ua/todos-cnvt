package parser

import (
	"encoding/xml"
	"fmt"
	"io"
)

type TodosNode struct{

}

func OntologyOwlFromParser(r io.Reader) ( err error) {
	decoder := xml.NewDecoder(r)
	//	total := 0
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		//fmt.Println(token)
		switch startElement := token.(type) {
		case xml.StartElement:
			if startElement.Name.Local == "entry" {
				// do what you need to do for each entry below
			}
			fmt.Println("startElement ", startElement)
		case xml.EndElement:
			fmt.Println("EndElement ", startElement)
		case xml.Attr:
			fmt.Println("Attr ", startElement)
		case xml.Comment:
			fmt.Println("Comment ", startElement)
		case xml.Name:
			fmt.Println("Name ", startElement)

		case xml.CharData:
			str := string([]byte(startElement))
			fmt.Println("CharData ", str)
		}
	}
	return
}