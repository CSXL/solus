package main

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Quit  key.Binding
	Help  key.Binding
}

type message struct {
	content string
	role    string
}

type model struct {
	messages []message
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}
