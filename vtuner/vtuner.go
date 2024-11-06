package vtuner

import (
	"encoding/xml"
	"io"
)

const (
	header = `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>`
)

// BogusToken is a bogus token used in the vTuner API.
var EncryptedToken = []byte(`<EncryptedToken>0000000000000000</EncryptedToken>`)

// Page is a struct that represents a page in the vTuner API.
type Page struct {
	XMLName     xml.Name `xml:"ListOfItems"`
	Items       []Item   `xml:"Item"`
	Count       int      `xml:"ItemsCount"`
	NoDataCache string   `xml:"NoDataCache"`
}

func NewPage(items []Item, cache bool) *Page {
	result := &Page{
		Items: items,
	}

	if cache {
		result.NoDataCache = "no"
	} else {
		result.NoDataCache = "yes"
	}

	return result
}

func (p *Page) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "ListOfItems" // Override the XML name

	p.Count = len(p.Items)

	return e.EncodeElement(*p, start)
}

func (p *Page) Write(w io.Writer) error {
	b, err := xml.MarshalIndent(p, "", " ")
	if err != nil {
		return err
	}

	b = append([]byte(header), b...)

	_, err = w.Write(b)
	return err
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
	// URL is used to generate the URLPrevious and URLPreviousBackup fields.
	// Set it to the URL of the previous item.
	Url string `xml:"-"`

	// Used for XML marshalling.
	UrlPrevious       string `xml:"UrlPrevious"`
	UrlPreviousBackup string `xml:"UrlPreviousBackup"` // defaults to UrlPrevious
}

func (p *Previous) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	p.ItemType = "Previous"

	p.UrlPrevious = addBogusParam(p.Url)
	p.UrlPreviousBackup = addBogusParam(p.Url)
	return e.EncodeElement(*p, start)
}

// Search is a struct that represents a search item in the vTuner API.
type Search struct {
	item
	Caption string `xml:"-"`
	URL     string `xml:"-"`

	// Used for XML marshalling.
	SearchURL          string `xml:"SearchURL"`
	SearchURLBackup    string `xml:"SearchURLBackup"`
	SearchCaption      string `xml:"SearchCaption"`
	SearchTextbox      string `xml:"SearchTextbox"`
	SearchButtonGo     string `xml:"SearchButtonGo"`
	SearchButtonCancel string `xml:"SearchButtonCancel"`
}

func (s *Search) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	s.ItemType = "Search"

	s.SearchButtonCancel = "Cancel"
	s.SearchButtonGo = "Search"

	s.SearchURL = addBogusParam(s.URL)
	s.SearchURLBackup = addBogusParam(s.URL)
	s.SearchCaption = s.Caption
	s.SearchTextbox = "" // TODO: Implement search textbox

	return e.EncodeElement(*s, start)
}

type Directory struct {
	item
	Title          string `xml:"Title"`
	DestinationURL string `xml:"-"`
	Count          int    `xml:"-"`

	// Used for XML marshalling.
	UrlDir       string `xml:"UrlDir"`
	UrlDirBackup string `xml:"UrlDirBackup"`
	DirCount     int    `xml:"DirCount"`
}

func (d *Directory) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	d.ItemType = "Directory"

	d.UrlDir = addBogusParam(d.DestinationURL)
	d.UrlDirBackup = addBogusParam(d.DestinationURL)
	d.DirCount = d.Count

	return e.EncodeElement(*d, start)
}

// def __init__(self, uid, name, description, url, icon, genre, location, mime, bitrate, bookmark)
type Station struct {
	item
	ID          string `xml:"StationId"`
	Name        string `xml:"StationName"`
	URL         string `xml:"StationUrl"`
	Description string `xml:"StationDescription"`
	Logo        string `xml:"StationLogo"`
	Format      string `xml:"StationFormat"`
	Bitrate     int    `xml:"StationBandwidth"`
	MIME        string `xml:"StationMime"`
	Relia       int    `xml:"Relia"` // is always 3 I think
	Bookmark    int    `xml:"Bookmark"`
}

func (s *Station) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// TODO: add trackURL
	s.ItemType = "Station"
	s.Relia = 3
	return e.EncodeElement(*s, start)
}

func addBogusParam(url string) string {
	// We need this bogus parameter because some (if not all) AVRs blindly append additional request parameters
	// with an ampersand. E.g.: '&mac=<REDACTED>&dlang=eng&fver=1.2&startitems=1&enditems=100'.
	// The original vTuner API hacks around that by adding a specific parameter or a bogus parameter like '?empty=' to
	// the target URL.
	return url + "?vtuner=true"
}
