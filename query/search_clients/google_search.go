package search_clients

import (
	"context"

	"google.golang.org/api/cloudsearch/v1"
	"google.golang.org/api/option"
)

type GoogleSearchClient struct {
	API_KEY string
	ctx     context.Context
	client  *cloudsearch.Service
}

func NewGoogleSearchClient(ctx context.Context, api_key string) (*GoogleSearchClient, error) {
	client, err := cloudsearch.NewService(ctx, option.WithAPIKey(api_key))
	if err != nil {
		return nil, err
	}
	return &GoogleSearchClient{
		API_KEY: api_key,
		client:  client,
		ctx:     ctx,
	}, nil
}

func (g *GoogleSearchClient) buildRequest(query string) *cloudsearch.SearchRequest {
	return &cloudsearch.SearchRequest{
		Query: query,
	}
}

func (g *GoogleSearchClient) Search(query string) ([]*cloudsearch.SearchResult, error) {
	request := g.buildRequest(query)
	search_request := g.client.Query.Search(request)
	search_response, err := search_request.Do()
	if err != nil {
		return nil, err
	}
	return search_response.Results, nil
}
