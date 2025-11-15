package main

import "encoding/xml"

// UnmarshalXML parses Wayland protocol XML data into a Protocol structure
func UnmarshalXML(data []byte) (Protocol, error) {
	var r Protocol
	err := xml.Unmarshal(data, &r)
	return r, err
}

// MarshalXML serializes a Protocol structure into XML format
func (r *Protocol) MarshalXML() ([]byte, error) {
	return xml.Marshal(r)
}

// Protocol represents a Wayland protocol definition
type Protocol struct {
	Interface []Interface `xml:"interface"`
	Name      string      `xml:"name,attr"`
}

// Interface represents a Wayland protocol interface
type Interface struct {
	Description *Description `xml:"description,omitempty"`
	Request     []*Request   `xml:"request"`
	Event       []*Event     `xml:"event"`
	Enum        []*Enum      `xml:"enum"`
	Name        string       `xml:"name,attr"`
	Version     string       `xml:"version,attr"`
}

// Request represents a protocol request message
type Request struct {
	Description *Description `xml:"description,omitempty"`
	Name        string       `xml:"name,attr"`
	Arg         []*Arg       `xml:"arg"`
	Type        *string      `xml:"type,attr,omitempty"`
	Since       *string      `xml:"since,attr,omitempty"`
}

// Description contains documentation for protocol elements
type Description struct {
	Summary string `xml:"summary,attr"`
	Text    string `xml:",chardata"`
}

// Enum represents an enumeration type in the protocol
type Enum struct {
	Description *Description `xml:"description,omitempty"`
	Entry       []*Entry     `xml:"entry"`
	Name        string       `xml:"name,attr"`
	Bitfield    *string      `xml:"bitfield,attr,omitempty"`
}

// Event represents a protocol event message
type Event struct {
	Description *Description `xml:"description,omitempty"`
	Arg         []*Arg       `xml:"arg"`
	Name        string       `xml:"name,attr"`
	Since       *string      `xml:"since,attr,omitempty"`
	Type        *string      `xml:"type,attr,omitempty"`
}

// Entry represents an enumeration entry
type Entry struct {
	Name    string  `xml:"name,attr"`
	Value   string  `xml:"value,attr"`
	Summary *string `xml:"summary,attr,omitempty"`
	Since   *string `xml:"since,attr,omitempty"`
}

// Arg represents a request or event argument
type Arg struct {
	Name      string  `xml:"name,attr"`
	Type      string  `xml:"type,attr"`
	Summary   string  `xml:"summary,attr"`
	Interface *string `xml:"interface,attr,omitempty"`
	AllowNull *string `xml:"allow-null,attr,omitempty"`
	Enum      *string `xml:"enum,attr,omitempty"`
}
