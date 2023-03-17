package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sashabaranov/go-openai"
)

func TestNewOpenAI(t *testing.T) {
	client := NewOpenAI("test")
	if client == nil {
		t.Error("NewOpenAI() returned nil")
	}
}

func TestGetCompletion(t *testing.T) {
	// Fake the response from the OpenAI API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Response from https://api.openai.com/v1/completions (prompt omitted)
		// Requested on 3/16/2023
		fakeResponse := `{"id":"cmpl-6uxqlIGXLAfPhBLZklA1uBLBho6gw","object":"text_completion","created":1679034451,"model":"code-davinci-002","choices":[{"text":"\nprint(\"Hello World\")","index":0,"logprobs":null,"finish_reason":"stop"}],"usage":{"prompt_tokens":6,"completion_tokens":6,"total_tokens":12}}`
		_, err := w.Write([]byte(fakeResponse))
		if err != nil {
			return
		}
	}))
	client := NewOpenAIWithBaseURL("test", ts.URL)
	completion, err := client.GetCompletion("# Print hello world in python.", openai.CodexCodeDavinci002)
	if err != nil {
		t.Errorf("GetCompletion() returned error: %v", err)
	}
	if completion != "\nprint(\"Hello World\")" {
		t.Errorf("GetCompletion() returned wrong completion: %v", completion)
	}
}
