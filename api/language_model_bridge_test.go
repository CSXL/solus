package api

import (
	"testing"
)

func TestNewAIMessage(t *testing.T) {
	aiMessage := NewAIMessage("message", "hello")
	if aiMessage.GetContent() != "hello" {
		t.Errorf("Expected content to be 'hello', got %s", aiMessage.GetContent())
	}
	if aiMessage.GetType() != "message" {
		t.Errorf("Expected type to be 'message', got %s", aiMessage.GetType())
	}
}

func TestAIMessageFromJSONString(t *testing.T) {
	jsonString := `{"type":"message","content":"hello"}`
	aiMessage, err := NewAIMessageFromJSONString(jsonString)
	if err != nil {
		t.Error(err)
	}
	if aiMessage.GetContent() != "hello" {
		t.Errorf("Expected content to be 'hello', got %s", aiMessage.GetContent())
	}
	if aiMessage.GetType() != "message" {
		t.Errorf("Expected type to be 'message', got %s", aiMessage.GetType())
	}
}

func TestAIMessageToJSONString(t *testing.T) {
	aiMessage := NewAIMessage("message", "hello")
	jsonString, err := aiMessage.ToJSONString()
	if err != nil {
		t.Error(err)
	}
	expected := `{"type":"message","content":"hello"}`
	if jsonString != expected {
		t.Errorf("Expected %s, got %s", expected, jsonString)
	}
}
