package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChatMessage__ToAIMessage(t *testing.T) {
	msg := ChatMessage{
		Content: `{"type": "message", "content": "Hello World"}`,
		Role:    "user",
	}
	aimsg, err := msg.ToAIMessage()
	if err != nil {
		t.Error(err)
	}
	if aimsg.GetContent() != "Hello World" {
		t.Errorf("ToAIMessage() returned wrong content: %v", aimsg.GetContent())
	}
	if !aimsg.IsMessage() {
		t.Errorf("ToAIMessage() returned wrong type: %v", aimsg.GetType())
	}
}

func TestNewChatClient(t *testing.T) {
	client := NewChatClient("test")
	if client == nil {
		t.Error("NewChatClient() returned nil")
	}
}

func TestGetSetandClearMessages(t *testing.T) {
	client := NewChatClient("test")
	messages := []ChatMessage{
		{
			Content: "Hello World",
			Role:    "user",
		},
		{
			Content: "Hello World",
			Role:    "assistant",
		},
	}
	client.SetMessages(messages)
	if len(client.GetMessages()) != 2 {
		t.Errorf("GetMessages() returned wrong number of messages: %v", len(client.GetMessages()))
	}
	client.ClearMessages()
	if len(client.GetMessages()) != 0 {
		t.Errorf("GetMessages() returned wrong number of messages: %v", len(client.GetMessages()))
	}
}

func TestUnmarhsalMessages(t *testing.T) {
	client := NewChatClient("test")
	marshalledMessages := `[{"Content":"Hello World","Role":"user"},{"Content":"Hello World","Role":"assistant"}]`
	expectedMessages := []ChatMessage{
		{
			Content: "Hello World",
			Role:    "user",
		},
		{
			Content: "Hello World",
			Role:    "assistant",
		},
	}
	r := strings.NewReader(marshalledMessages)
	actualMessages, err := client.unmarshalMessages(r)
	if err != nil {
		t.Errorf("unmarshalMessages() returned error: %v", err)
	}
	if len(actualMessages) != len(expectedMessages) {
		t.Errorf("unmarshalMessages() returned wrong number of messages: %v", len(actualMessages))
	}
	for i := range actualMessages {
		if actualMessages[i].Content != expectedMessages[i].Content {
			t.Errorf("unmarshalMessages() returned wrong message content: %v", actualMessages[i].Content)
		}
		if actualMessages[i].GetRole() != expectedMessages[i].GetRole() {
			t.Errorf("unmarshalMessages() returned wrong message role: %v", actualMessages[i].GetRole())
		}
	}
}

func TestGetLastMessage(t *testing.T) {
	client := NewChatClient("test")
	messages := []ChatMessage{
		{
			Content: "Hello World",
			Role:    "user",
		},
		{
			Content: "Hello World",
			Role:    "assistant",
		},
	}
	client.SetMessages(messages)
	lastMessage := client.GetLastMessage()
	if lastMessage.Content != "Hello World" {
		t.Errorf("GetLastMessage() returned wrong message content: %v", lastMessage.Content)
	}
	if lastMessage.GetRole() != "assistant" {
		t.Errorf("GetLastMessage() returned wrong message role: %v", lastMessage.GetRole())
	}
}

func TestCreateChatCompletion(t *testing.T) {
	client := NewChatClient("test")
	// Fake the response from the OpenAI API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-Type", "application/json")
		// Response from https://api.openai.com/v1/chat/completions (prompt omitted)
		// Requested on 3/20/2023
		fakeResponse := `{"id":"chatcmpl-123","object":"chat.completion","created":1679367552,"model":"gpt-3.5-turbo-0301","usage":{"prompt_tokens":9,"completion_tokens":11,"total_tokens":20},"choices":[{"message":{"role":"assistant","content":"\n\nHi there! How may I assist you today?"},"finish_reason":"stop","index":0}]}`
		_, err := w.Write([]byte(fakeResponse))
		if err != nil {
			return
		}
	}))
	client.openAIClient = NewOpenAIWithBaseURL("test", ts.URL)
	messages := []ChatMessage{
		{
			Content: "Hello World",
			Role:    "user",
		},
		{
			Content: "Hello World",
			Role:    "assistant",
		},
	}
	newMessages, err := client.CreateChatCompletion(messages, "gpt-3.5-turbo")
	if err != nil {
		t.Errorf("CreateChatCompletion() returned error: %v", err)
	}
	lastMessageContent := newMessages[len(newMessages)-1].GetContent()
	if lastMessageContent != "\n\nHi there! How may I assist you today?" {
		t.Errorf("CreateChatCompletion() returned wrong completion: %v", lastMessageContent)
	}
}

func TestSendMessage(t *testing.T) {
	client := NewChatClient("test")
	// Fake the response from the OpenAI API
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Response from https://api.openai.com/v1/chat/completions (prompt omitted)
		// Requested on 3/20/2023
		fakeResponse := `{"id":"chatcmpl-123","object":"chat.completion","created":1679367552,"model":"gpt-3.5-turbo-0301","usage":{"prompt_tokens":9,"completion_tokens":11,"total_tokens":20},"choices":[{"message":{"role":"assistant","content":"\n\nHi there! How may I assist you today?"},"finish_reason":"stop","index":0}]}`
		_, err := w.Write([]byte(fakeResponse))
		if err != nil {
			return
		}
	}))
	client.openAIClient = NewOpenAIWithBaseURL("test", ts.URL)
	err := client.SendMessage("Hello World", "user")
	if err != nil {
		t.Errorf("SendMessage() returned error: %v", err)
	}
	lastMessageContent := client.GetLastMessage().GetContent()
	if lastMessageContent != "\n\nHi there! How may I assist you today?" {
		t.Errorf("SendMessage() returned wrong completion: %v", lastMessageContent)
	}
}
