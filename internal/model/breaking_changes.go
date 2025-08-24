package model

import (
	"github.com/cffnpwr/git-cz-go/pkg/component/confirm"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type BreakingStage string

const (
	BreakingStageConfirm  BreakingStage = "confirm"
	BreakingStageInput    BreakingStage = "input"
	BreakingStageFinished BreakingStage = "finished"
)

const (
	defaultBreakingConfirmPrompt = "Are there any breaking changes? (Y/n)"
	defaultBreakingMessagePrompt = "Describe breaking changes..."
)

type BreakingChangesModel struct {
	stage     BreakingStage
	confirm   confirm.Model
	textinput textinput.Model
}

func NewBreakingChangesModel(confirmPrompt, messagePrompt string) BreakingChangesModel {
	confirmModel := confirm.New()

	// 設定からプロンプトメッセージを取得
	prompt := defaultBreakingConfirmPrompt
	if confirmPrompt != "" {
		prompt = confirmPrompt
	}
	confirmModel.Prompt = prompt

	textinputModel := textinput.New()
	prompt = defaultBreakingMessagePrompt
	if messagePrompt != "" {
		prompt = messagePrompt
	}
	textinputModel.Prompt = prompt
	textinputModel.Focus()

	return BreakingChangesModel{
		stage:     BreakingStageConfirm,
		confirm:   confirmModel,
		textinput: textinputModel,
	}
}

func (m BreakingChangesModel) GetConfirmPrompt() string {
	return m.confirm.Prompt
}

func (m BreakingChangesModel) GetMessagePrompt() string {
	return m.textinput.Prompt
}

func (m BreakingChangesModel) IsFinished() bool {
	return m.stage == BreakingStageFinished
}

func (m BreakingChangesModel) GetValue() string {
	if m.confirm.IsConfirmed() && m.confirm.GetValue() {
		return m.textinput.Value()
	}
	return ""
}

func (m BreakingChangesModel) HasBreakingChanges() bool {
	return m.confirm.IsConfirmed() && m.confirm.GetValue()
}

func (m BreakingChangesModel) Init() tea.Cmd {
	return m.confirm.Init()
}

func (m BreakingChangesModel) Update(msg tea.Msg) (BreakingChangesModel, tea.Cmd) {
	switch m.stage {
	case BreakingStageConfirm:
		var cmd tea.Cmd
		m.confirm, cmd = m.confirm.Update(msg)
		if m.confirm.IsConfirmed() {
			if m.confirm.GetValue() {
				m.stage = BreakingStageInput
			} else {
				m.stage = BreakingStageFinished
			}
		}
		return m, cmd
	case BreakingStageInput:
		if msg, ok := msg.(tea.KeyMsg); ok {
			if key.Matches(msg, submitKey) {
				m.stage = BreakingStageFinished
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.textinput, cmd = m.textinput.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m BreakingChangesModel) View() string {
	if m.stage == BreakingStageConfirm {
		return m.confirm.View()
	}
	return m.confirm.View() + "\n" + m.textinput.View()
}
