// Endpoints for mocking OpenAI API calls.
package openai

type OpenAIResponse string

const (
	SampleChatCompletion OpenAIResponse = `{"id":"chatcmpl-123","object":"chat.completion","created":1679367552,"model":"gpt-3.5-turbo-0301","usage":{"prompt_tokens":9,"completion_tokens":11,"total_tokens":20},"choices":[{"message":{"role":"assistant","content":"\n\nHi there! How may I assist you today?"},"finish_reason":"stop","index":0}]}`
)
