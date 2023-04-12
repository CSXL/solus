package db

import (
	"context"

	chromadb "github.com/CSXL/go-chroma"
)

// ChromaClient is a client for interacting with the Chroma database
type ChromaClient struct {
	context *context.Context
	db      *chromadb.Client
}

// NewChromaClient creates a new ChromaClient
func NewChromaClient(ctx *context.Context, basePath string) (*ChromaClient, error) {
	chromadbClient, err := chromadb.NewClient(basePath)
	return &ChromaClient{
		context: ctx,
		db:      chromadbClient,
	}, err
}

// GetChromaClient returns the internal chromadb client
func (c *ChromaClient) GetChromaDB() *chromadb.Client {
	return c.db
}

// GetContext returns the context
func (c *ChromaClient) GetContext() *context.Context {
	return c.context
}

// AddDocument adds a document to the database
