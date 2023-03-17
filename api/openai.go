package api

import (
	openai "github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	apiKey string
	client *openai.Client
}

func NewOpenAI(apiKey string) *OpenAI {
	return &OpenAI{
		apiKey: apiKey,
		client: openai.NewClient(apiKey),
	}
}
