package db

import (
	"context"

	chromadb "github.com/CSXL/go-chroma"
	"github.com/CSXL/solus/ai"
	"github.com/CSXL/solus/ai/openai"
)

// ChromaClient is a client for interacting with the Chroma database
type ChromaClient struct {
	context      *context.Context
	db           *chromadb.ChromaClient
	aiConfig     *ai.AIConfig
	openAIClient *openai.OpenAI
}

// NewChromaClient creates a new ChromaClient
func NewChromaClient(ctx *context.Context, basePath string, aiConfig ai.AIConfig) (*ChromaClient, error) {
	chromadbClient := chromadb.NewChromaClient(basePath)
	openAIClient := openai.NewOpenAI(aiConfig.OpenAIAPIKey)
	return &ChromaClient{
		context:      ctx,
		db:           chromadbClient,
		aiConfig:     &aiConfig,
		openAIClient: openAIClient,
	}, nil
}

// GetChromaClient returns the internal chromadb client
func (c *ChromaClient) GetChromaDB() *chromadb.ChromaClient {
	return c.db
}

// GhangeOpenAIBaseURL changes the base URL for the OpenAI client
func (c *ChromaClient) ChangeOpenAIBaseURL(baseURL string) {
	c.openAIClient = openai.NewOpenAIWithBaseURL(c.aiConfig.OpenAIAPIKey, baseURL)
}

// GetContext returns the context
func (c *ChromaClient) GetContext() *context.Context {
	return c.context
}

// GetEmbeddings returns the embeddings for a given text
func (c *ChromaClient) GetEmbeddings(text string) ([]float32, error) {
	embeddings, err := c.openAIClient.GetEmbeddings([]string{text})
	if err != nil {
		return nil, err
	}
	if len(embeddings) > 0 {
		return embeddings[0], nil
	}
	return []float32{}, nil
}

// AddDocuments adds documents to the database
func (c *ChromaClient) AddDocuments(collection string, documents []*Document) error {
	embeddings := make([]interface{}, len(documents))
	metadatas := make([]interface{}, len(documents))
	ids := make([]interface{}, len(documents))
	contents := make([]interface{}, len(documents))
	for i, document := range documents {
		embeddings[i] = document.Embedding
		metadatas[i] = document.Metadata
		ids[i] = document.ID
		contents[i] = document.Content
	}
	AddEmbeddingRequest := chromadb.AddEmbedding{
		Embeddings:     embeddings,
		IncrementIndex: true,
		Metadatas:      metadatas,
		Ids:            ids,
		Documents:      contents,
	}
	_, err := c.db.Add(collection, &AddEmbeddingRequest)
	return err
}

type Document struct {
	ID        string
	Metadata  map[string]interface{}
	Embedding []interface{}
	Content   string
}

func NewDocument(id string, metadata map[string]interface{}, content string) *Document {
	return &Document{
		ID:       id,
		Metadata: metadata,
		Content:  content,
	}
}

func NewDocumentWithEmbedding(id string, metadata map[string]interface{}, content string, embedding []float32) *Document {
	embeddingInterface := make([]interface{}, len(embedding))
	for i, v := range embedding {
		embeddingInterface[i] = v
	}
	return &Document{
		ID:        id,
		Metadata:  metadata,
		Content:   content,
		Embedding: embeddingInterface,
	}
}
