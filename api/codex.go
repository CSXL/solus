package api

import openai "github.com/sashabaranov/go-openai"

var CODEX_MODEL = openai.CodexCodeDavinci002

type CodexClient struct {
	apiKey       string
	openAIClient *OpenAI
}

func NewCodexClient(apiKey string) *CodexClient {
	return &CodexClient{
		apiKey:       apiKey,
		openAIClient: NewOpenAI(apiKey),
	}
}

func (c *CodexClient) ExecuteCodexCompletion(prompt string) (string, error) {
	response, err := c.openAIClient.GetCompletion(prompt, CODEX_MODEL)
	if err != nil {
		return "", err
	}
	return response, nil
}
