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

func TestCreateCollection(t *testing.T) {
	ctx := context.Background()
	aiConfig := ai.NewAIConfig("testKey")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	client, err := NewChromaClient(&ctx, ts.URL, *aiConfig)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	err = client.CreateCollection("test")
	assert.Nil(t, err)
}

func TestListCollections(t *testing.T) {
	ctx := context.Background()
	aiConfig := ai.NewAIConfig("testKey")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Example response pulled from internal test
		// Requested on 5/14/2023
		fakeResponse := `[{"name":"test","metadata":null},{"name":"test2","metadata":null}]`
		_, err := w.Write([]byte(fakeResponse))
		if err != nil {
			return
		}
	}))
	defer ts.Close()
	client, err := NewChromaClient(&ctx, ts.URL, *aiConfig)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	collections, err := client.ListCollections()
	assert.Nil(t, err)
	assert.NotNil(t, collections)
}

func TestUpdateDocumentsMetadata(t *testing.T) {
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
		NewDocumentWithEmbedding("docid", Metadatas{"mymetadatakey": "myvalue"}, "documentcontent", []float32{0.1, 0.2, 0.3}),
	}
	err = client.UpdateDocumentsMetadata("test", documents)
	assert.Nil(t, err)
}

func TestUpdateDocumentsContent(t *testing.T) {
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
		NewDocumentWithEmbedding("docid", Metadatas{"mymetadatakey": "myvalue"}, "documentcontent", []float32{0.1, 0.2, 0.3}),
	}
	err = client.UpdateDocumentsContent("test", documents)
	assert.Nil(t, err)
}

func TestDeleteCollection(t *testing.T) {
	ctx := context.Background()
	aiConfig := ai.NewAIConfig("testKey")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	client, err := NewChromaClient(&ctx, ts.URL, *aiConfig)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	err = client.DeleteCollection("test")
	assert.Nil(t, err)
}

func TestReset(t *testing.T) {
	ctx := context.Background()
	aiConfig := ai.NewAIConfig("testKey")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	client, err := NewChromaClient(&ctx, ts.URL, *aiConfig)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	err = client.Reset()
	assert.Nil(t, err)
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
		NewDocumentWithEmbedding("docid", Metadatas{"mymetadatakey": "myvalue"}, "documentcontent", []float32{0.1, 0.2, 0.3}),
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

func TestRemoveDocumentsByMetadata(t *testing.T) {
	ctx := context.Background()
	aiConfig := ai.NewAIConfig("testKey")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	client, err := NewChromaClient(&ctx, ts.URL, *aiConfig)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	err = client.RemoveDocumentsByMetadata("test", Metadatas{"mymetadatakey": "myvalue"})
	assert.Nil(t, err)
}

func TestRemoveDocumentsByIDs(t *testing.T) {
	ctx := context.Background()
	aiConfig := ai.NewAIConfig("testKey")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()
	client, err := NewChromaClient(&ctx, ts.URL, *aiConfig)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	err = client.RemoveDocumentsByIDs("test", []string{"docid"})
	assert.Nil(t, err)
}

func TestSearch(t *testing.T) {
	ctx := context.Background()
	aiConfig := ai.NewAIConfig("testKey")
	// Fake the response from the OpenAI API
	// Fake the response from the OpenAI API
	ts_openai := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Example response pulled from https://platform.openai.com/docs/guides/embeddings/how-to-get-embeddings
		// Requested on 4/11/2023
		fakeResponse := `{"data":[{"embedding":[-0.006929283495992422,-0.005336422007530928,-0.00004547132266452536,-0.024047505110502243],"index":0,"object":"embedding"}],"model":"text-embedding-ada-002","object":"list","usage":{"prompt_tokens":5,"total_tokens":5}}`
		_, err := w.Write([]byte(fakeResponse))
		if err != nil {
			return
		}
	}))
	defer ts_openai.Close()
	// Fake the response from the Chroma Service
	ts_chroma := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Example response pulled from internal test
		// Requested on 5/14/2023
		fakeResponse := `{"ids":[["csxlabs.org","afd.enterprises"]],"embeddings":[[[1.0,2.0,3.0,4.0],[4.5,12.4,1.2,4.1]]],"documents":[["CSX Labs is a collection of open research laboratories...","AFD Enterprises is a Gs..."]],"metadatas":[[{"type":"website"},{"type":"website"}]],"distances":[[3,131]]}`
		_, err := w.Write([]byte(fakeResponse))
		if err != nil {
			return
		}
	}))
	defer ts_chroma.Close()
	client, err := NewChromaClient(&ctx, ts_chroma.URL, *aiConfig)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	client.ChangeOpenAIBaseURL(ts_openai.URL)
	results, err := client.Search("test", "test", 2)
	assert.Nil(t, err)
	assert.NotNil(t, results)
}
