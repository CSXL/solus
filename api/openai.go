package api

import (
	"context"
	"encoding/json"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type ChatClient struct {
	API_KEY  string
	messages []Message
	client   *OpenAI
}

func NewChatClient(api_key string) *ChatClient {
	return &ChatClient{
		API_KEY:  api_key,
		messages: []Message{},
		client:   NewOpenAI(api_key),
	}
}

func (c *ChatClient) GetMessages() []Message {
	return c.messages
}

func (c *ChatClient) SetMessages(messages []Message) {
	c.messages = messages
}

func (c *ChatClient) ClearMessages() {
	c.messages = []Message{}
}

func (c *ChatClient) LoadMessages(filename string) error {
	messages := []Message{}
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
	c.messages = append(c.messages, Message{content, role})
	messages, err := c.client.CreateChatCompletion(c.messages, openai.GPT3Dot5Turbo)
	c.messages = messages
	return err
}

func (c *ChatClient) GetLastMessage() Message {
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

type OpenAI struct {
	API_KEY string
	client  *openai.Client
}

func NewOpenAI(api_key string) *OpenAI {
	return &OpenAI{
		API_KEY: api_key,
		client:  openai.NewClient(api_key),
	}
}

func (o *OpenAI) CreateChatCompletion(messages []Message, model string) ([]Message, error) {
	openaiMessages := []openai.ChatCompletionMessage{}
	for _, message := range messages {
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Content: message.Content,
			Role:    message.Role,
		})
	}
	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: openaiMessages,
		},
	)
	if err != nil {
		return nil, err
	}
	new_message := resp.Choices[0].Message
	messages = append(messages, Message{
		Content: new_message.Content,
		Role:    new_message.Role,
	})
	return messages, nil
}
