package owl

import "encoding/xml"

const (
	//XMLNS    = "http://inhost.com.ua/ontologies/"
	XML_BASE = "http://inhost.com.ua/ontologies/"
	//xmlns:medical_test="http://inhost.com.ua/ontologies/medical_test#"
	XMLNS_RDF  = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	XMLNS_OWL  = "http://www.w3.org/2002/07/owl#"
	XMLNS_XML  = "http://www.w3.org/XML/1998/namespace"
	XMLNS_XSD  = "http://www.w3.org/2001/XMLSchema#"
	XMLNS_RDFS = "http://www.w3.org/2000/01/rdf-schema#"
)

//type Rdf struct {
//	XMLName xml.Name `xml:"rdf:RDF"`
//	Xmlns   string   `xml:"xmlns,attr"`
//	Base    string   `xml:"xml:base,attr"`
//	Rdf     string   `xml:"xmlns:rdf,attr"`
//	Owl     string   `xml:"xmlns:owl,attr"`
//	Xml     string   `xml:"xmlns:xml,attr"`
//	Xsd     string   `xml:"xmlns:xsd,attr"`
//	Rdfs    string   `xml:"xmlns:rdfs,attr"`
//}

type OntologyTag struct {
	XMLName xml.Name `xml:"owl:Ontology"`
	About   string   `xml:"rdf:about,attr"`

	VersionIri VersionIriTag
	Label      LabelTag
}

type VersionIriTag struct {
	XMLName  xml.Name `xml:"owl:versionIRI"`
	Resource string   `xml:"rdf:resource,attr"`
}

type LabelTag struct {
	XMLName  xml.Name `xml:"rdfs:label"`
	Datatype string   `xml:"rdf:datatype,attr"`
	InnerXML string   `xml:",innerxml"`
}

type Root struct {
	XMLName xml.Name `xml:"rdf:RDF"`
	Xmlns   string   `xml:"xmlns,attr"`
	Base    string   `xml:"xml:base,attr"`
	Rdf     string   `xml:"xmlns:rdf,attr"`
	Owl     string   `xml:"xmlns:owl,attr"`
	Xml     string   `xml:"xmlns:xml,attr"`
	Xsd     string   `xml:"xmlns:xsd,attr"`
	Rdfs    string   `xml:"xmlns:rdfs,attr"`

	Ontology OntologyTag
}

func BuildHeader(name string) Root {
	ver := VersionIriTag{
		XMLName:  xml.Name{},
		Resource: XML_BASE + name,
	}
	l := LabelTag{
		XMLName:  xml.Name{},
		Datatype: XMLNS_XSD + "string",
		InnerXML: name,
	}
	o := OntologyTag{
		XMLName:    xml.Name{},
		About:      XML_BASE + name,
		VersionIri: ver,
		Label:      l,
	}
	h := Root{
		XMLName: xml.Name{},
		Xmlns:   XML_BASE + name + "#",
		Base:    XML_BASE + name + "#",
		Rdf:     XMLNS_RDF,
		Owl:     XMLNS_OWL,
		Xml:     XMLNS_XML,
		Xsd:     XMLNS_XSD,
		Rdfs:    XMLNS_RDFS,

		Ontology: o,
	}
	return h
}

type ClassTag struct {
	XMLName  xml.Name        `xml:"owl:Class"`
	About    string          `xml:"rdf:about,attr"`
	SubClass []SubClassOfTag `xml:",omitempty"`
	Label    LabelTag        `xml:",omitempty"`
	Comment  CommentTag      `xml:",omitempty"`
}

type SubClassOfTag struct {
	XMLName  xml.Name `xml:"rdfs:subClassOf"`
	Resource string   `xml:"rdf:resource,attr"`
}

type CommentTag struct {
	XMLName  xml.Name `xml:"rdfs:comment,omitempty"`
	Lang     string   `xml:"xml:lang,attr,omitempty"`
	CharData string   `xml:",chardata"`
}

type EquivalentClassTag struct {
	XMLName xml.Name `xml:"owl:equivalentClass"`
}
type UnionOfTag struct {
	XMLName     xml.Name `xml:"owl:unionOf"`
	ParseType   string   `xml:"rdf:parseType,attr"`
	Description []DescriptionTag
}
type DescriptionTag struct {
	XMLName xml.Name `xml:"rdf:Description"`
	About   string   `xml:"rdf:about,attr"`
}
type NamedIndividualTag struct {
	XMLName xml.Name `xml:"owl:NamedIndividual"`
	About   string   `xml:"rdf:about,attr"`
}
