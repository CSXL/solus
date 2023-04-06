package openai

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"github.com/sashabaranov/go-openai"
)

type AIMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (msg *AIMessage) GetContent() string {
	return msg.Content
}

func (msg *AIMessage) GetType() string {
	return msg.Type
}

func (msg *AIMessage) IsQuery() bool {
	return msg.Type == "query"
}

func (msg *AIMessage) IsMessage() bool {
	return msg.Type == "message"
}

type ChatMessage struct {
	Content string
	Role    string
}

func (msg ChatMessage) GetContent() string {
	return msg.Content
}

func (msg ChatMessage) GetRole() string {
	return msg.Role
}

func (msg *ChatMessage) ToAIMessage() (AIMessage, error) {
	marshalledContent := msg.GetContent()
	var unMarshalledContent AIMessage
	err := json.Unmarshal([]byte(marshalledContent), &unMarshalledContent)
	if err != nil {
		return unMarshalledContent, nil
	}
	return unMarshalledContent, err
}

type ChatClient struct {
	apiKey       string
	messages     []ChatMessage
	openAIClient *OpenAI
}

func NewChatClient(apiKey string) *ChatClient {
	return &ChatClient{
		apiKey:       apiKey,
		messages:     []ChatMessage{},
		openAIClient: NewOpenAI(apiKey),
	}
}

func (c *ChatClient) GetMessages() []ChatMessage {
	return c.messages
}

func (c *ChatClient) SetMessages(messages []ChatMessage) {
	c.messages = messages
}

func (c *ChatClient) ClearMessages() {
	c.messages = []ChatMessage{}
}

func (c *ChatClient) LoadMessages(filename string) error {
	handle, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer handle.Close()
	messages, err := c.unmarshalMessages(handle)
	c.messages = messages
	if err != nil {
		return err
	}
	return nil
}

func (c *ChatClient) unmarshalMessages(handle io.Reader) ([]ChatMessage, error) {
	var messages []ChatMessage
	err := json.NewDecoder(handle).Decode(&messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (c *ChatClient) SaveMessages(filename string) error {
	messages, _ := json.Marshal(c.GetMessages())
	err := os.WriteFile(filename, messages, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *ChatClient) AddMessage(role string, content string) {
	c.messages = append(c.messages, ChatMessage{content, role})
}

func (c *ChatClient) SendMessage(content string, role string) error {
	c.AddMessage(role, content)
	messages, err := c.CreateChatCompletion(c.messages, openai.GPT3Dot5Turbo)
	c.messages = messages
	return err
}

func (c *ChatClient) GetLastMessage() ChatMessage {
	if len(c.messages) == 0 {
		return ChatMessage{}
	}
	return c.messages[len(c.messages)-1]
}

func (c *ChatClient) SendUserMessage(msg string) error {
	err := c.SendMessage(msg, "user")
	return err
}

func (c *ChatClient) SendAssistantMessage(msg string) error {
	err := c.SendMessage(msg, "assistant")
	return err
}

func (c *ChatClient) SendSystemMessage(msg string) error {
	err := c.SendMessage(msg, "system")
	return err
}

func (c *ChatClient) CreateChatCompletion(messages []ChatMessage, model string) ([]ChatMessage, error) {
	var openaiMessages []openai.ChatCompletionMessage
	for _, message := range messages {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Content: message.GetContent(),
			Role:    message.GetRole(),
		})
	}
	resp, err := c.openAIClient.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    model,
			Messages: openaiMessages,
		},
	)
	if err != nil {
		return nil, err
	}
	newMessage := resp.Choices[0].Message
	messages = append(messages, ChatMessage{
		Content: newMessage.Content,
		Role:    newMessage.Role,
	})
	return messages, nil
}
