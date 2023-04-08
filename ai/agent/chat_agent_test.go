package agent

import (
	"testing"

	"github.com/CSXL/solus/ai"
	"github.com/CSXL/solus/ai/openai"
	"github.com/stretchr/testify/assert"
)

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
	assert.False(t, chatAgentMessage.IsQueryMessage())
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

func TestChatAgentMessage_FromJSON(t *testing.T) {
	chatAgentMessage := NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content")
	json, err := chatAgentMessage.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, `{"Type":"text","Role":"user","Content":"test-content"}`, json)
	chatAgentMessage, err = ChatAgentMessageFromJSON(json)
	assert.Nil(t, err)
	assert.Equal(t, ChatAgentMessageTypeText, chatAgentMessage.Type)
	assert.Equal(t, ChatAgentMessageRoleUser, chatAgentMessage.Role)
	assert.Equal(t, "test-content", chatAgentMessage.Content)
}

func TestNewChatAgentMessageContent(t *testing.T) {
	msgType := "user_message"
	msgContent := "test-content"
	chatAgentMessageContent := NewChatAgentMessageContent(msgType, msgContent)
	assert.Equal(t, msgContent, chatAgentMessageContent.Content)
}

func TestChatAgentMessageContent_ToJSON(t *testing.T) {
	msgType := "user_message"
	msgContent := "test-content"
	chatAgentMessageContent := NewChatAgentMessageContent(msgType, msgContent)
	json, err := chatAgentMessageContent.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, `{"type":"user_message","content":"test-content"}`, json)
}

func TestChatAgentMessageContent_FromJSON(t *testing.T) {
	msgType := "user_message"
	msgContent := "test-content"
	chatAgentMessageContent := NewChatAgentMessageContent(msgType, msgContent)
	json, err := chatAgentMessageContent.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, `{"type":"user_message","content":"test-content"}`, json)
	chatAgentMessageContent, err = ChatAgentMessageContentFromJSON(json)
	assert.Nil(t, err)
	assert.Equal(t, msgType, chatAgentMessageContent.Type)
	assert.Equal(t, msgContent, chatAgentMessageContent.Content)
}

func TestNewChatAgent(t *testing.T) {
	NewChatAgent("testAgent", ai.NewAIConfig("test-key"))
}

func TestChatAgent_AddMessage(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", ai.NewAIConfig("test-key"))
	chatAgent.AddMessage(*NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content"))
	assert.Equal(t, 1, len(chatAgent.Messages))
}

func TestChatAgent_GetMessages(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", ai.NewAIConfig("test-key"))
	chatAgent.AddMessage(*NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content"))
	assert.Equal(t, 1, len(chatAgent.GetMessages()))
}

func TestChatAgent_SetMessages(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", ai.NewAIConfig("test-key"))
	chatAgent.AddMessage(*NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content"))
	chatAgent.SetMessages([]ChatAgentMessage{})
	assert.Equal(t, 0, len(chatAgent.GetMessages()))
}

func TestChatAgent_ResetMessages(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", ai.NewAIConfig("test-key"))
	chatAgent.AddMessage(*NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content"))
	chatAgent.ResetMessages()
	assert.Equal(t, 0, len(chatAgent.GetMessages()))
}

func TestChatAgent_GetLastMessage(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", ai.NewAIConfig("test-key"))
	chatAgent.AddMessage(*NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content"))
	assert.Equal(t, "test-content", chatAgent.GetLastMessage().Content)
}

func TestChatAgent_GetLastMessageContent(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", ai.NewAIConfig("test-key"))
	chatAgent.AddMessage(*NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content"))
	assert.Equal(t, "test-content", chatAgent.GetLastMessageContent())
}

func TestChatAgent_GetLastMessageRole(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", ai.NewAIConfig("test-key"))
	chatAgent.AddMessage(*NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content"))
	assert.Equal(t, ChatAgentMessageRoleUser, chatAgent.GetLastMessageRole())
}

func TestChatAgent_GetLastMessageWithEmptyMessages(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", ai.NewAIConfig("test-key"))
	assert.Equal(t, ChatAgentMessage{}, chatAgent.GetLastMessage())
}

func TestChatAgent_SendMessageToAgent(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", ai.NewAIConfig("test-key"))
	chatAgent.Start()
	defer chatAgent.Kill()
	ts := openai.StartHTTPTestServer(openai.SampleChatCompletion)
	defer ts.Close()
	chatAgent.OpenAIChatClient.SetBaseURL(ts.URL)
	msg := NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content")
	messageTask, err := chatAgent.sendMessageToAgent(*msg)
	assert.Nil(t, err)
	messageTask.AwaitCompletion()
	assert.Equal(t, 2, len(chatAgent.Messages))
}

func TestChatAgent_SendMessage(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", ai.NewAIConfig("test-key"))
	chatAgent.Start()
	defer chatAgent.Kill()
	ts := openai.StartHTTPTestServer(openai.SampleChatCompletion)
	defer ts.Close()
	chatAgent.OpenAIChatClient.SetBaseURL(ts.URL)
	msg := NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content")
	aiResponse, err := chatAgent.sendMessage(*msg)
	assert.Nil(t, err)
	assert.NotNil(t, aiResponse)
	assert.Equal(t, 2, len(chatAgent.Messages))
}

func TestChatAgent_SendChatMessage(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", ai.NewAIConfig("test-key"))
	chatAgent.Start()
	defer chatAgent.Kill()
	ts := openai.StartHTTPTestServer(openai.SampleChatCompletion)
	defer ts.Close()
	chatAgent.OpenAIChatClient.SetBaseURL(ts.URL)
	msg := NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content")
	aiResponse, err := chatAgent.SendChatMessage(*msg)
	assert.Nil(t, err)
	assert.NotNil(t, aiResponse)
	assert.Equal(t, 2, len(chatAgent.Messages))
}

func TestChatAgent_SendChatMessageAndWriteResponseToChannel(t *testing.T) {
	chatAgent := NewChatAgent("testAgent", ai.NewAIConfig("test-key"))
	chatAgent.Start()
	defer chatAgent.Kill()
	ts := openai.StartHTTPTestServer(openai.SampleChatCompletion)
	defer ts.Close()
	chatAgent.OpenAIChatClient.SetBaseURL(ts.URL)
	msg := NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRoleUser, "test-content")
	messageChannel := make(chan ChatAgentMessage)
	go chatAgent.SendChatMessageAndWriteResponseToChannel(*msg, messageChannel)
	aiResponse := <-messageChannel
	assert.NotNil(t, aiResponse)
	assert.Equal(t, 2, len(chatAgent.Messages))
}
