package model

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultFooterPrompt = "Enter footer ('word: content' or 'word #content')"
)

var (
	footerSubmitKey = key.NewBinding(
		key.WithKeys("alt+enter", "ctrl+enter"),
		key.WithHelp("⌘+Enter/Ctrl+Enter", "submit footer"),
	)

	footerInfoColor  = lipgloss.Color("#696969")
	footerErrorColor = lipgloss.Color("#ff0000")
	footerInfoStyle  = lipgloss.NewStyle().Foreground(footerInfoColor)
	footerErrorStyle = lipgloss.NewStyle().Foreground(footerErrorColor)
)

// footerStartPattern
// 空白以外で構成された1つのトークンから始まり、`: `あるいは` #`が続くパターンで始まる
var footerStartPattern = regexp.MustCompile(`^\S+((:\s)|(\s#))`)

type footerValidationResult struct {
	valid    bool
	errorMsg string
}

type FooterModel struct {
	textarea textarea.Model
	finished bool
	valid    bool
	errorMsg string
}

func NewFooterModel(prompt string) FooterModel {
	ta := textarea.New()

	// プロンプトメッセージ設定
	if prompt == "" {
		prompt = defaultFooterPrompt
	}
	ta.Placeholder = prompt
	ta.Focus()

	return FooterModel{
		textarea: ta,
	}
}

func (m FooterModel) GetPrompt() string {
	return m.textarea.Placeholder
}

func (m FooterModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m FooterModel) Update(msg tea.Msg) (FooterModel, tea.Cmd) {
	if m.finished {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, footerSubmitKey) && (m.valid || strings.TrimSpace(m.textarea.Value()) == "") {
			m.finished = true
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)

	res := m.validateInput()
	m.valid = res.valid
	m.errorMsg = res.errorMsg
	return m, cmd
}

func (m FooterModel) View() string {
	view := m.textarea.View()

	if m.errorMsg != "" {
		view += "\n" + footerErrorStyle.Render("✕ "+m.errorMsg)
	}

	if (m.valid || strings.TrimSpace(m.textarea.Value()) == "") && !m.finished {
		view += "\n" + footerInfoStyle.Render("Press ⌘+Enter/Ctrl+Enter to continue")
	}

	return view
}

func (m FooterModel) IsFinished() bool {
	return m.finished
}

func (m FooterModel) GetValue() string {
	return strings.TrimSpace(m.textarea.Value())
}

func (m FooterModel) validateInput() footerValidationResult {
	value := strings.TrimSpace(m.textarea.Value())

	lines := strings.Split(value, "\n")
	if len(lines) == 0 {
		return footerValidationResult{
			valid: true,
		}
	}

	// 次の行が現在の行の続きかどうか
	isContinue := false
	for _, l := range lines {
		// 空行は無効
		if l == "" {
			return footerValidationResult{
				valid:    false,
				errorMsg: "Footer cannot contain empty line",
			}
		}

		// 続きの行でなくてFooterの開始パターンに一致しない場合は無効
		if !isContinue && !footerStartPattern.MatchString(l) {
			return footerValidationResult{
				valid:    false,
				errorMsg: "Footer must start with 'word: ' or 'word # ' format",
			}
		}
		isContinue = true
	}
	return footerValidationResult{
		valid: true,
	}
}
