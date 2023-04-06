package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/CSXL/solus/ai"
	"github.com/CSXL/solus/config"
	"github.com/CSXL/solus/query"
	"github.com/CSXL/solus/query/search_clients"
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
	primary     lipgloss.Color
	secondary   lipgloss.Color
	specialText lipgloss.Color
	title       lipgloss.Color
}

type styleMap struct {
	primary     lipgloss.Style
	secondary   lipgloss.Style
	title       lipgloss.Style
	specialText lipgloss.Style
	body        lipgloss.Style
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
		primary:     lipgloss.Color("#f8f8f2"),
		secondary:   lipgloss.Color("#b581c7"),
		title:       lipgloss.Color("#b581c7"),
		specialText: lipgloss.Color("#FFD580"),
	}
	styles = styleMap{
		primary: lipgloss.NewStyle().
			Foreground(colors.primary),
		secondary: lipgloss.NewStyle().
			Foreground(colors.secondary),
		specialText: lipgloss.NewStyle().
			Foreground(colors.specialText),
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

type TUIConfig struct {
	SavedMessagesFile    string
	DiscoveryMessage     string
	APIKey               string // In environment variable OPENAI_API_KEY
	LoadMessagesFromFile bool
	Debug                bool
}

type model struct {
	ChatClient  *ai.ChatClient
	QueryClient *query.QueryBuilder
	screen      screen
	input       textinput.Model
	viewport    viewport.Model
	tui_config  TUIConfig
	err         error
}

func NewModel(tui_config TUIConfig, query_client *query.QueryBuilder) model {
	ti := textinput.New()
	ti.Prompt = ""
	ti.Placeholder = "Enter your message here..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 80
	return model{
		ChatClient:  ai.NewChatClient(tui_config.APIKey),
		input:       ti,
		viewport:    viewport.New(80, 20),
		tui_config:  tui_config,
		QueryClient: query_client,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
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
				cmd := tea.ClearScreen
				cmds = append(cmds, cmd)
			} else {
				m.input.Focus()
			}
		case key.Matches(msg, keybindings.Help):
			if !m.input.Focused() {
				return m, tea.Println("This is a temporary help message that will be replaced by a help view.")
			}
		case key.Matches(msg, keybindings.Enter):
			if m.input.Value() != "" {
				messageToSend := newTUIMessage("message", "user", m.input.Value())
				jsonMessage, err := (&messageToSend).ToJSON()
				if err != nil {
					m.err = err
					return m, nil
				}
				m.ChatClient.SendUserMessage(jsonMessage)
				m.input.SetValue("")
			}
		case key.Matches(msg, keybindings.Save):
			m.ChatClient.SaveMessages(m.tui_config.SavedMessagesFile)
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

type tuiMessage struct {
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content string `json:"content"`
}

func newTUIMessage(_type string, role string, content string) tuiMessage {
	return tuiMessage{
		Role:    role,
		Type:    _type,
		Content: content,
	}
}

func (tmsg *tuiMessage) GetType() string {
	return tmsg.Type
}

func (tmsg *tuiMessage) GetRole() string {
	return tmsg.Role
}

func (tmsg *tuiMessage) GetContent() string {
	return tmsg.Content
}

func (tmsg *tuiMessage) ToJSON() (string, error) {
	jsonMessage, err := json.Marshal(tmsg)
	if err != nil {
		return "", err
	}
	return string(jsonMessage), nil
}

func processMessage(msg ai.ChatMessage) (tuiMessage, error) {
	var tuiMsg tuiMessage
	if msg.GetRole() == "assistant" || msg.GetRole() == "user" {
		AIMessage, err := msg.ToAIMessage()
		if err != nil {
			return tuiMsg, err
		}
		tuiMsg = newTUIMessage(AIMessage.GetType(), msg.GetRole(), AIMessage.GetContent())
		return tuiMsg, nil
	}
	tuiMsg = newTUIMessage("message", msg.GetRole(), msg.GetContent())
	return tuiMsg, nil
}

func (m model) repromptAI(prompt string) error {
	currentMessages := m.ChatClient.GetMessages()
	messagesBeforeLast := currentMessages[:len(currentMessages)-1]
	err := m.ChatClient.SendSystemMessage(prompt)
	if err != nil {
		return err
	}
	newMessage := m.ChatClient.GetLastMessage()
	slicedMessages := append(messagesBeforeLast, newMessage)
	m.ChatClient.SetMessages(slicedMessages)
	return nil
}

func (m model) resolveLastMessage() error {
	lastMessage := m.ChatClient.GetLastMessage()
	if lastMessage.GetRole() == "assistant" {
		AIMessage, err := lastMessage.ToAIMessage()
		content := AIMessage.GetContent()
		// If the chat model doesn't obey the JSON format, Solus will reprompt the AI until it does.
		// The reprompt messages and the offending message are deleted from the chat history to avoid clutter.
		for content == "" {
			messageToReprompt := newTUIMessage("reprompt", "system", "Send the last message again wrapped in JSON.")
			jsonMessageToReprompt, err := (&messageToReprompt).ToJSON()
			if err != nil {
				return err
			}
			m.repromptAI(jsonMessageToReprompt)
			lastMessage = m.ChatClient.GetLastMessage()
			AIMessage, err = lastMessage.ToAIMessage()
			if err != nil {
				return err
			}
			content = AIMessage.GetContent()
		}
		if err != nil {
			return err
		}
		if AIMessage.IsQuery() {
			queryResults, err := m.QueryClient.SetType("search").SetQueryText(AIMessage.GetContent()).Execute().GetResults()
			if err != nil {
				return err
			}
			err = m.ChatClient.SendSystemMessage(queryResults)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m model) ChatView() string {
	var s string

	err := m.resolveLastMessage()
	if err != nil {
		m.err = err
		return m.err.Error()
	}

	for _, msg := range m.ChatClient.GetMessages() {
		if msg.GetRole() != "system" || m.tui_config.Debug {
			tuiMsg, _ := processMessage(msg)
			formattedMessage := m.formatMessage(tuiMsg)
			s += styles.secondary.Render(formattedMessage)
			s += "\n"
		}
	}

	s += styles.secondary.Render("[USER]: ")
	s += styles.primary.Render(m.input.View())

	if !m.input.Focused() {
		m.viewport.SetContent(s)
		return m.viewport.View()
	}

	return s
}

func (m model) formatMessage(tuiMsg tuiMessage) string {
	if tuiMsg.GetType() == "query" {
		return m.formatQueryMessage(tuiMsg)
	}

	return m.formatNonQueryMessage(tuiMsg)
}

func (m model) formatQueryMessage(tuiMsg tuiMessage) string {
	coloredQuery := styles.specialText.Render(strings.Trim(tuiMsg.GetContent(), " \n"))
	formatted_message := fmt.Sprintf("Searching: %s\n\n", coloredQuery)

	return formatted_message
}

func (m model) formatNonQueryMessage(tuiMsg tuiMessage) string {
	formatted_role := strings.ToUpper(tuiMsg.GetRole())
	markdown_renderer, _ := glamour.NewTermRenderer(glamour.WithAutoStyle())
	markdown_content, _ := markdown_renderer.Render(tuiMsg.GetContent())
	formatted_message := fmt.Sprintf("[%s]: %s", formatted_role, markdown_content)

	return formatted_message
}

func readTUIConfig() (TUIConfig, error) {
	config_reader := config.New()
	err := config_reader.Read("tui_config", ".")
	if err != nil {
		return TUIConfig{}, err
	}
	tui_config := TUIConfig{}
	tui_config.SavedMessagesFile = config_reader.Get("saved_messages_file").(string)
	tui_config.DiscoveryMessage = config_reader.Get("discovery_message").(string)
	tui_config.LoadMessagesFromFile = config_reader.Get("load_messages_from_file").(bool)
	tui_config.Debug = config_reader.Get("debug").(bool)
	return tui_config, nil
}

func loadTUIConfig() (TUIConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return TUIConfig{}, err
	}
	openai_api_key := os.Getenv("OPENAI_API_KEY")
	tui_config, err := readTUIConfig()
	if err != nil {
		return TUIConfig{}, err
	}
	tui_config.APIKey = openai_api_key
	return tui_config, nil
}

func prepareChatClient(config TUIConfig, chatClient *ai.ChatClient) error {
	if config.LoadMessagesFromFile {
		err := chatClient.LoadMessages(config.SavedMessagesFile)
		if err != nil {
			return err
		}
	} else {
		systemTUIMessage := newTUIMessage("system", "system", config.DiscoveryMessage)
		systemTUIMessageJSON, err := systemTUIMessage.ToJSON()
		if err != nil {
			return err
		}
		err = chatClient.SendSystemMessage(systemTUIMessageJSON)
		if err != nil {
			return err
		}
	}
	return nil
}

func loadSearchEngineConfig() (search_clients.SearchClientConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return search_clients.SearchClientConfig{}, err
	}
	search_engine_config := search_clients.SearchClientConfig{
		GoogleSearchAPIKey:   os.Getenv("GOOGLE_API_KEY"),
		GoogleSearchEngineID: os.Getenv("GOOGLE_PROGRAMMABLE_SEARCH_ENGINE_ID"),
	}
	return search_engine_config, nil
}

func Run() (tea.Model, error) {
	ctx := context.Background()
	tui_config, err := loadTUIConfig()
	if err != nil {
		return nil, err
	}
	search_engine_config, err := loadSearchEngineConfig()
	if err != nil {
		return nil, err
	}
	query_client := query.NewQuery(ctx, search_engine_config)
	m := NewModel(tui_config, query_client)
	err = prepareChatClient(tui_config, m.ChatClient)
	if err != nil {
		return nil, err
	}
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	return p.Run()
}
