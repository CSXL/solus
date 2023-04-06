package openai

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	apiKey string
	ctx    context.Context
	client *openai.Client
}

func NewOpenAI(apiKey string) *OpenAI {
	return &OpenAI{
		apiKey: apiKey,
		ctx:    context.Background(),
		client: openai.NewClient(apiKey),
	}
}

func NewOpenAIWithBaseURL(apiKey string, baseURL string) *OpenAI {
	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = baseURL
	return &OpenAI{
		apiKey: apiKey,
		ctx:    context.Background(),
		client: openai.NewClientWithConfig(cfg),
	}
}

func (o *OpenAI) GetCompletion(prompt string, model string) (string, error) {
	resp, err := o.client.CreateCompletion(
		o.ctx,
		openai.CompletionRequest{
			Prompt: prompt,
			Model:  model,
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Text, nil
}
