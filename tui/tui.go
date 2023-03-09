package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/CSXL/solus/api"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Enter key.Binding
	Quit  key.Binding
	Help  key.Binding
}

type colorMap struct {
	primary   lipgloss.Color
	secondary lipgloss.Color
	title     lipgloss.Color
}

type styleMap struct {
	primary   lipgloss.Style
	secondary lipgloss.Style
	title     lipgloss.Style
	body      lipgloss.Style
}

var (
	keybindings = keyMap{
		Up:    key.NewBinding(key.WithKeys("up"), key.WithHelp("up", "Move up")),
		Down:  key.NewBinding(key.WithKeys("down"), key.WithHelp("down", "Move down")),
		Left:  key.NewBinding(key.WithKeys("left"), key.WithHelp("left", "Move left")),
		Right: key.NewBinding(key.WithKeys("right"), key.WithHelp("right", "Move right")),
		Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Send message")),
		Quit:  key.NewBinding(key.WithKeys("ctrl+c", "esc"), key.WithHelp("ctrl+c", "Quit")),
		Help:  key.NewBinding(key.WithKeys("h", "?"), key.WithHelp("h", "Show help")),
	}
	colors = colorMap{
		primary:   lipgloss.Color("#f8f8f2"),
		secondary: lipgloss.Color("#b581c7"),
		title:     lipgloss.Color("#b581c7"),
	}
	styles = styleMap{
		primary: lipgloss.NewStyle().
			Foreground(colors.primary),
		secondary: lipgloss.NewStyle().
			Foreground(colors.secondary),
		title: lipgloss.NewStyle().
			Foreground(colors.title).
			Bold(true).
			Border(lipgloss.NormalBorder()).
			Width(100).
			Align(lipgloss.Center).
			BorderForeground(colors.secondary).
			Margin(1).
			MarginBottom(0),
		body: lipgloss.NewStyle().
			Foreground(colors.primary).
			Border(lipgloss.NormalBorder()).
			BorderForeground(colors.secondary).
			Width(100).
			Padding(1).
			Margin(1),
	}
)

type screen struct {
	width  int
	height int
}

type model struct {
	ChatClient *api.ChatClient
	screen     screen
	input      textinput.Model
	err        error
}

func NewModel(OPENAI_API_KEY string) model {
	ti := textinput.New()
	ti.Prompt = ""
	ti.Placeholder = "Enter your message here..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 80
	return model{
		ChatClient: api.NewChatClient(OPENAI_API_KEY),
		screen:     screen{},
		input:      ti,
		err:        nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.screen.width = msg.Width
		m.screen.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keybindings.Quit):
			return m, tea.Quit
		case key.Matches(msg, keybindings.Help):
			if !m.input.Focused() {
				return m, tea.Println("This is a temporary help message that will be replaced by a help view.")
			}
		case key.Matches(msg, keybindings.Enter):
			if m.input.Value() != "" {
				m.ChatClient.SendUserMessage(m.input.Value())
				m.input.SetValue("")
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var s string
	s += styles.title.Render("Solus")
	s += styles.body.Render(m.ChatView())
	return s
}

func (m model) ChatView() string {
	var s string
	for _, msg := range m.ChatClient.GetMessages() {
		formatted_message := fmt.Sprintf("[%s]: %s", strings.ToTitle(msg.Role), msg.Content)
		s += styles.secondary.Render(formatted_message)
		s += "\n"
	}
	s += styles.secondary.Render("[User]: ")
	s += styles.primary.Render(m.input.View())
	return s
}

func Run() (tea.Model, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	openai_api_key := os.Getenv("OPENAI_API_KEY")
	m := NewModel(openai_api_key)
	// TODO: Put discovery message into YAML config
	DISCOVERY_MESSAGE := "You are Solus, a helpful AI coding assistant by CSX Labs (Computer Science Exploration Laboratories). Start by greeting the user:"
	err = m.ChatClient.SendSystemMessage(DISCOVERY_MESSAGE)
	if err != nil {
		return nil, err
	}
	p := tea.NewProgram(m)
	return p.Run()
}
