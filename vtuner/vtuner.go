package vtuner

import (
	"encoding/xml"
)

const (
	_      = `<EncryptedToken>0000000000000000</EncryptedToken>`
	header = `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>`
)

// Page is a struct that represents a page in the vTuner API.
type Page struct {
	XMLName     xml.Name `xml:"ListOfItems"`
	Items       []Item   `xml:">Item"`
	Count       int
	NoDataCache bool `xml:"NoDataCache"`
}

func NewPage(items []Item, cache bool) *Page {
	return &Page{
		Items:       items,
		NoDataCache: !cache,
	}
}

func (p *Page) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	p.Count = len(p.Items)
	return e.EncodeElement(*p, start)
}

func (p *Page) Render() ([]byte, error) {
	b, err := xml.MarshalIndent(p, "", " ")
	if err != nil {
		return nil, err
	}

	b = append([]byte(header), b...)
	return b, nil
}

// Item is an interface that represents a generic item in the vTuner API.
type Item interface {
	Type() string
}

// item is a concrete type that implements the Item interface.
type item struct {
	XMLName  xml.Name `xml:"Item"`
	ItemType string   `xml:"ItemType"`
}

// Type returns the type of the item.
func (i *item) Type() string {
	return i.ItemType
}

// Display is a struct that represents a display item in the vTuner API.
type Display struct {
	item
	Display string `xml:"Display"`
}

func (d *Display) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	d.ItemType = "Display"
	return e.EncodeElement(*d, start)
}

// Previous is a struct that represents a previous item in the vTuner API.
type Previous struct {
	item
	UrlPrevious       string `xml:"UrlPrevious"`
	UrlPreviousBackup string `xml:"UrlPreviousBackup"` // defaults to UrlPrevious
}

func (p *Previous) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	p.ItemType = "Previous"

	p.UrlPrevious = addBogusParam(p.UrlPrevious)
	p.UrlPreviousBackup = p.UrlPrevious
	return e.EncodeElement(*p, start)
}

// Search is a struct that represents a search item in the vTuner API.
type Search struct {
	item
	SearchURL          string `xml:"SearchURL"`
	SearchURLBackup    string `xml:"SearchURLBackup"`
	SearchCaption      string `xml:"SearchCaption"`
	SearchTextbox      string `xml:"SearchTextbox"`
	SearchButtonGo     string `xml:"SearchButtonGo"`
	SearchButtonCancel string `xml:"SearchButtonCancel"`
}

func (s *Search) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	s.ItemType = "Search"

	if s.SearchButtonCancel == "" {
		s.SearchButtonCancel = "Cancel"
	}

	if s.SearchButtonGo == "" {
		s.SearchButtonGo = "Search"
	}

	s.SearchURL = addBogusParam(s.SearchURL)
	s.SearchURLBackup = s.SearchURL

	return e.EncodeElement(*s, start)
}

type Directory struct {
	item
	Title        string `xml:"Title"`
	UrlDir       string `xml:"UrlDir"`
	UrlDirBackup string `xml:"UrlDirBackup"`
	DirCount     int    `xml:"DirCount"`
}

func (d *Directory) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	d.ItemType = "Directory"

	d.UrlDir = addBogusParam(d.UrlDir)
	d.UrlDirBackup = d.UrlDir

	return e.EncodeElement(*d, start)
}

func addBogusParam(url string) string {
	// We need this bogus parameter because some (if not all) AVRs blindly append additional request parameters
	// with an ampersand. E.g.: '&mac=<REDACTED>&dlang=eng&fver=1.2&startitems=1&enditems=100'.
	// The original vTuner API hacks around that by adding a specific parameter or a bogus parameter like '?empty=' to
	// the target URL.
	return url + "?vtuner=true"
}
