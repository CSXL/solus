package requirements

import (
	"testing"

	"github.com/CSXL/solus/ai/openai"
	"github.com/stretchr/testify/assert"
)

func TestNewRequirementsGenerator(t *testing.T) {
	testInputPrompt := "test input prompt"
	testInputConversation := "test input conversation"
	testOpenAIKey := "test key"
	testConfig := NewRequirementsConfig(testInputPrompt, testOpenAIKey)
	testGenerator := NewRequirementsGenerator(testInputConversation, testConfig)
	assert.NotNil(t, testGenerator)
}

func TestRequirementsGenerator_Generate(t *testing.T) {
	testInputPrompt := "test input prompt"
	testInputConversation := "test input conversation"
	testOpenAIKey := "test key"
	testConfig := NewRequirementsConfig(testInputPrompt, testOpenAIKey)
	testGenerator := NewRequirementsGenerator(testInputConversation, testConfig)
	assert.NotNil(t, testGenerator)
	ts := openai.StartHTTPTestServer(openai.SampleChatYAMLCompletion)
	defer ts.Close()
	testGenerator.Conversation.GetAgent().OpenAIChatClient.SetBaseURL(ts.URL)
	_, _ = testGenerator.Generate()
	assert.NotNil(t, testGenerator.GeneratedRequirements)
}
