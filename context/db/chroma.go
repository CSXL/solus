package db

import (
	"context"
	"encoding/json"

	chromadb "github.com/CSXL/go-chroma"
	"github.com/CSXL/solus/ai"
	"github.com/CSXL/solus/ai/openai"
)

type Metadatas map[string]interface{}

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

// CreateCollection creates a new collection
func (c *ChromaClient) CreateCollection(collection string) error {
	createCollectionRequest := chromadb.NewCreateCollection(collection)
	_, err := c.db.CreateCollection(createCollectionRequest)
	return err
}

type collectionsResponse struct {
	Name     string      `json:"name"`
	Metadata interface{} `json:"metadata"`
}

// ListCollections lists the collections
func (c *ChromaClient) ListCollections() ([]string, error) {
	response, err := c.db.ListCollections()
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var rawCollections []collectionsResponse
	err = json.NewDecoder(response.Body).Decode(&rawCollections)
	if err != nil {
		return nil, err
	}
	collections := make([]string, len(rawCollections))
	for i, collection := range rawCollections {
		collections[i] = collection.Name
	}
	return collections, nil
}

// DeleteCollection deletes a collection
func (c *ChromaClient) DeleteCollection(collection string) error {
	_, err := c.db.DeleteCollection(collection)
	return err
}

// Reset resets the database
func (c *ChromaClient) Reset() error {
	_, err := c.db.Reset()
	return err
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

// UpdateDocumentMetadatas updates the metadatas for the given documents
func (c *ChromaClient) UpdateDocumentsMetadata(collection string, documents []*Document) error {
	metadatas := make([]interface{}, len(documents))
	ids := make([]interface{}, len(documents))
	for i, document := range documents {
		metadatas[i] = document.Metadata
		ids[i] = document.ID
	}
	UpdateEmbeddingRequest := chromadb.UpdateEmbedding{
		Ids:       ids,
		Metadatas: metadatas,
	}
	_, err := c.db.Update(collection, &UpdateEmbeddingRequest)
	return err
}

// UpdateDocumentsContent updates the content and embeddings for the given documents
func (c *ChromaClient) UpdateDocumentsContent(collection string, documents []*Document) error {
	embeddings := make([]interface{}, len(documents))
	ids := make([]interface{}, len(documents))
	contents := make([]interface{}, len(documents))
	for i, document := range documents {
		embeddings[i] = document.Embedding
		ids[i] = document.ID
		contents[i] = document.Content
	}
	UpdateEmbeddingRequest := chromadb.UpdateEmbedding{
		Ids:        ids,
		Embeddings: embeddings,
		Documents:  contents,
	}
	_, err := c.db.Update(collection, &UpdateEmbeddingRequest)
	return err
}

// RemoveDocumentsByMetadata removes documents from the database by metadata
func (c *ChromaClient) RemoveDocumentsByMetadata(collection string, metadata Metadatas) error {
	DeleteEmbeddingRequest := chromadb.DeleteEmbedding{
		Where: metadata,
	}
	_, err := c.db.Delete(collection, &DeleteEmbeddingRequest)
	return err
}

// RemoveDocumentsByIDs removes documents from the database by IDs
func (c *ChromaClient) RemoveDocumentsByIDs(collection string, ids []string) error {
	idsInterface := make([]interface{}, len(ids))
	for i, id := range ids {
		idsInterface[i] = id
	}
	DeleteEmbeddingRequest := chromadb.DeleteEmbedding{
		Ids: idsInterface,
	}
	_, err := c.db.Delete(collection, &DeleteEmbeddingRequest)
	return err
}

type searchResponse struct {
	Ids       [][]string    `json:"ids"`
	Documents [][]string    `json:"documents"`
	Distances [][]float32   `json:"distances"`
	Metadatas [][]Metadatas `json:"metadatas"`
}

// Search searches the database based on a query
func (c *ChromaClient) Search(collection string, query string, limit int) ([]*Document, error) {
	queryEmbeddings, err := c.GetEmbeddings(query)
	if err != nil {
		return nil, err
	}
	queryInterface := make([]interface{}, len(queryEmbeddings))
	for i, v := range queryEmbeddings {
		queryInterface[i] = v
	}
	QueryEmbeddingRequest := chromadb.QueryEmbedding{
		QueryEmbeddings: queryInterface,
		NResults:        limit,
		Include:         []string{"metadatas", "documents", "distances", "embeddings"},
	}
	response, err := c.db.GetNearestNeighbors(collection, &QueryEmbeddingRequest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var searchResponse searchResponse
	err = json.NewDecoder(response.Body).Decode(&searchResponse)
	if err != nil {
		return nil, err
	}
	documents := make([]*Document, len(searchResponse.Ids))
	for i := range searchResponse.Ids {
		documents[i] = &Document{
			ID:        searchResponse.Ids[i][0],
			Metadata:  searchResponse.Metadatas[i][0],
			Content:   searchResponse.Documents[i][0],
			Embedding: nil,
		}
	}
	return documents, nil
}

type Document struct {
	ID        string
	Metadata  Metadatas
	Embedding []interface{}
	Content   string
}

func NewDocument(id string, metadata Metadatas, content string) *Document {
	return &Document{
		ID:       id,
		Metadata: metadata,
		Content:  content,
	}
}

func NewDocumentWithEmbedding(id string, metadata Metadatas, content string, embedding []float32) *Document {
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
