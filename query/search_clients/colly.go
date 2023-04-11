package search_clients

import (
	"github.com/PuerkitoBio/goquery"
	colly "github.com/gocolly/colly/v2"
)

type HTMLElement = colly.HTMLElement
type Request = colly.Request
type HTMLCallback = colly.HTMLCallback
type RequestCallback = colly.RequestCallback
type ResponseCallback = colly.ResponseCallback
type ErrorCallback = colly.ErrorCallback

type Website struct {
	Title       string
	URL         string
	Links       []string
	MIME        string
	HTMLContent string
	TextContent string
}

func (w *Website) GetTitle() string {
	return w.Title
}

func (w *Website) GetURL() string {
	return w.URL
}

func (w *Website) GetLinks() []string {
	return w.Links
}

func (w *Website) GetMIME() string {
	return w.MIME
}

func (w *Website) GetHTMLContent() string {
	return w.HTMLContent
}

func (w *Website) GetTextContent() string {
	return w.TextContent
}

// Scraper wraps the Colly scraper
type Scraper struct {
	c *colly.Collector
}

// NewScraper creates a new Scraper
func NewScraper() *Scraper {
	return &Scraper{
		c: colly.NewCollector(),
	}
}

func (s *Scraper) ScrapePage(url string) (*Website, error) {
	var err error
	var w Website
	w.URL = url
	s.onHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		w.Links = append(w.Links, link)
	})
	s.onHTML("html", func(e *colly.HTMLElement) {
		w.Title = e.DOM.Find("title").Text()
		w.HTMLContent, err = e.DOM.Html()
		w.TextContent = s.getTextContent(*e)
	})
	s.onResponse(func(r *colly.Response) {
		w.MIME = r.Headers.Get("Content-Type")
	})
	s.onError(func(r *colly.Response, page_err error) {
		err = page_err
	})
	if err != nil {
		return nil, err
	}
	err = s.visit(url)
	if err != nil {
		return nil, err
	}
	return &w, err
}

// Gets text content from an HTML element
func (s *Scraper) getTextContent(e HTMLElement) string {
	doc := goquery.NewDocumentFromNode(e.DOM.Nodes[0])
	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		// Remove any non-text nodes
		if s.Is("script, style, head, iframe, input, textarea") {
			s.Remove()
		}
	})
	return doc.Text()
}

// Scrape scrapes the given URL
func (s *Scraper) visit(url string) error {
	return s.c.Visit(url)
}

// OnHTML registers a callback function for the OnHTML event
func (s *Scraper) onHTML(selector string, f HTMLCallback) {
	s.c.OnHTML(selector, f)
}

// OnError registers a callback function for the OnError event
func (s *Scraper) onError(f ErrorCallback) {
	s.c.OnError(f)
}

// OnResponse registers a callback function for the OnResponse event
func (s *Scraper) onResponse(f ResponseCallback) {
	s.c.OnResponse(f)
}
