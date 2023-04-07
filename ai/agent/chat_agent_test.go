package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CSXL/solus/ai/openai"
	"github.com/stretchr/testify/assert"
)

func TestNewChatAgentConfig(t *testing.T) {
	chatAgentConfig := NewChatAgentConfig("test-key")
	assert.Equal(t, "test-key", chatAgentConfig.OpenAIAPIKey)
}

func TestNewChatAgentMessage(t *testing.T) {
	chatAgentMessage := NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content")
	assert.Equal(t, ChatAgentMessageTypeText, chatAgentMessage.Type)
	assert.Equal(t, ChatAgentMessageRoleUser, chatAgentMessage.Role)
	assert.Equal(t, "test-content", chatAgentMessage.Content)
}

func TestChatAgentMessage_IsMessageOfType(t *testing.T) {
	chatAgentMessage := NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content")
	assert.True(t, chatAgentMessage.IsMessageOfType(ChatAgentMessageTypeText))
	assert.False(t, chatAgentMessage.IsMessageOfType(ChatAgentMessageTypeFile))
}

func TestChatAgentMessage_IsXTypeMessage(t *testing.T) {
	chatAgentMessage := NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content")
	assert.True(t, chatAgentMessage.IsTextMessage())
	assert.False(t, chatAgentMessage.IsFileMessage())
	assert.False(t, chatAgentMessage.IsLinkMessage())
}

func TestChatAgentMessage_IsMessageOfRole(t *testing.T) {
	chatAgentMessage := NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content")
	assert.True(t, chatAgentMessage.IsMessageOfRole(ChatAgentMessageRoleUser))
	assert.False(t, chatAgentMessage.IsMessageOfRole(ChatAgentMessageRoleAssistant))
}

func TestChatAgentMessage_IsXRoleMessage(t *testing.T) {
	chatAgentMessage := NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content")
	assert.True(t, chatAgentMessage.IsUserMessage())
	assert.False(t, chatAgentMessage.IsAssistantMessage())
	assert.False(t, chatAgentMessage.IsSystemMessage())
}

func TestChatAgentMessage_ToJSON(t *testing.T) {
	chatAgentMessage := NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content")
	json, err := chatAgentMessage.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, `{"Type":"text","Role":"user","Content":"test-content"}`, json)
}

func TestNewChatAgent(t *testing.T) {
	NewChatAgent("testAgent", NewChatAgentConfig("test-key"))
}

func TestChatAgent_AddMessage(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", NewChatAgentConfig("test-key"))
	chatAgent.AddMessage(*NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content"))
	assert.Equal(t, 1, len(chatAgent.Messages))
}

func TestChatAgent_GetMessages(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", NewChatAgentConfig("test-key"))
	chatAgent.AddMessage(*NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content"))
	assert.Equal(t, 1, len(chatAgent.GetMessages()))
}

func TestChatAgent_GetLastMessage(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", NewChatAgentConfig("test-key"))
	chatAgent.AddMessage(*NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content"))
	assert.Equal(t, "test-content", chatAgent.GetLastMessage().Content)
}

func TestChatAgent_GetLastMessageContent(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", NewChatAgentConfig("test-key"))
	chatAgent.AddMessage(*NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content"))
	assert.Equal(t, "test-content", chatAgent.GetLastMessageContent())
}

func TestChatAgent_GetLastMessageRole(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", NewChatAgentConfig("test-key"))
	chatAgent.AddMessage(*NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content"))
	assert.Equal(t, ChatAgentMessageRoleUser, chatAgent.GetLastMessageRole())
}

func TestChatAgent_GetLastMessageWithEmptyMessages(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", NewChatAgentConfig("test-key"))
	assert.Equal(t, ChatAgentMessage{}, chatAgent.GetLastMessage())
}

func TestChatAgent_SendMessage(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", NewChatAgentConfig("test-key"))
	chatAgent.Start()
	defer chatAgent.Kill()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fakeResponse := openai.SampleChatCompletion
		w.Write([]byte(fakeResponse))
	}))
	defer ts.Close()
	chatAgent.OpenAIChatClient.SetBaseURL(ts.URL)
	msg := NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content")
	messageTask, err := chatAgent.SendMessage(*msg)
	assert.Nil(t, err)
	messageTask.AwaitCompletion()
	assert.Equal(t, 2, len(chatAgent.Messages))
}
