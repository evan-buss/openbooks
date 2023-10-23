package opds

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// Feed link types
	FeedAcquisitionLinkType = "application/atom+xml;profile=opds-catalog;kind=acquisition"
	FeedNavigationLinkType  = "application/atom+xml;profile=opds-catalog;kind=navigation"
	FeedSearchLinkType      = "application/opensearchdescription+xml"
	// Feed link relations
	FeedStartLinkRel      = "start"
	FeedSelfLinkRel       = "self"
	FeedSearchLinkRel     = "search"
	FeedFirstLinkRel      = "first"
	FeedLastLinkRel       = "last"
	FeedNextLinkRel       = "next"
	FeedPrevLinkRel       = "prev"
	FeedSubsectionLinkRel = "subsection"

	// Content types
	FeedTextContentType = "text"
	FeedHtmlContentType = "html"
)

type Feed struct {
	XMLName      xml.Name `xml:"feed"`
	Xmlns        string   `xml:"xmlns,attr"`
	XmlnsDC      string   `xml:"xmlns:dcterms,attr,omitempty"`
	XmlnsOS      string   `xml:"xmlns:opensearch,attr,omitempty"`
	XmlnsOPDS    string   `xml:"xmlns:opds,attr,omitempty"`
	Title        string   `xml:"title"`
	ID           string   `xml:"id"`
	Updated      TimeStr  `xml:"updated"`
	Link         []Link   `xml:"link"`
	Author       []Author `xml:"author,omitempty"`
	Entry        []*Entry `xml:"entry"`
	Category     string   `xml:"category,omitempty"`
	Icon         string   `xml:"icon,omitempty"`
	Logo         string   `xml:"logo,omitempty"`
	Content      string   `xml:"content,omitempty"`
	Subtitle     string   `xml:"subtitle,omitempty"`
	SearchResult uint     `xml:"opensearch:totalResults,omitempty"`
}

type Entry struct {
	// XMLName   xml.Name `xml:"entry"`
	// Xmlns     string   `xml:"xmlns,attr,omitempty"`
	Title     string   `xml:"title"`
	ID        string   `xml:"id"`
	Link      []Link   `xml:"link"`
	Published string   `xml:"published,omitempty"`
	Updated   TimeStr  `xml:"updated"`
	Category  string   `xml:"category,omitempty"`
	Authors   []Author `xml:"author"`
	Summary   *Summary `xml:"summary"`
	Content   *Content `xml:"content"`
	Rights    string   `xml:"rights,omitempty"`
	Source    string   `xml:"source,omitempty"`
}

type Link struct {
	// XMLName xml.Name `xml:"link"`
	Type   string `xml:"type,attr,omitempty"`
	Title  string `xml:"title,attr,omitempty"`
	Href   string `xml:"href,attr"`
	Rel    string `xml:"rel,attr,omitempty"`
	Length string `xml:"length,attr,omitempty"`
}

type Author struct {
	// XMLName xml.Name `xml:"author"`
	Name string `xml:"name,omitempty"`
	Uri  string `xml:"uri,omitempty"`
}

type Summary struct {
	// XMLName xml.Name `xml:"summary"`
	Content string `xml:",chardata"`
	Type    string `xml:"type,attr"`
}

type Content struct {
	// XMLName xml.Name `xml:"content"`
	Content string `xml:",chardata"`
	Type    string `xml:"type,attr"`
}

type TimeStr string

func (f *Feed) Time(t time.Time) TimeStr {
	return TimeStr(t.Format("2006-01-02T15:04:05-07:00"))
}

func NewFeed(title, subtitle, self string) *Feed {
	f := &Feed{
		XMLName:   xml.Name{},
		Xmlns:     "http://www.w3.org/2005/Atom",
		XmlnsDC:   "http://purl.org/dc/terms/",
		XmlnsOS:   "http://a9.com/-/spec/opensearch/1.1/",
		XmlnsOPDS: "http://opds-spec.org/2010/catalog",
		Title:     title,
		ID:        self,
		Link: []Link{
			{Rel: FeedSearchLinkRel, Href: "/opds/opensearch", Type: FeedSearchLinkType, Title: "Search on catalog"},
			{Rel: FeedStartLinkRel, Href: "/opds", Type: FeedNavigationLinkType},
			{Rel: FeedSelfLinkRel, Href: self, Type: FeedNavigationLinkType},
		},
		Subtitle: subtitle,
	}
	f.Updated = f.Time(time.Now())
	return f
}

func writeFeed(w http.ResponseWriter, statusCode int, f Feed) {
	data, err := xml.MarshalIndent(f, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Internal server error")
		return
	}
	s := fmt.Sprintf("%s%s", xml.Header, data)
	w.Header().Add("Content-Type", "application/atom+xml")
	w.WriteHeader(statusCode)
	io.WriteString(w, s)
}
