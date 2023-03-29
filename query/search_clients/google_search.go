package search_clients

import (
	"context"
	"encoding/json"

	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

type GoogleSearchClient struct {
	apiKey               string
	googleSearchEngineID string
	ctx                  context.Context
	client               *customsearch.Service
}

func NewGoogleSearchClient(ctx context.Context, apiKey string, googleSearchEngineID string) (*GoogleSearchClient, error) {
	client, err := customsearch.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GoogleSearchClient{
		apiKey:               apiKey,
		googleSearchEngineID: googleSearchEngineID,
		client:               client,
		ctx:                  ctx,
	}, nil
}

type GoogleSearchResult struct {
	Title    string
	Url      string
	Summary  string
	MIMEType string
}

// GetBasePath returns the base path of the GoogleSearchClient's HTTP client
func (gsc *GoogleSearchClient) GetBasePath() string {
	return gsc.client.BasePath
}

// SetBasePath sets the base path of the GoogleSearchClient's HTTP client
//
// This is useful for unit testing.
func (gsc *GoogleSearchClient) SetBasePath(basePath string) {
	gsc.client.BasePath = basePath
}

func (gsc *GoogleSearchClient) Search(query string) ([]*GoogleSearchResult, error) {
	response, err := gsc.client.Cse.List().Q(query).Cx(gsc.googleSearchEngineID).Do()
	if err != nil {
		return nil, err
	}
	results := make([]*GoogleSearchResult, len(response.Items))
	for i, item := range response.Items {
		results[i] = &GoogleSearchResult{
			Title:    item.Title,
			Url:      item.Link,
			Summary:  item.Snippet,
			MIMEType: item.Mime,
		}
	}
	return results, nil
}

func GoogleSearchResultsToJSON(results []*GoogleSearchResult) (string, error) {
	json, err := json.Marshal(results)
	if err != nil {
		return "", err
	}
	return string(json), nil
}
