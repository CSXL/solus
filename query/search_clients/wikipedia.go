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

func (c *WikipediaClient) doParseRequest(pageTitle string) (*http.Response, error) {
	url_query := url.Values{
		"action": {"parse"},
		"page":   {pageTitle},
		"format": {"json"},
	}
	response, err := c.doRequest(url_query)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type WikipediaParseResult struct {
	Data string `json:"*"`
}

type wikipediaParseWikiText struct {
	WikiText WikipediaParseResult `json:"wikitext"`
}

type wikipediaParse struct {
	Parse wikipediaParseWikiText `json:"parse"`
}

func (c *WikipediaClient) GetPage(pageTitle string) (string, error) {
	response, err := c.doParseRequest(pageTitle)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	var result wikipediaParse
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	data := result.Parse.WikiText.Data
	return data, nil
}

func (c *WikipediaClient) doPageSummaryRequest(pageTitle string) (*http.Response, error) {
	url_query := url.Values{
		"action":      {"query"},
		"prop":        {"extracts"},
		"exintro":     {"true"},
		"explaintext": {"true"},
		"titles":      {pageTitle},
		"format":      {"json"},
	}
	response, err := c.doRequest(url_query)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type WikipediaQuerySummary struct {
	Title   string `json:"title"`
	Summary string `json:"extract"`
}

type wikipediaQuerySummaryPages struct {
	Pages map[string]WikipediaQuerySummary `json:"pages"`
}

type wikipediaQuerySummary struct {
	Query wikipediaQuerySummaryPages `json:"query"`
}

// GetPageSummary returns the summary of a Wikipedia page.
// See https://en.wikipedia.org/w/api.php?action=help&modules=query%2Bextracts for more information.
func (c *WikipediaClient) GetPageSummary(pageTitle string) (string, error) {
	response, err := c.doPageSummaryRequest(pageTitle)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	var result wikipediaQuerySummary
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	// The response is a map of pages, but we only requested one page so we can just take the first one.
	var summary string
	for _, page := range result.Query.Pages {
		summary = page.Summary
		break
	}
	return summary, nil
}
