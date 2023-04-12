package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewChromaClient(t *testing.T) {
	ctx := context.Background()
	client, err := NewChromaClient(&ctx, "https://example.com")
	assert.Nil(t, err)
	assert.NotNil(t, client)
}
