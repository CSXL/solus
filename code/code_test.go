package code

import (
	"os"
	"testing"

	"github.com/CSXL/solus/ai/openai"
	"github.com/stretchr/testify/assert"
)

func TestNewCodeGenerator(t *testing.T) {
	testGenerationFolder := "test generation folder"
	testOpenAIKey := "test key"
	testConfig := NewCodeConfig(testGenerationFolder, testOpenAIKey)
	testGenerator := NewCodeGenerator(testGenerationFolder, testConfig)
	assert.NotNil(t, testGenerator)
}

func TestCodeGenerator_Generate(t *testing.T) {
	testGenerationFolder, err := os.MkdirTemp("", "test_generation_folder")
	assert.Nil(t, err)
	defer os.RemoveAll(testGenerationFolder)
	testOpenAIKey := "test key"
	testConfig := NewCodeConfig(testGenerationFolder, testOpenAIKey)
	testGenerator := NewCodeGenerator(testGenerationFolder, testConfig)
	assert.NotNil(t, testGenerator)
	ts := openai.StartHTTPTestServer(openai.SampleChatFileCompletion)
	defer ts.Close()
	testGenerator.Conversation.GetAgent().OpenAIChatClient.SetBaseURL(ts.URL)
	err = testGenerator.Generate()
	assert.Nil(t, err)
	assert.NotNil(t, testGenerator.ProjectState)
}
