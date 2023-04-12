package db

import (
	"context"
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
	client, err := NewChromaClient(&ctx, "https://example.com", *aiConfig)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	err = client.AddDocuments("test", []*Document{})
	assert.Nil(t, err)
}
