package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Quit  key.Binding
	Help  key.Binding
}

type styleMap struct {
	primary   lipgloss.Style
	secondary lipgloss.Style
}

var (
	keybindings = keyMap{
		Up:    key.NewBinding(key.WithKeys("up"), key.WithHelp("up", "Move up")),
		Down:  key.NewBinding(key.WithKeys("down"), key.WithHelp("down", "Move down")),
		Left:  key.NewBinding(key.WithKeys("left"), key.WithHelp("left", "Move left")),
		Right: key.NewBinding(key.WithKeys("right"), key.WithHelp("right", "Move right")),
		Quit:  key.NewBinding(key.WithKeys("ctrl+c", "esc"), key.WithHelp("ctrl+c", "Quit")),
		Help:  key.NewBinding(key.WithKeys("h", "?"), key.WithHelp("h", "Show help")),
	}
	styles = styleMap{
		primary:   lipgloss.NewStyle().Foreground(lipgloss.Color("#f8f8f2")),
		secondary: lipgloss.NewStyle().Foreground(lipgloss.Color("#4E0069")),
	}
)

type message struct {
	content string
	role    string
}

type model struct {
	messages []message
}

func NewModel() model {
	return model{
		messages: []message{},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybindings.Quit):
			return m, tea.Quit
		case key.Matches(msg, keybindings.Help):
			return m, tea.Println("This is a temporary help message that will be replaced by a help view.")
		}
	}
	return m, nil
}

func (m model) View() string {
	var s string
	s += styles.primary.Render("Hello, World!")
	s += "\n"
	s += styles.secondary.Render("Hola, tierra!")
	s += "\n"
	return s
}

func Run() (tea.Model, error) {
	p := tea.NewProgram(NewModel())
	return p.Run()
}
