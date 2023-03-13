package search_clients

import (
	"context"

	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"
)

type GoogleSearchClient struct {
	api_key                string
	programmable_search_id string
	ctx                    context.Context
	client                 *customsearch.Service
}

func NewGoogleSearchClient(ctx context.Context, api_key string, programmable_search_id string) (*GoogleSearchClient, error) {
	client, err := customsearch.NewService(ctx, option.WithAPIKey(api_key))
	if err != nil {
		return nil, err
	}
	return &GoogleSearchClient{
		api_key:                api_key,
		programmable_search_id: programmable_search_id,
		client:                 client,
		ctx:                    ctx,
	}, nil
}

type GoogleSearchResult struct {
	Title    string
	Url      string
	Summary  string
	MIMEType string
}

func (gsc *GoogleSearchClient) Search(query string) ([]*GoogleSearchResult, error) {
	response, err := gsc.client.Cse.List().Q(query).Cx(gsc.programmable_search_id).Do()
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
