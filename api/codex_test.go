package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewCodexClient(t *testing.T) {
	client := NewCodexClient("test")
	if client == nil {
		t.Error("NewCodexClient() returned nil")
	}
}

func TestExecuteCodexCompletion(t *testing.T) {
	client := NewCodexClient("test")
	// Fake the response from the OpenAI API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Response from https://api.openai.com/v1/completions (prompt omitted)
		// Requested on 3/16/2023
		fake_response := `{"id":"cmpl-6uxqlIGXLAfPhBLZklA1uBLBho6gw","object":"text_completion","created":1679034451,"model":"code-davinci-002","choices":[{"text":"\nprint(\"Hello World\")","index":0,"logprobs":null,"finish_reason":"stop"}],"usage":{"prompt_tokens":6,"completion_tokens":6,"total_tokens":12}}`
		w.Write([]byte(fake_response))
	}))
	client.openAIClient = NewOpenAIWithBaseURL("test", ts.URL)
	completion, err := client.ExecuteCodexCompletion("# Print hello world in python.")
	if err != nil {
		t.Errorf("ExecuteCodexCompletion() returned error: %v", err)
	}
	if completion != "\nprint(\"Hello World\")" {
		t.Errorf("ExecuteCodexCompletion() returned wrong completion: %v", completion)
	}
}
