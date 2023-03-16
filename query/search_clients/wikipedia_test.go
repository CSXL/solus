package search_clients

import (
	"context"
	"testing"

	"net/http"
	"net/http/httptest"
)

func TestNewWikipediaClient(t *testing.T) {
	ctx := context.Background()
	_, err := NewWikipediaClient(ctx)
	if err != nil {
		t.Errorf("NewWikipediaClient() returned error: %v", err)
	}
}

func TestWikipediaClient_Search(t *testing.T) {
	ctx := context.Background()
	client, err := NewWikipediaClient(ctx)
	if err != nil {
		t.Errorf("NewWikipediaClient() returned error: %v", err)
	}
	query := "Computing"
	// Fake the response from the Wikipedia API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Response from https://en.wikipedia.org/w/api.php?action=query&list=search&format=json&srsearch=Computing&srlimit=2
		// Requested on 3/12/2023
		fake_response := `{"batchcomplete":"","continue":{"sroffset":2,"continue":"-||"},"query":{"searchinfo":{"totalhits":59315,"suggestion":"competing","suggestionsnippet":"competing"},"search":[{"ns":0,"title":"Computing","pageid":5213,"size":46982,"wordcount":4896,"snippet":"<span class=\"searchmatch\">Computing</span> is any goal-oriented activity requiring, benefiting from, or creating <span class=\"searchmatch\">computing</span> machinery. It includes the study and experimentation of algorithmic","timestamp":"2023-02-21T17:55:27Z"},{"ns":0,"title":"Cloud computing","pageid":19541494,"size":107970,"wordcount":10625,"snippet":"Cloud <span class=\"searchmatch\">computing</span> is the on-demand availability of computer system resources, especially data storage (cloud storage) and <span class=\"searchmatch\">computing</span> power, without direct","timestamp":"2023-03-11T19:06:24Z"}]}}`
		w.Write([]byte(fake_response))
	}))
	client.wikipediaActionBaseUrl = ts.URL
	results, err := client.Search(query)
	if err != nil {
		t.Errorf("Search() returned error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Search() returned wrong number of results: %d", len(results))
	}
	if results[0].Title != "Computing" || results[1].Title != "Cloud computing" {
		t.Errorf("Search() returned wrong results: %v", results)
	}
}

func TestWikipediaClient_GetPage(t *testing.T) {
	ctx := context.Background()
	client, err := NewWikipediaClient(ctx)
	if err != nil {
		t.Errorf("NewWikipediaClient() returned error: %v", err)
	}
	page := "2005_Azores_subtropical_storm"
	// Fake the response from the Wikipedia API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Truncated response from https://en.wikipedia.org/w/api.php?action=parse&prop=wikitext&page=2005_Azores_subtropical_storm&format=json
		// Requested on 3/15/2023
		fake_response := `{"parse":{"title":"2005 Azores subtropical storm","pageid":7715205,"wikitext":{"*":"Text truncated for testing purposes."}}}`
		w.Write([]byte(fake_response))
	}))
	client.wikipediaActionBaseUrl = ts.URL
	page_text, err := client.GetPage(page)
	if err != nil {
		t.Errorf("GetPage() returned error: %v", err)
	}
	if page_text != "Text truncated for testing purposes." {
		t.Errorf("GetPage() returned wrong text: %s", page_text)
	}
}
