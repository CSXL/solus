package db

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CSXL/solus/ai"
	"github.com/stretchr/testify/assert"
)

func TestNewChromaClient(t *testing.T) {
	ctx := context.Background()
	aiConfig := ai.NewAIConfig("testKey")
	client, err := NewChromaClient(&ctx, "https://example.com", *aiConfig)
	assert.Nil(t, err)
	assert.NotNil(t, client)
}

func TestGetChromaDB(t *testing.T) {
	ctx := context.Background()
	aiConfig := ai.NewAIConfig("testKey")
	client, err := NewChromaClient(&ctx, "https://example.com", *aiConfig)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.GetChromaDB())
}

func TestGetContext(t *testing.T) {
	ctx := context.Background()
	aiConfig := ai.NewAIConfig("testKey")
	client, err := NewChromaClient(&ctx, "https://example.com", *aiConfig)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.GetContext())
}

func TestAddDocuments(t *testing.T) {
	ctx := context.Background()
	aiConfig := ai.NewAIConfig("testKey")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	client, err := NewChromaClient(&ctx, ts.URL, *aiConfig)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	documents := []*Document{
		NewDocumentWithEmbedding("docid", map[string]interface{}{"mymetadatakey": "myvalue"}, "documentcontent", []float32{0.1, 0.2, 0.3}),
	}
	err = client.AddDocuments("test", documents)
	assert.Nil(t, err)
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
	ctx := context.Background()
	aiConfig := ai.NewAIConfig("testKey")
	client, err := NewChromaClient(&ctx, "https://example.com", *aiConfig)
	client.ChangeOpenAIBaseURL(ts.URL)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	embeddings, err := client.GetEmbeddings("Hello World")
	assert.Nil(t, err)
	assert.NotNil(t, embeddings)
}
