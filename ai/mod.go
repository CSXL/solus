package ai

type AIConfig struct {
	OpenAIAPIKey string
}

func NewAIConfig(openAIAPIKey string) *AIConfig {
	return &AIConfig{
		OpenAIAPIKey: openAIAPIKey,
	}
}
