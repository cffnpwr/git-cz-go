package model

import (
	"regexp"
	"strings"

	"github.com/cffnpwr/git-cz-go/config"
	"github.com/cffnpwr/git-cz-go/internal/interface/repo"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultTicketPrompt = "Enter ticket number"
)

var (
	infoColor  = lipgloss.Color("#696969")
	errorColor = lipgloss.Color("#ff0000")
	infoStyle  = lipgloss.NewStyle().Foreground(infoColor)
	errorStyle = lipgloss.NewStyle().Foreground(errorColor)
)

type tnValidationResult struct {
	valid    bool
	errorMsg string
}

type TicketNumberModel struct {
	input    textinput.Model
	config   config.TicketNumber
	gitRepo  repo.GitRepository
	finished bool
	valid    bool
	errorMsg string
}

func NewTicketNumberModel(tnPrompt string, tnCfg config.TicketNumber, gitRepo repo.GitRepository) TicketNumberModel {
	input := textinput.New()

	// プロンプトメッセージ設定
	prompt := defaultTicketPrompt
	if tnPrompt != "" {
		prompt = tnPrompt
	}
	input.Prompt = prompt + "> "
	input.Focus()

	model := TicketNumberModel{
		input:   input,
		config:  tnCfg,
		gitRepo: gitRepo,
	}

	// ブランチ名から自動抽出
	if tnCfg.FromBranchName.Enable {
		autoValue := model.extractTicketFromBranch()
		if autoValue != "" {
			model.input.SetValue(autoValue)
		}
	}

	return model
}

func (m TicketNumberModel) GetPrompt() string {
	return m.input.Prompt
}

func (m TicketNumberModel) GetValue() string {
	value := strings.TrimSpace(m.input.Value())
	if value != "" && m.config.Prefix != "" {
		return m.config.Prefix + value
	}
	return value
}

func (m TicketNumberModel) IsFinished() bool {
	return m.finished
}

func (m TicketNumberModel) Focus() {
	m.input.Focus()
}

func (m TicketNumberModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TicketNumberModel) Update(msg tea.Msg) (TicketNumberModel, tea.Cmd) {
	if m.finished {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, enterKey) && m.valid {
			m.finished = true
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	res := m.validateInput()
	m.valid = res.valid
	m.errorMsg = res.errorMsg
	return m, cmd
}

func (m TicketNumberModel) View() string {
	view := m.input.View()

	if m.errorMsg != "" {
		view += "\n" + errorStyle.Render("✕ "+m.errorMsg)
	}

	if m.valid && !m.finished {
		view += "\n" + infoStyle.Render("Press Enter to continue")
	}

	return view
}

func (m TicketNumberModel) validateInput() tnValidationResult {
	value := strings.TrimSpace(m.input.Value())

	if !m.config.Required {
		return tnValidationResult{
			valid: true,
		}
	}
	if m.config.MatchPattern == nil && value == "" {
		return tnValidationResult{
			valid:    false,
			errorMsg: "Ticket number is required",
		}
	}

	re := (*regexp.Regexp)(m.config.MatchPattern)
	if !re.MatchString(value) {
		return tnValidationResult{
			valid:    false,
			errorMsg: "Invalid ticket number format",
		}
	}
	return tnValidationResult{
		valid: true,
	}
}

func (m TicketNumberModel) extractTicketFromBranch() string {
	if !m.config.FromBranchName.Enable || m.config.FromBranchName.ExtractRegexp == nil {
		return ""
	}

	// Repository経由でブランチ名を取得
	branchName, err := m.gitRepo.GetCurrentBranch()
	if err != nil {
		return ""
	}

	// 正規表現で抽出
	re := (*regexp.Regexp)(m.config.FromBranchName.ExtractRegexp)
	matches := re.FindStringSubmatch(branchName)
	if len(matches) > 1 {
		// 名前付きキャプチャグループから抽出
		// ticket_numberが存在しない場合は空文字を返す
		for i, name := range re.SubexpNames() {
			if name == "ticket_number" && i < len(matches) {
				return matches[i]
			}
		}
	}

	return ""
}
