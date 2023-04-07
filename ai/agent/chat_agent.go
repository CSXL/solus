package agent

import (
	"encoding/json"

	"github.com/CSXL/solus/ai/openai"
)

type ChatAgentMessageRole string
type ChatAgentMessageType string

const (
	// ChatAgent Metadata
	ChatAgentType = "chat"
	// ChatAgent Message Roles
	ChatAgentMessageRoleUser      ChatAgentMessageRole = "user"
	ChatAgentMessageRoleAssistant ChatAgentMessageRole = "assistant"
	ChatAgentMessageRoleSystem    ChatAgentMessageRole = "system"
	// ChatAgent Message Types
	ChatAgentMessageTypeText ChatAgentMessageType = "text"
	ChatAgentMessageTypeFile ChatAgentMessageType = "file"
	ChatAgentMessageTypeLink ChatAgentMessageType = "link"
)

type ChatAgentConfig struct {
	OpenAIAPIKey string
}

func (c *ChatAgentConfig) NewChatAgentConfig(openAIAPIKey string) *ChatAgentConfig {
	return &ChatAgentConfig{
		OpenAIAPIKey: openAIAPIKey,
	}
}

type ChatAgentMessage struct {
	Type    ChatAgentMessageType
	Role    ChatAgentMessageRole
	Content string
}

func (c *ChatAgentMessage) NewChatAgentMessage(msgType ChatAgentMessageType, role ChatAgentMessageRole, content string) *ChatAgentMessage {
	return &ChatAgentMessage{
		Type:    msgType,
		Role:    role,
		Content: content,
	}
}

func (c *ChatAgentMessage) IsMessageOfType(msgType ChatAgentMessageType) bool {
	return c.Type == msgType
}

func (c *ChatAgentMessage) IsTextMessage() bool {
	return c.IsMessageOfType(ChatAgentMessageTypeText)
}

func (c *ChatAgentMessage) IsFileMessage() bool {
	return c.IsMessageOfType(ChatAgentMessageTypeFile)
}

func (c *ChatAgentMessage) IsLinkMessage() bool {
	return c.IsMessageOfType(ChatAgentMessageTypeLink)
}

func (c *ChatAgentMessage) IsMessageOfRole(role ChatAgentMessageRole) bool {
	return c.Role == role
}

func (c *ChatAgentMessage) IsUserMessage() bool {
	return c.IsMessageOfRole(ChatAgentMessageRoleUser)
}

func (c *ChatAgentMessage) IsAssistantMessage() bool {
	return c.IsMessageOfRole(ChatAgentMessageRoleAssistant)
}

func (c *ChatAgentMessage) IsSystemMessage() bool {
	return c.IsMessageOfRole(ChatAgentMessageRoleSystem)
}

func (c *ChatAgentMessage) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(c)
	return string(jsonBytes), err
}

type ChatAgent struct {
	*Agent
	OpenAIChatClient *openai.ChatClient
	Messages         []ChatAgentMessage
}

func (c *ChatAgent) NewChatAgent(id string, name string, config *ChatAgentConfig) *ChatAgent {
	return &ChatAgent{
		Agent:            NewAgent(name, ChatAgentType, config),
		OpenAIChatClient: openai.NewChatClient(config.OpenAIAPIKey),
	}
}

func (c *ChatAgent) AddMessage(msg ChatAgentMessage) {
	c.Messages = append(c.Messages, msg)
}

func (c *ChatAgent) GetMessages() []ChatAgentMessage {
	return c.Messages
}

func (c *ChatAgent) GetLastMessage() ChatAgentMessage {
	return c.Messages[len(c.Messages)-1]
}

func (c *ChatAgent) GetLastMessageContent() string {
	return c.GetLastMessage().Content
}

func (c *ChatAgent) GetLastMessageRole() ChatAgentMessageRole {
	return c.GetLastMessage().Role
}
