package search_clients

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	wikipedia "github.com/CSXL/go-wikipedia"
)

var (
	// Unique User-Agent header value for this client, see
	// https://en.wikipedia.org/api/rest_v1/#/ for more information.
	UserAgent = "CSXL/Solus/1.0.0"
)

type WikipediaClient struct {
	ctx                    context.Context
	httpclient             *http.Client
	wikipediaActionBaseUrl string
	client                 *wikipedia.APIClient
}

func NewWikipediaClient(ctx context.Context) (*WikipediaClient, error) {
	cfg := wikipedia.NewConfiguration()
	cfg.UserAgent = UserAgent
	client := wikipedia.NewAPIClient(cfg)
	httpclient := &http.Client{}
	return &WikipediaClient{
		ctx:                    ctx,
		client:                 client,
		httpclient:             httpclient,
		wikipediaActionBaseUrl: "https://en.wikipedia.org/w/api.php",
	}, nil
}

func (c *WikipediaClient) buildSearchRequest(query string) (*http.Request, error) {
	requestUrl, err := url.Parse(c.wikipediaActionBaseUrl)
	if err != nil {
		return nil, err
	}
	url_query := url.Values{
		"action":   {"query"},
		"list":     {"search"},
		"srsearch": {query},
		"format":   {"json"},
	}
	requestUrl.RawQuery = url_query.Encode()
	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	return req, nil
}

type WikipediaSearchResult struct {
	Title  string `json:"title"`
	Url    string `json:"url"`
	PageID int    `json:"pageid"`
	Size   int    `json:"size"`
	// A snippet includes html higlighting matching text on the page from the search query.
	// Example:
	// <span class=\"searchmatch\">Computing</span> is any goal-oriented activity requiring, benefiting from...
	Snippet string `json:"snippet"`
}

type wikipediaQuery struct {
	Search []WikipediaSearchResult `json:"search"`
}

type wikipediaQueryResult struct {
	Query wikipediaQuery `json:"query"`
}

// Search performs a search on Wikipedia and returns a list of results.
// See https://en.wikipedia.org/w/api.php?action=help&modules=query%2Bsearch for more information.
func (c *WikipediaClient) Search(query string) ([]WikipediaSearchResult, error) {
	request, err := c.buildSearchRequest(query)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpclient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result wikipediaQueryResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	results := result.Query.Search
	return results, nil
}
