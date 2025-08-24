package confirm

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultPrompt = "Confirm"
)

var (
	fgColor = lipgloss.Color("#bb9af7")
	bgColor = lipgloss.Color("#1e1e2e")

	border = lipgloss.NormalBorder()

	style         = lipgloss.NewStyle().Background(bgColor).Padding(0, 2).Margin(1, 1)
	selectedStyle = style.Foreground(lipgloss.Color("#ffffff")).Background(fgColor).Bold(true)
)

type KeyMap struct {
	Toggle      key.Binding
	Affirmative key.Binding
	Negative    key.Binding
	Select      key.Binding
	Quit        key.Binding
}

var DefaultKeyMap = KeyMap{
	Toggle: key.NewBinding(
		key.WithKeys("tab", "l", "right", "shift+tab", "h", "left"),
		key.WithHelp("←/→", "toggle selection"),
	),
	Affirmative: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "affirmative selection"),
	),
	Negative: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "negative selection"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "confirm selection"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q", "esc"),
		key.WithHelp("ctrl+c/q/esc", "quit"),
	),
}

type Model struct {
	Prompt string // Question prompt

	value     bool // Selected value
	confirmed bool // Confirmation status

	keyMap KeyMap // Key map
}

func New() Model {
	return Model{
		Prompt:    defaultPrompt,
		value:     false,
		confirmed: false,
		keyMap:    DefaultKeyMap,
	}
}

func (m Model) SetKeyMap(km KeyMap) Model {
	m.keyMap = km
	return m
}

func (m Model) GetValue() bool {
	return m.value
}

func (m Model) IsConfirmed() bool {
	return m.confirmed
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		km := m.keyMap
		switch {
		case key.Matches(msg, km.Toggle):
			m.value = !m.value
		case key.Matches(msg, km.Affirmative):
			m.value = true
			m.confirmed = true
		case key.Matches(msg, km.Negative):
			m.value = false
			m.confirmed = true
		case key.Matches(msg, km.Select):
			m.confirmed = true
		case key.Matches(msg, km.Quit):
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	var aff, neg string
	if m.value {
		aff = selectedStyle.Render("Yes")
		neg = style.Render("No")
	} else {
		aff = style.Render("Yes")
		neg = selectedStyle.Render("No")
	}

	return m.Prompt + "\n" + lipgloss.JoinHorizontal(lipgloss.Top, aff, neg)
}
