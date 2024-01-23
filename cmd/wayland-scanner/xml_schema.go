package main

import "encoding/xml"

func UnmarshalXml(data []byte) (Xml, error) {
	var r Xml
	err := xml.Unmarshal(data, &r)
	return r, err
}

func (r *Xml) MarshalXml() ([]byte, error) {
	return xml.Marshal(r)
}

type Xml Protocol

type Protocol struct {
	Interface []Interface `xml:"interface"`
	Name      string      `xml:"name,attr"`
}

type Interface struct {
	Description Description `xml:"description"`
	Request     []*Request  `xml:"request"`
	Event       []*Event    `xml:"event"`
	Enum        []*Enum     `xml:"enum"`
	Name        string      `xml:"name,attr"`
	Version     string      `xml:"version,attr"`
}

type Request struct {
	Description Description `xml:"description"`
	Name        string      `xml:"name,attr"`
	Arg         []*Arg      `xml:"arg"`
	Type        *EventType  `xml:"type,attr,omitempty"`
	Since       *string     `xml:"since,attr,omitempty"`
}

type Description struct {
	Summary string  `xml:"summary,attr"`
	Text    *string `xml:"_text,attr,omitempty"`
}

type Enum struct {
	Description *Description `xml:"description,omitempty"`
	Entry       []*Entry     `xml:"entry"`
	Name        string       `xml:"name,attr"`
	Bitfield    *string      `xml:"bitfield,attr,omitempty"`
}

type Event struct {
	Description Description `xml:"description"`
	Arg         []*Arg      `xml:"arg"`
	Name        string      `xml:"name,attr"`
	Since       *string     `xml:"since,attr,omitempty"`
	Type        *EventType  `xml:"type,attr,omitempty"`
}

type EventType string

type Entry struct {
	Name    string  `xml:"name,attr"`
	Value   string  `xml:"value,attr"`
	Summary *string `xml:"summary,attr,omitempty"`
	Since   *string `xml:"since,attr,omitempty"`
}

type Arg struct {
	Name      string  `xml:"name,attr"`
	Type      string  `xml:"type,attr"`
	Summary   string  `xml:"summary,attr"`
	Interface *string `xml:"interface,attr,omitempty"`
	AllowNull *string `xml:"allow-null,attr,omitempty"`
	Enum      *string `xml:"enum,attr,omitempty"`
}
