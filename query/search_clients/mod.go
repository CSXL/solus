package search_clients

type SearchClientConfig struct {
	GoogleSearchAPIKey   string
	GoogleSearchEngineID string
}

func NewSearchClientConfig(googleSearchAPIKey string, googleSearchEngineID string) *SearchClientConfig {
	return &SearchClientConfig{
		GoogleSearchAPIKey:   googleSearchAPIKey,
		GoogleSearchEngineID: googleSearchEngineID,
	}
}

func (sc *SearchClientConfig) GetGoogleSearchAPIKey() string {
	return sc.GoogleSearchAPIKey
}

func (sc *SearchClientConfig) GetGoogleSearchEngineID() string {
	return sc.GoogleSearchEngineID
}
