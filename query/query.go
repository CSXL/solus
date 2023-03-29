package query

import (
	"context"

	"github.com/CSXL/solus/query/search_clients"
)

type QueryBuilder struct {
	ctx                context.Context
	searchClientConfig search_clients.SearchClientConfig
	googleSearchClient *search_clients.GoogleSearchClient
	queryText          string
	_type              string
	results            string
	err                error
}

// NewQuery returns a new QueryBuilder.
// It uses the builder design pattern to allow method chaining.
func NewQuery(ctx context.Context, searchClientConfig search_clients.SearchClientConfig) *QueryBuilder {
	googleSheetAPIKey := searchClientConfig.GetGoogleSearchAPIKey()
	googleSheetID := searchClientConfig.GetGoogleSearchEngineID()
	googleSearchClient, err := search_clients.NewGoogleSearchClient(ctx, googleSheetAPIKey, googleSheetID)
	return &QueryBuilder{ctx: ctx, searchClientConfig: searchClientConfig, googleSearchClient: googleSearchClient, err: err}
}

func (q *QueryBuilder) SetQueryText(text string) *QueryBuilder {
	q.queryText = text
	return q
}

func (q *QueryBuilder) SetType(t string) *QueryBuilder {
	q._type = t
	return q
}

func (q *QueryBuilder) GetQueryText() string {
	return q.queryText
}

func (q *QueryBuilder) GetType() string {
	return q._type
}

func (q *QueryBuilder) GetResults() (string, error) {
	return q.results, q.err
}

// Execute executes the query and stores the results in the QueryBuilder.
func (q *QueryBuilder) Execute() *QueryBuilder {
	if q.err != nil {
		return q
	}
	googleSearchResults, err := q.googleSearchClient.Search(q.queryText)
	if err != nil {
		q.err = err
		return q
	}
	googleSearchResultsString, err := search_clients.GoogleSearchResultsToJSON(googleSearchResults)
	if err != nil {
		q.err = err
	}
	q.results = googleSearchResultsString
	return q
}
