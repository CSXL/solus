package search_clients

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
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
	UserAgent              string
}

func NewWikipediaClient(ctx context.Context) (*WikipediaClient, error) {
	httpclient := &http.Client{}
	return &WikipediaClient{
		ctx:                    ctx,
		UserAgent:              UserAgent,
		httpclient:             httpclient,
		wikipediaActionBaseUrl: "https://en.wikipedia.org/w/api.php",
	}, nil
}

func (c *WikipediaClient) doRequest(query url.Values) (*http.Response, error) {
	requestUrl, err := url.Parse(c.wikipediaActionBaseUrl)
	if err != nil {
		return nil, err
	}
	requestUrl.RawQuery = query.Encode()
	req, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	return c.httpclient.Do(req)
}

func (c *WikipediaClient) doSearchRequest(query string) (*http.Response, error) {
	url_query := url.Values{
		"action":   {"query"},
		"list":     {"search"},
		"srsearch": {query},
		"format":   {"json"},
	}
	response, err := c.doRequest(url_query)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type WikipediaQuerySearchResult struct {
	Title  string `json:"title"`
	Url    string `json:"url"`
	PageID int    `json:"pageid"`
	Size   int    `json:"size"`
	// A snippet includes html higlighting matching text on the page from the search query.
	// Example:
	// <span class=\"searchmatch\">Computing</span> is any goal-oriented activity requiring, benefiting from...
	Snippet string `json:"snippet"`
}

type wikipediaQuerySearch struct {
	Search []WikipediaQuerySearchResult `json:"search"`
}

type wikipediaQuery struct {
	Query wikipediaQuerySearch `json:"query"`
}

// Search performs a search on Wikipedia and returns a list of results.
// See https://en.wikipedia.org/w/api.php?action=help&modules=query%2Bsearch for more information.
func (c *WikipediaClient) Search(query string) ([]WikipediaQuerySearchResult, error) {
	response, err := c.doSearchRequest(query)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var result wikipediaQuery
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	results := result.Query.Search
	return results, nil
}
