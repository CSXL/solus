package agent

import (
	"encoding/json"
	"fmt"

	"github.com/CSXL/solus/ai/openai"
	"github.com/google/logger"
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

func NewChatAgentConfig(openAIAPIKey string) *ChatAgentConfig {
	return &ChatAgentConfig{
		OpenAIAPIKey: openAIAPIKey,
	}
}

type ChatAgentMessage struct {
	Type    ChatAgentMessageType
	Role    ChatAgentMessageRole
	Content string
}

func NewChatAgentMessage(msgType ChatAgentMessageType, role ChatAgentMessageRole, content string) *ChatAgentMessage {
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

func ChatAgentMessageFromJSON(jsonStr string) (*ChatAgentMessage, error) {
	var msg ChatAgentMessage
	err := json.Unmarshal([]byte(jsonStr), &msg)
	return &msg, err
}

func ChatAgentMessageFromOpenAIChatMessage(msg openai.ChatMessage) *ChatAgentMessage {
	return NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRole(msg.Role), msg.Content)
}

type ChatAgentMessageContent struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func NewChatAgentMessageContent(msgType string, content string) *ChatAgentMessageContent {
	return &ChatAgentMessageContent{
		Type:    msgType,
		Content: content,
	}
}

func (c *ChatAgentMessageContent) ToJSON() (string, error) {
	json, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

func ChatAgentMessageContentFromJSON(jsonMessage string) (*ChatAgentMessageContent, error) {
	var content ChatAgentMessageContent
	err := json.Unmarshal([]byte(jsonMessage), &content)
	if err != nil {
		return nil, err
	}
	return &content, nil
}

type ChatAgent struct {
	*Agent
	OpenAIChatClient *openai.ChatClient
	Messages         []ChatAgentMessage
}

func NewChatAgent(name string, config *ChatAgentConfig) *ChatAgent {
	return &ChatAgent{
		Agent:            NewAgent(name, ChatAgentType, config),
		OpenAIChatClient: openai.NewChatClient(config.OpenAIAPIKey),
		Messages:         []ChatAgentMessage{},
	}
}

func (c *ChatAgent) AddMessage(msg ChatAgentMessage) {
	c.Messages = append(c.Messages, msg)
}

func (c *ChatAgent) GetMessages() []ChatAgentMessage {
	return c.Messages
}

func (c *ChatAgent) GetLastMessage() ChatAgentMessage {
	if len(c.Messages) == 0 {
		return ChatAgentMessage{}
	}
	return c.Messages[len(c.Messages)-1]
}

func (c *ChatAgent) GetLastMessageContent() string {
	return c.GetLastMessage().Content
}

func (c *ChatAgent) GetLastMessageRole() ChatAgentMessageRole {
	return c.GetLastMessage().Role
}

func (c *ChatAgent) SendMessageToAgent(msg ChatAgentMessage) (*ChatAgentTask, error) {
	if !c.IsRunning() {
		logger.Info("Note: Agent is not running, message will be queued but not sent.")
	}
	sendTask, err := NewChatAgentTask(c, ChatAgentTaskTypeSendMessage, msg)
	if err != nil {
		return nil, err
	}
	err = c.AddTask(sendTask)
	if err != nil {
		return nil, err
	}
	return sendTask, err
}

func (c *ChatAgent) SendMessage(msg ChatAgentMessage) (*ChatAgentMessage, error) {
	messageTask, err := c.SendMessageToAgent(msg)
	if err != nil {
		return nil, err
	}
	messageTask.AwaitCompletion()
	aiResponseMessage := messageTask.GetResult().(*ChatAgentMessage)
	return aiResponseMessage, nil
}

type ChatAgentTaskType string
type ChatAgentTaskPayload interface{}

const (
	ChatAgentTaskTypeSendMessage ChatAgentTaskType = "send_message"
)

type ChatAgentTask struct {
	*AgentTask
}

func NewChatAgentTask(agent *ChatAgent, taskType ChatAgentTaskType, payload ChatAgentTaskPayload) (*ChatAgentTask, error) {
	isSequential := true
	agentTaskType := NewAgentTaskType(string(taskType), isSequential)
	handler, err := BuildChatAgentHandler(agent, taskType, payload)
	if err != nil {
		return nil, err
	}
	return &ChatAgentTask{
		AgentTask: NewAgentTask(string(taskType), agentTaskType, handler),
	}, nil
}

func BuildChatAgentHandler(agent *ChatAgent, taskType ChatAgentTaskType, payload ChatAgentTaskPayload) (HandlerFunction, error) {
	switch taskType {
	case ChatAgentTaskTypeSendMessage:
		msg := payload.(ChatAgentMessage)
		return BuildChatAgentMessageHandler(agent, msg), nil
	default:
		return nil, fmt.Errorf("unknown task type: %s", taskType)
	}
}

func BuildChatAgentMessageHandler(agent *ChatAgent, msg ChatAgentMessage) HandlerFunction {
	return func(kill chan bool) interface{} {
		agent.AddMessage(msg)
		err := agent.OpenAIChatClient.SendMessage(msg.Content, string(msg.Role))
		if err != nil {
			return err
		}
		openaiResponse := agent.OpenAIChatClient.GetLastMessage()
		serializedResponse := ChatAgentMessageFromOpenAIChatMessage(openaiResponse)
		agent.Messages = append(agent.Messages, *serializedResponse)
		return serializedResponse
	}
}
