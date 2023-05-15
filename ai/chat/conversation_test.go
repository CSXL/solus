package chat

import (
	"testing"

	"github.com/CSXL/solus/ai"
	"github.com/CSXL/solus/ai/agent"
	"github.com/CSXL/solus/ai/openai"

	"github.com/stretchr/testify/assert"
)

func TestNewConversation(t *testing.T) {
	convName := "test-conv"
	config := ai.NewAIConfig("test-openai-api-key")
	conversation := NewConversation(convName, config)
	assert.NotNil(t, conversation)
}

func TestConversationMessageOperations(t *testing.T) {
	convName := "test-conv"
	config := ai.NewAIConfig("test-openai-api-key")
	conversation := NewConversation(convName, config)
	message := agent.NewChatAgentMessage(agent.ChatAgentMessageTypeText, agent.ChatAgentMessageRoleUser, "test-content")
	for i := 0; i < 10; i++ {
		conversation.AddMessage(*message)
	}
	assert.Equal(t, 10, conversation.GetMessageCount())
	assert.Equal(t, 10, len(conversation.GetMessages()))
	conversation.ResetMessages()
	assert.Equal(t, 0, conversation.GetMessageCount())
	testMessages := []agent.ChatAgentMessage{
		*agent.NewChatAgentMessage(agent.ChatAgentMessageTypeText, agent.ChatAgentMessageRoleUser, "test-content"),
		*agent.NewChatAgentMessage(agent.ChatAgentMessageTypeText, agent.ChatAgentMessageRoleUser, "test-content"),
		*agent.NewChatAgentMessage(agent.ChatAgentMessageTypeText, agent.ChatAgentMessageRoleUser, "test-content"),
	}
	conversation.SetMessages(testMessages)
	assert.Equal(t, 3, conversation.GetMessageCount())
}

func TestConversation_PreemptiveStart(t *testing.T) {
	convName := "test-conv"
	config := ai.NewAIConfig("test-openai-api-key")
	conversation := NewConversation(convName, config)
	conversation.PreemptiveStart()
}

func TestConversation_startIfNotStarted(t *testing.T) {
	convName := "test-conv"
	config := ai.NewAIConfig("test-openai-api-key")
	conversation := NewConversation(convName, config)
	assert.NotNil(t, conversation)
	conversation.startIfNotStarted()
}

func TestConversation_Close(t *testing.T) {
	convName := "test-conv"
	config := ai.NewAIConfig("test-openai-api-key")
	conversation := NewConversation(convName, config)
	conversation.Close()
}

func TestConversation_Kill(t *testing.T) {
	convName := "test-conv"
	config := ai.NewAIConfig("test-openai-api-key")
	conversation := NewConversation(convName, config)
	conversation.Kill()
}

func TestConversation_SendUserMessage(t *testing.T) {
	convName := "test-conv"
	config := ai.NewAIConfig("test-openai-api-key")
	conversation := NewConversation(convName, config)
	testContent := "test-content"
	ts := openai.StartHTTPTestServer(openai.SampleChatCompletion)
	defer ts.Close()
	conversation.chatAgent.OpenAIChatClient.SetBaseURL(ts.URL)
	_, _ = conversation.SendUserMessage(testContent)
	assert.Equal(t, 2, conversation.GetMessageCount())
	assert.NotNil(t, conversation.GetLastMessage())
}

func TestConversation_SendSystemMessage(t *testing.T) {
	convName := "test-conv"
	config := ai.NewAIConfig("test-openai-api-key")
	conversation := NewConversation(convName, config)
	testContent := "test-content"
	ts := openai.StartHTTPTestServer(openai.SampleChatCompletion)
	defer ts.Close()
	conversation.chatAgent.OpenAIChatClient.SetBaseURL(ts.URL)
	_, err := conversation.SendSystemMessage(testContent)
	assert.Nil(t, err)
	assert.Equal(t, 2, conversation.GetMessageCount())
	assert.NotNil(t, conversation.GetLastMessage())
}
