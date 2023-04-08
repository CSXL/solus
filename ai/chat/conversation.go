package chat

import (
	"github.com/CSXL/solus/ai"
	"github.com/CSXL/solus/ai/agent"
)

type Conversation struct {
	chatAgent *agent.ChatAgent
	config    *ai.AIConfig
}

// NewConversation creates a new conversation with the given name and config.
// The underlying agent will be started automatically on the first message.
// If you want to start the agent preemptively, use the PreemptiveStart()
// function.
func NewConversation(name string, config *ai.AIConfig) *Conversation {
	return &Conversation{
		chatAgent: agent.NewChatAgent(name, config),
		config:    config,
	}
}

// Starts the agent underlying the conversation.
// As the conversation is started automatically on the first message, this
// function is only useful if you want to start the conversation preemptively
// due to computational constraints.
func (c *Conversation) PreemptiveStart() {
	c.chatAgent.Start()
}

// Close the conversation.
// This will close the underlying agent.
func (c *Conversation) Close() {
	c.chatAgent.Stop()
}

// Kills and deletes the conversation.
// This will kill and delete the underlying agent.
func (c *Conversation) Kill() {
	c.chatAgent.Kill()
	c.chatAgent.ResetMessages()
}

func (c *Conversation) startIfNotStarted() {
	if !c.chatAgent.IsRunning() {
		c.chatAgent.Start()
	}
}

func (c *Conversation) LoadFromFile(filename string) error {
	return c.chatAgent.OpenAIChatClient.LoadMessages(filename)
}

func (c *Conversation) SaveToFile(filename string) error {
	return c.chatAgent.OpenAIChatClient.SaveMessages(filename)
}

func (c *Conversation) GetMessages() []agent.ChatAgentMessage {
	return c.chatAgent.GetMessages()
}

func (c *Conversation) GetMessageCount() int {
	return len(c.chatAgent.GetMessages())
}

func (c *Conversation) GetLastMessage() agent.ChatAgentMessage {
	return c.chatAgent.GetLastMessage()
}

// Adds a message to the conversation without sending it to the agent.
func (c *Conversation) AddMessage(msg agent.ChatAgentMessage) {
	c.chatAgent.AddMessage(msg)
}

func (c *Conversation) SetMessages(messages []agent.ChatAgentMessage) {
	c.chatAgent.SetMessages(messages)
}

func (c *Conversation) ResetMessages() {
	c.chatAgent.ResetMessages()
}

func (c *Conversation) GetAgent() *agent.ChatAgent {
	return c.chatAgent
}

func (c *Conversation) GetConfig() *ai.AIConfig {
	return c.config
}

// Send a message to the conversation.
// The message will be sent to the agent and the agent will respond with a
// completion.
func (c *Conversation) SendUserMessage(msgContent string) (agent.ChatAgentMessage, error) {
	c.startIfNotStarted()
	agentMsg := agent.NewChatAgentMessage(agent.ChatAgentMessageTypeText, agent.ChatAgentMessageRoleUser, msgContent)
	aiResponse, err := c.chatAgent.SendChatMessage(*agentMsg)
	if aiResponse == nil {
		return agent.ChatAgentMessage{}, err
	}
	return *aiResponse, err
}

// Send a system message to the conversation.
// The message will be sent to the agent and the agent will respond with a
// completion.
func (c *Conversation) SendSystemMessage(msgContent string) (agent.ChatAgentMessage, error) {
	c.startIfNotStarted()
	agentMsg := agent.NewChatAgentMessage(agent.ChatAgentMessageTypeText, agent.ChatAgentMessageRoleSystem, msgContent)
	aiResponse, err := c.chatAgent.SendChatMessage(*agentMsg)
	if aiResponse == nil {
		return agent.ChatAgentMessage{}, err
	}
	return *aiResponse, err
}
