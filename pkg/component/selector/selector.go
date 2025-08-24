package selector

import (
	"errors"
	"fmt"
	"math"

	"github.com/cffnpwr/git-cz-go/internal/util"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultPrompt = "Select"
	basePadding   = 0
)

var (
	fgColor = lipgloss.Color("#bb9af7")

	style             = lipgloss.NewStyle().PaddingTop(2).PaddingBottom(2).PaddingLeft(2)
	questionStyle     = lipgloss.NewStyle().Bold(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(basePadding + 2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(basePadding).Foreground(fgColor)
)

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
	Quit   key.Binding
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "mode up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "select item"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "esc"),
		key.WithHelp("Ctrl + C/Esc", "quit"),
	),
}

type SelectItem interface {
	fmt.Stringer
}

type Model struct {
	Prompt string // Question prompt

	showSelectedItem bool // Flag for show selected item
	selected         bool // Flag for selected
	cyclic           bool // Flag for circular view

	items        []SelectItem // Selectable items
	cursor       int          // Cursor position
	displayRange [2]int       // Item display range
	displaySize  int          // Item display size

	keyMap KeyMap // key map
}

func New(items []SelectItem, displaySize int) (Model, error) {
	if displaySize < 1 {
		return Model{}, errors.New("invalid display size, must be positive integer")
	}
	return Model{
		Prompt:           defaultPrompt,
		showSelectedItem: false,
		selected:         false,
		cyclic:           false,
		items:            items,
		cursor:           0,
		displayRange:     [2]int{0, displaySize},
		displaySize:      displaySize,
		keyMap:           DefaultKeyMap,
	}, nil
}

func (m Model) SetShowSelectedItem(b bool) Model {
	m.showSelectedItem = b
	return m
}

func (m Model) SetCyclic(b bool) Model {
	m.cyclic = b
	return m
}

func (m Model) SetKeyMap(km KeyMap) Model {
	m.keyMap = km
	return m
}

func (m Model) IsSelected() bool {
	return m.selected
}

func (m Model) GetSelectedItem() SelectItem {
	if m.selected {
		return m.items[m.cursor]
	}
	return nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		km := m.keyMap
		mod := util.GenMod(len(m.items))
		// middle of displayed items
		middle := mod(int(math.Floor(float64(m.displaySize)/2)) + m.displayRange[0])
		switch {
		case key.Matches(msg, km.Up):
			if m.selected {
				return m, nil
			}
			if m.cyclic {
				m.cursor = mod(m.cursor - 1)

				m.displayRange[0] = mod(m.displayRange[0] - 1)
				m.displayRange[1] = mod(m.displayRange[1] - 1)
			} else {
				if m.displayRange[0] != 0 && m.cursor <= middle {
					m.displayRange[0] -= 1
					m.displayRange[1] -= 1
				}
				if m.cursor != 0 {
					m.cursor -= 1
				}
			}
		case key.Matches(msg, km.Down):
			if m.selected {
				return m, nil
			}
			if m.cyclic {
				m.cursor = mod(m.cursor + 1)

				if m.cursor >= middle {
					m.displayRange[0] = mod(m.displayRange[0] + 1)
					m.displayRange[1] = mod(m.displayRange[1] + 1)
				}
			} else {
				if m.displayRange[1] != len(m.items) && m.cursor >= middle {
					m.displayRange[0] += 1
					m.displayRange[1] += 1
				}
				if m.cursor != len(m.items)-1 {
					m.cursor += 1
				}
			}
		case key.Matches(msg, km.Select):
			m.selected = true
		case key.Matches(msg, km.Quit):
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	mod := util.GenMod(len(m.items))

	qStr := m.Prompt + "> "
	qStr = questionStyle.Render(qStr)
	if m.selected {
		if !m.showSelectedItem {
			return ""
		}

		selected := m.items[m.cursor]
		return qStr + selected.String()
	}

	start := m.displayRange[0]
	end := m.displayRange[1]
	displayItems := []SelectItem{}
	if start >= end {
		displayItems = append(displayItems, m.items[start:]...)
		displayItems = append(displayItems, m.items[:end]...)
	} else {
		displayItems = m.items[start:end]
	}

	var selectStr string
	for index, i := range displayItems {
		itemStr := itemStyle.Render(i.String())
		if m.cursor == mod(index+start) {
			itemStr = selectedItemStyle.Render("> " + i.String())
		}
		selectStr += itemStr + "\n"
	}

	return qStr + style.Render(selectStr)
}
