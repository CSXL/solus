package query

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CSXL/solus/query/search_clients"
	"google.golang.org/api/customsearch/v1"
)

func TestNewQuery(t *testing.T) {
	ctx := context.Background()
	searchClientConfig := search_clients.NewSearchClientConfig("test", "test")
	q := NewQuery(ctx, *searchClientConfig)
	if q == nil {
		t.Errorf("NewQuery() returned nil")
	}
}

// TODO: Split this into multiple tests.
func TestQueryBuilder_SettersandGetters(t *testing.T) {
	ctx := context.Background()
	searchClientConfig := search_clients.NewSearchClientConfig("test", "test")
	q := NewQuery(ctx, *searchClientConfig)
	queryText := "test_query"
	q.SetQueryText(queryText)
	if q.GetQueryText() != queryText {
		t.Errorf("SetQueryText() or GetQueryText() failed")
	}
	_type := "test_type"
	q.SetType(_type)
	if q.GetType() != _type {
		t.Errorf("SetType() or GetType() failed")
	}
}

func TestQueryBuilder_Execute(t *testing.T) {
	ctx := context.Background()
	googleSearchAPIKey := "test"
	googleSearchEngineID := "test"
	searchClientConfig := search_clients.NewSearchClientConfig(googleSearchAPIKey, googleSearchEngineID)
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
		// We know the response is valid so we don't need error checking.
		// trunk-ignore(golangci-lint/errcheck)
		json.NewEncoder(w).Encode(response)
	}))
	q := NewQuery(ctx, *searchClientConfig)
	q.googleSearchClient.SetBasePath(ts.URL)
	queryText := "test_query"
	q.SetQueryText(queryText)
	q.Execute()

	results, err := q.GetResults()
	if err != nil {
		t.Errorf("Execute() returned error: %v", err)
	}
	if results == "" {
		t.Errorf("Execute() returned empty results")
	}
}
