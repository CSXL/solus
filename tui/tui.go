package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/CSXL/solus/api"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
)

type keyMap struct {
	ToggleFocus key.Binding
	Save        key.Binding
	Enter       key.Binding
	Quit        key.Binding
	Down        key.Binding
	Up          key.Binding
	Help        key.Binding
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
		ToggleFocus: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "Toggle focus")),
		Save:        key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "Save")),
		Enter:       key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Send message")),
		Quit:        key.NewBinding(key.WithKeys("ctrl+c", "esc"), key.WithHelp("ctrl+c", "Quit")),
		Help:        key.NewBinding(key.WithKeys("h", "?"), key.WithHelp("h", "Show help")),
		Up:          key.NewBinding(key.WithKeys("j", "up"), key.WithHelp("j/↑", "Scroll up")),
		Down:        key.NewBinding(key.WithKeys("k", "down"), key.WithHelp("k/↓", "Scroll down")),
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
	viewport   viewport.Model
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
		input:      ti,
		viewport:   viewport.New(80, 20),
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
		case key.Matches(msg, keybindings.ToggleFocus):
			if m.input.Focused() {
				m.input.Blur()
			} else {
				m.input.Focus()
			}
		case key.Matches(msg, keybindings.Help):
			if !m.input.Focused() {
				return m, tea.Println("This is a temporary help message that will be replaced by a help view.")
			}
		case key.Matches(msg, keybindings.Enter):
			if m.input.Value() != "" {
				m.ChatClient.SendUserMessage(m.input.Value())
				m.input.SetValue("")
			}
		case key.Matches(msg, keybindings.Save):
			m.ChatClient.SaveMessages("messages.json")
		case key.Matches(msg, keybindings.Down):
			m.viewport.YOffset++
			if m.viewport.ScrollPercent() >= 100 {
				m.viewport.GotoBottom()
			}
		case key.Matches(msg, keybindings.Up):
			m.viewport.YOffset--
			if m.viewport.ScrollPercent() <= 0 {
				m.viewport.GotoTop()
			}
		}
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd
	if m.input.Focused() {
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
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
		formatted_role := strings.ToUpper(msg.Role)
		markdown_renderer, _ := glamour.NewTermRenderer(glamour.WithAutoStyle())
		markdown_content, _ := markdown_renderer.Render(msg.Content)
		formatted_message := fmt.Sprintf("[%s]: %s", formatted_role, markdown_content)
		s += styles.secondary.Render(formatted_message)
		s += "\n"
	}
	s += styles.secondary.Render("[USER]: ")
	s += styles.primary.Render(m.input.View())
	if !m.input.Focused() {
		m.viewport.SetContent(s)
		return m.viewport.View()
	}
	return s
}

func Run() (tea.Model, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	openai_api_key := os.Getenv("OPENAI_API_KEY")
	m := NewModel(openai_api_key)
	debug := os.Getenv("DEBUG")
	if strings.ToUpper(debug) == "TRUE" {
		m.ChatClient.LoadMessages("messages.json")
	} else {
		// TODO: Put discovery message into YAML config
		DISCOVERY_MESSAGE := "You are Solus, a helpful AI coding assistant by CSX Labs (Computer Science Exploration Laboratories). Start by greeting the user:"
		err = m.ChatClient.SendSystemMessage(DISCOVERY_MESSAGE)
		if err != nil {
			return nil, err
		}
	}
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	return p.Run()
}
