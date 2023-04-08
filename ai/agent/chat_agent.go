package agent

import (
	"encoding/json"
	"fmt"

	"github.com/CSXL/solus/ai"
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
	ChatAgentMessageTypeText  ChatAgentMessageType = "text"
	ChatAgentMessageTypeFile  ChatAgentMessageType = "file"
	ChatAgentMessageTypeLink  ChatAgentMessageType = "link"
	ChatAgentMessageTypeQuery ChatAgentMessageType = "query"
)

type ChatAgentMessage struct {
	Type    ChatAgentMessageType
	Role    ChatAgentMessageRole
	Content string
}

// NewChatAgentMessage creates a new ChatAgentMessage. This message can be sent
// to the ChatAgent to get a response and completion.
//
// msgType: The type of message. Can be one of the following:
// - ChatAgentMessageTypeText
// - ChatAgentMessageTypeFile
// - ChatAgentMessageTypeLink
// - ChatAgentMessageTypeQuery
//
// role: The role of the message. Can be one of the following:
// - ChatAgentMessageRoleUser
// - ChatAgentMessageRoleAssistant
// - ChatAgentMessageRoleSystem
//
// content: The content of the message. The content can only be a string.
func NewChatAgentMessage(msgType ChatAgentMessageType, role ChatAgentMessageRole, content string) *ChatAgentMessage {
	return &ChatAgentMessage{
		Type:    msgType,
		Role:    role,
		Content: content,
	}
}

func (c *ChatAgentMessage) GetType() ChatAgentMessageType {
	return c.Type
}

func (c *ChatAgentMessage) GetRole() ChatAgentMessageRole {
	return c.Role
}

