package search_clients

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/api/customsearch/v1"
)

func TestNewGoogleSearchClient(t *testing.T) {
	ctx := context.Background()
	api_key := "test"
	programmable_search_id := "test"
	_, err := NewGoogleSearchClient(ctx, api_key, programmable_search_id)
	if err != nil {
		t.Errorf("NewGoogleSearchClient() returned error: %v", err)
	}
}

func TestGoogleSearchClient_Search(t *testing.T) {
	ctx := context.Background()
	api_key := "test"
	programmable_search_id := "test"
	client, err := NewGoogleSearchClient(ctx, api_key, programmable_search_id)
	if err != nil {
		t.Errorf("NewGoogleSearchClient() returned error: %v", err)
	}
	query := "test_query"
	// Fake the response from the Google Programmble Search Engine API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := &customsearch.Search{
			Items: []*customsearch.Result{
				{
					Title:   "test_title",
					Link:    "test_url",
					Mime:    "text/plain",
					Snippet: "test_summary",
				},
				{
					Title:   "test_title2",
					Link:    "test_url2",
					Mime:    "text/plain",
					Snippet: "test_summary2",
				},
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))

	// Set the client's HTTP client to the test server
	client.SetBasePath(ts.URL)

	// Perform the search
	results, err := client.Search(query)
	if err != nil {
		t.Errorf("Search() returned error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("Search() returned wrong number of results: %d", len(results))
	}
	if results[0].Title != "test_title" || results[0].Url != "test_url" || results[0].Summary != "test_summary" || results[0].MIMEType != "text/plain" {
		t.Errorf("Search() returned wrong results: %v", results)
	}
}

func TestGoogleSearchResultsToJSON(t *testing.T) {
	results := []*GoogleSearchResult{
		{
			Title:    "test_title",
			Url:      "test_url",
			Summary:  "test_summary",
			MIMEType: "text/plain",
		},
		{
			Title:    "test_title2",
			Url:      "test_url2",
			Summary:  "test_summary2",
			MIMEType: "text/plain",
		},
	}
	json, err := GoogleSearchResultsToJSON(results)
	if err != nil {
		t.Errorf("GoogleSearchResultsToJSON() returned error: %v", err)
	}
	expected := `[{"Title":"test_title","Url":"test_url","Summary":"test_summary","MIMEType":"text/plain"},{"Title":"test_title2","Url":"test_url2","Summary":"test_summary2","MIMEType":"text/plain"}]`
	if json != expected {
		t.Errorf("GoogleSearchResultsToJSON() returned wrong JSON: %s", json)
	}
}
