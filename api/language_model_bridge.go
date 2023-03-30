package api

import "encoding/json"

type AIMessage struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func NewAIMessage(_type string, content string) *AIMessage {
	return &AIMessage{
		Type:    _type,
		Content: content,
	}
}

func NewAIMessageFromJSONString(jsonString string) (*AIMessage, error) {
	var aiMessage AIMessage
	err := json.Unmarshal([]byte(jsonString), &aiMessage)
	if err != nil {
		return nil, err
	}
	return &aiMessage, nil
}

func (msg *AIMessage) ToJSONString() (string, error) {
	marshalled, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	return string(marshalled), nil
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
