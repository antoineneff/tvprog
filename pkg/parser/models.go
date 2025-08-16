package parser

import "encoding/xml"

type TV struct {
	XMLName   xml.Name    `xml:"tv"`
	Channels  []Channel   `xml:"channel"`
	Programmes []Programme `xml:"programme"`
}

type Channel struct {
	ID          string `xml:"id,attr"`
	DisplayName string `xml:"display-name"`
}

type Programme struct {
	Title   string `xml:"title"`
	Channel string `xml:"channel,attr"`
	Start   string `xml:"start,attr"`
	Stop    string `xml:"stop,attr"`
}