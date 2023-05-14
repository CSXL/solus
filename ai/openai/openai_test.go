package openai

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

func TestGetEmbeddings(t *testing.T) {
	// Fake the response from the OpenAI API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Example response pulled from https://platform.openai.com/docs/guides/embeddings/how-to-get-embeddings
		// Requested on 4/11/2023
		fakeResponse := `{"data":[{"embedding":[-0.006929283495992422,-0.005336422007530928,-0.00004547132266452536,-0.024047505110502243],"index":0,"object":"embedding"}],"model":"text-embedding-ada-002","object":"list","usage":{"prompt_tokens":5,"total_tokens":5}}`
		_, err := w.Write([]byte(fakeResponse))
		if err != nil {
			return
		}
	}))
	client := NewOpenAIWithBaseURL("test", ts.URL)
	embeddings, err := client.GetEmbeddings([]string{"Hello World"})
	if err != nil {
		t.Errorf("GetEmbeddings() returned error: %v", err)
	}
	if len(embeddings) != 1 {
		t.Errorf("GetEmbeddings() returned wrong number of embeddings: %v", len(embeddings))
	}
	if len(embeddings[0]) != 4 {
		t.Errorf("GetEmbeddings() returned wrong number of dimensions: %v", len(embeddings[0]))
	}
}