func (c *ChatAgentMessage) GetContent() string {
	return c.Content
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

func (c *ChatAgentMessage) IsQueryMessage() bool {
	return c.IsMessageOfType(ChatAgentMessageTypeQuery)
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

func (c *ChatAgentMessage) mutateContentFromNonJSONMessage() error {
	mutatedMsgContent := chatAgentMessageContentFromChatAgentMessage(*c)
	jsonMsgContent, err := mutatedMsgContent.ToJSON()
	if err != nil {
		return err
	}
	c.Type = ChatAgentMessageTypeText
	c.Content = jsonMsgContent
	return nil
}

func (c *ChatAgentMessage) Serialize() error {
	msgContent, err := chatAgentMessageContentFromJSON(c.GetContent())
	if err != nil {
		return nil
	}
	c.Type = ChatAgentMessageType(msgContent.Type)
	c.Content = msgContent.Content
	return nil
}

func (c *ChatAgentMessage) Marshal() error {
	err := c.Serialize()
	if err != nil {
		return err
	}
	return c.mutateContentFromNonJSONMessage()
}

func (c *ChatAgentMessage) ToOpenAIChatMessage() *openai.ChatMessage {
	return &openai.ChatMessage{
		Role:    string(c.Role),
		Content: c.Content,
	}
}

func ChatAgentMessageFromJSON(jsonStr string) (*ChatAgentMessage, error) {
	var msg ChatAgentMessage
	err := json.Unmarshal([]byte(jsonStr), &msg)
	return &msg, err
}

func ChatAgentMessageFromOpenAIChatMessage(c openai.ChatMessage) *ChatAgentMessage {
	return NewChatAgentMessage(ChatAgentMessageTypeText, ChatAgentMessageRole(c.Role), c.Content)
}

type chatAgentMessageContent struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func newChatAgentMessageContent(msgType string, content string) *chatAgentMessageContent {
	return &chatAgentMessageContent{
		Type:    msgType,
		Content: content,
	}
}

func (c *chatAgentMessageContent) ToJSON() (string, error) {
	json, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

func chatAgentMessageContentFromJSON(jsonMessage string) (*chatAgentMessageContent, error) {
	var content chatAgentMessageContent
	err := json.Unmarshal([]byte(jsonMessage), &content)
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func chatAgentMessageContentFromChatAgentMessage(msg ChatAgentMessage) *chatAgentMessageContent {
	return newChatAgentMessageContent(string(msg.Type), msg.Content)
}

type ChatAgent struct {
	*Agent
	OpenAIChatClient *openai.ChatClient
	Messages         []ChatAgentMessage
}

// NewChatAgent creates a new ChatAgent. The ChatAgent can be used to hold a
// conversation between a language model (currently using OpenAI's APIs), and
// the user.
//
// Remember to call the Start() method on the ChatAgent!
func NewChatAgent(name string, config *ai.AIConfig) *ChatAgent {
	return &ChatAgent{
		Agent:            NewAgent(name, ChatAgentType, config),
		OpenAIChatClient: openai.NewChatClient(config.OpenAIAPIKey),
		Messages:         []ChatAgentMessage{},
	}
}

func (c *ChatAgent) AddMessage(msg ChatAgentMessage) {
	msg.Serialize()
	c.Messages = append(c.Messages, msg)
	c.syncMessages()
}

func (c *ChatAgent) GetMessages() []ChatAgentMessage {
	return c.Messages
}

func (c *ChatAgent) SetMessages(msgs []ChatAgentMessage) {
	c.Messages = msgs
	c.syncMessages()
}

func (c *ChatAgent) getMarshalledMessages() []ChatAgentMessage {
	marshalledMessages := []ChatAgentMessage{}
	for _, msg := range c.Messages {
		msg.Marshal()
		marshalledMessages = append(marshalledMessages, msg)
	}
	return marshalledMessages
}

func (c *ChatAgent) getSerializedMessages() []ChatAgentMessage {
	serializedMessages := []ChatAgentMessage{}
	for _, msg := range c.Messages {
		msg.Serialize()
		serializedMessages = append(serializedMessages, msg)
	}
	return serializedMessages
}

func (c *ChatAgent) serializeAllMessages() {
	c.Messages = c.getSerializedMessages()
}

func (c *ChatAgent) syncMessages() {
	convertedMessages := []openai.ChatMessage{}
	marshalledMessages := c.getMarshalledMessages()
	for _, msg := range marshalledMessages {
		convertedMessages = append(convertedMessages, *msg.ToOpenAIChatMessage())
	}
	c.OpenAIChatClient.SetMessages(convertedMessages)
}

func (c *ChatAgent) ResetMessages() {
	c.Messages = []ChatAgentMessage{}
	c.syncMessages()
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

func (c *ChatAgent) sendMessageToAgent(msg ChatAgentMessage) (*ChatAgentTask, error) {
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

func (c *ChatAgent) sendMessage(msg ChatAgentMessage) (*ChatAgentMessage, error) {
	messageTask, err := c.sendMessageToAgent(msg)
	if err != nil {
		return nil, err
	}
	messageTask.AwaitCompletion()
	// Checks for error in result by attempting to type assert to ChatAgentMessage
	// If the type assertion fails, the task failed and we return the error
	_, ok := messageTask.GetResult().(*ChatAgentMessage)
	if !ok {
		return nil, messageTask.GetResult().(error)
	}
	aiResponseMessage := messageTask.GetResult().(*ChatAgentMessage)
	logger.Infof("Received chat message from ChatAgent <ID: %s, Name: %s>: %s", c.GetID(), c.GetName(), aiResponseMessage.Content)
	return aiResponseMessage, nil
}

// SendChatMessage marshalls and sends a ChatAgentMessage to the agent.
// The agent will respond with a ChatAgentMessage that is returned.
// Note: The content is serialized to JSON before sending in the schema:
//
//	{
//	  "type": string,   // Your message type (e.g. "text", "file", "link", "query")
//	  "content": string // Your message content (e.g. "Hello", "https://example.com", "What is the weather like in 2023?")
//	}
func (c *ChatAgent) SendChatMessage(msg ChatAgentMessage) (*ChatAgentMessage, error) {
	logger.Infof("Sending chat message to ChatAgent <ID: %s, Name: %s>: %s", c.GetID(), c.GetName(), msg.Content)
	msg.Marshal()
	aiMessage, err := c.sendMessage(msg)
	if err != nil {
		return nil, err
	}
	processedAIMessage, err := c.ProcessChatMessage(*aiMessage)
	if err != nil {
		return nil, err
	}
	c.Messages[len(c.Messages)-1] = processedAIMessage
	return &processedAIMessage, nil
}

// SendChatMessageAndWriteResponseToChannel serializes and sends a
// ChatAgentMessage to the agent.
// The agent will write the response to the provided channel.
// Note: The content is serialized to JSON before sending in the schema:
//
//	{
//	  "type": string,   // Your message type (e.g. "text", "file", "link", "query")
//	  "content": string // Your message content (e.g. "Hello", "https://example.com", "What is the weather like in 2023?")
//	}
func (c *ChatAgent) SendChatMessageAndWriteResponseToChannel(msg ChatAgentMessage, channel chan ChatAgentMessage) error {
	response, err := c.SendChatMessage(msg)
	if err != nil {
		return err
	}
	channel <- *response
	return nil
}

func (c *ChatAgent) ProcessChatMessage(msg ChatAgentMessage) (ChatAgentMessage, error) {
	err := msg.Serialize()
	return msg, err
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
	handler, err := buildChatAgentHandler(agent, taskType, payload)
	if err != nil {
		return nil, err
	}
	return &ChatAgentTask{
		AgentTask: NewAgentTask(string(taskType), agentTaskType, handler),
	}, nil
}

func buildChatAgentHandler(agent *ChatAgent, taskType ChatAgentTaskType, payload ChatAgentTaskPayload) (HandlerFunction, error) {
	switch taskType {
	case ChatAgentTaskTypeSendMessage:
		msg := payload.(ChatAgentMessage)
		return buildChatAgentMessageHandler(agent, msg), nil
	default:
		return nil, fmt.Errorf("unknown task type: %s", taskType)
	}
}

func buildChatAgentMessageHandler(agent *ChatAgent, msg ChatAgentMessage) HandlerFunction {
	return func(kill chan bool) interface{} {
		agent.AddMessage(msg)
		err := agent.OpenAIChatClient.SendMessage(msg.Content, string(msg.Role))
		if err != nil {
			return err
		}
		openaiResponse := agent.OpenAIChatClient.GetLastMessage()
		serializedResponse := ChatAgentMessageFromOpenAIChatMessage(openaiResponse)
		processedResponse, err := agent.ProcessChatMessage(*serializedResponse)
		if err != nil {
			return err
		}
		agent.Messages = append(agent.Messages, processedResponse)
		agent.syncMessages()
		agent.serializeAllMessages()
		return serializedResponse
	}
}
