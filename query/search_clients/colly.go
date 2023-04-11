package search_clients

import (
	colly "github.com/gocolly/colly/v2"
)

type HTMLElement = colly.HTMLElement
type Request = colly.Request
type HTMLCallback = colly.HTMLCallback
type RequestCallback = colly.RequestCallback
type ResponseCallback = colly.ResponseCallback
type ErrorCallback = colly.ErrorCallback

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

// Scrape scrapes the given URL
func (s *Scraper) Scrape(url string) error {
	return s.c.Visit(url)
}

// OnHTML registers a callback function for the OnHTML event
func (s *Scraper) OnHTML(selector string, f HTMLCallback) {
	s.c.OnHTML(selector, f)
}

// OnRequest registers a callback function for the OnRequest event
func (s *Scraper) OnRequest(f RequestCallback) {
	s.c.OnRequest(f)
}

// OnError registers a callback function for the OnError event
func (s *Scraper) OnError(f ErrorCallback) {
	s.c.OnError(f)
}

// OnResponse registers a callback function for the OnResponse event
func (s *Scraper) OnResponse(f ResponseCallback) {
	s.c.OnResponse(f)
}
