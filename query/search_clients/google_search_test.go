package search_clients

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/api/cloudsearch/v1"
)

func TestNewGoogleSearchClient(t *testing.T) {
	ctx := context.Background()
	API_KEY := "test"
	_, err := NewGoogleSearchClient(ctx, API_KEY)
	if err != nil {
		t.Errorf("NewGoogleSearchClient() returned error: %v", err)
	}
}

func TestGoogleSearchClient_Search(t *testing.T) {
	ctx := context.Background()
	API_KEY := "test"
	client, err := NewGoogleSearchClient(ctx, API_KEY)
	if err != nil {
		t.Errorf("NewGoogleSearchClient() returned error: %v", err)
	}
	query := "test_query"
	// Fake the response from the Google Search API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := &cloudsearch.SearchResponse{
			Results: []*cloudsearch.SearchResult{
				{
					Title: "test_title",
					Url:   "test_url",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))

	// Set the client's HTTP client to the test server
	client.client.BasePath = ts.URL

	// Perform the search
	results, err := client.Search(query)
	if err != nil {
		t.Errorf("Search() returned error: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Search() returned wrong number of results: %d", len(results))
	}
	if results[0].Title != "test_title" || results[0].Url != "test_url" {
		t.Errorf("Search() returned wrong results: %v", results)
	}
}
