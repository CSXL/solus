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
	db           *chromadb.Client
	aiConfig     *ai.AIConfig
	openAIClient *openai.OpenAI
}

// NewChromaClient creates a new ChromaClient
func NewChromaClient(ctx *context.Context, basePath string, aiConfig ai.AIConfig) (*ChromaClient, error) {
	chromadbClient, err := chromadb.NewClient(basePath)
	if err != nil {
		return nil, err
	}
	openAIClient := openai.NewOpenAI(aiConfig.OpenAIAPIKey)
	return &ChromaClient{
		context:      ctx,
		db:           chromadbClient,
		aiConfig:     &aiConfig,
		openAIClient: openAIClient,
	}, nil
}

// GetChromaClient returns the internal chromadb client
func (c *ChromaClient) GetChromaDB() *chromadb.Client {
	return c.db
}

// GetContext returns the context
func (c *ChromaClient) GetContext() *context.Context {
	return c.context
}

// GetEmbeddings returns the embeddings for the given texts
func (c *ChromaClient) GetEmbeddings(texts []string) ([][]float32, error) {
	return c.openAIClient.GetEmbeddings(texts)
}

// AddDocuments adds documents to the database
func (c *ChromaClient) AddDocuments(collection string, documents []*Document) error {
	IDs := make([]interface{}, len(documents))
	Metadatas := make([]interface{}, len(documents))
	Embeddings := make([]interface{}, len(documents))
	Contents := make([]interface{}, len(documents))
	for i, document := range documents {
		IDs[i] = document.ID
		Metadatas[i] = document.Metadata
		if document.Embedding == nil {
			generatedEmbeddings, err := c.GetEmbeddings([]string{document.Content})
			if err != nil {
				return err
			}
			Embeddings[i] = generatedEmbeddings[0]
		} else {
			Embeddings[i] = document.Embedding
		}
		Contents[i] = document.Content
	}
	chromaDocuments := chromadb.AddEmbedding_Documents{}
	err := chromaDocuments.FromAddEmbeddingDocuments1(chromadb.AddEmbeddingDocuments1(Contents))
	if err != nil {
		return err
	}
	chromaEmbeddings := Embeddings
	chromaIds := chromadb.AddEmbedding_Ids{}
	err = chromaIds.FromAddEmbeddingIds1(IDs)
	if err != nil {
		return err
	}
	chromaMetadatas := chromadb.AddEmbedding_Metadatas{}
	err = chromaMetadatas.FromAddEmbeddingMetadatas0(Metadatas)
	if err != nil {
		return err
	}
	IncrementIndex := true
	addRequest := chromadb.AddEmbedding{
		Documents:      &chromaDocuments,
		Embeddings:     chromaEmbeddings,
		Ids:            &chromaIds,
		IncrementIndex: &IncrementIndex,
		Metadatas:      &chromaMetadatas,
	}
	_, err = c.db.Add(*c.GetContext(), collection, addRequest)
	return err
}

type Document struct {
	ID        string
	Metadata  []interface{}
	Embedding []interface{}
	Content   string
}

func NewDocument(id string, metadata []interface{}, content string) *Document {
	return &Document{
		ID:       id,
		Metadata: metadata,
		Content:  content,
	}
}

func NewDocumentWithEmbedding(id string, metadata []interface{}, content string, embedding []interface{}) *Document {
	return &Document{
		ID:        id,
		Metadata:  metadata,
		Content:   content,
		Embedding: embedding,
	}
}
