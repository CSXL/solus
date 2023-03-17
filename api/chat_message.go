package api

import (
	"context"
	"encoding/json"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type ChatMessage struct {
	Content string
	Role    string
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
	messages := []ChatMessage{}
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &messages)
	if err != nil {
		return err
	}
	c.messages = messages
	return nil
}

func (c *ChatClient) SaveMessages(filename string) error {
	messages, _ := json.Marshal(c.GetMessages())
	err := os.WriteFile(filename, messages, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *ChatClient) SendMessage(content string, role string) error {
	c.messages = append(c.messages, ChatMessage{content, role})
	messages, err := c.CreateChatCompletion(c.messages, openai.GPT3Dot5Turbo)
	c.messages = messages
	return err
}

func (c *ChatClient) GetLastMessage() ChatMessage {
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
	openaiMessages := []openai.ChatCompletionMessage{}
	for _, message := range messages {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Content: message.Content,
			Role:    message.Role,
		})
	}
	resp, err := c.openAIClient.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
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
