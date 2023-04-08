// Endpoints for mocking OpenAI API calls.
package openai

import (
	"net/http"
	"net/http/httptest"
)

type OpenAIResponse string

const (
	SampleChatCompletion     OpenAIResponse = `{"id":"chatcmpl-123","object":"chat.completion","created":1679367552,"model":"gpt-3.5-turbo-0301","usage":{"prompt_tokens":9,"completion_tokens":11,"total_tokens":20},"choices":[{"message":{"role":"assistant","content":"\n\nHi there! How may I assist you today?"},"finish_reason":"stop","index":0}]}`
	SampleChatJSONCompletion OpenAIResponse = `{"id":"chatcmpl-123","object":"chat.completion","created":1679367552,"model":"gpt-3.5-turbo-0301","usage":{"prompt_tokens":9,"completion_tokens":11,"total_tokens":20},"choices":[{"message":{"role":"assistant","content":"{\"type\":\"message\",\"content\":\"CSX Labs is an amazing organization.\"}"},"finish_reason":"stop","index":0}]}`
)

func StartHTTPTestServer(response OpenAIResponse) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-Type", "application/json")
		_, err := w.Write([]byte(response))
		if err != nil {
			return
		}
	}))
}
