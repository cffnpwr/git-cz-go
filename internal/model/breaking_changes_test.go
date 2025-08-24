package model

import (
	"testing"
	"time"

	"github.com/cffnpwr/git-cz-go/pkg/component/confirm"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// ModelWrapper wraps BreakingChangesModel to implement tea.Model interface properly
type ModelWrapper struct {
	BreakingChangesModel
}

func (w ModelWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := w.BreakingChangesModel.Update(msg)
	return ModelWrapper{model}, cmd
}

func createTestBreakingChangesModel() BreakingChangesModel {
	cm := confirm.New()
	cm.Prompt = "Are there any breaking changes?"
	return BreakingChangesModel{
		stage:     BreakingStageConfirm,
		confirm:   cm,
		textinput: textinput.New(),
	}
}

func TestNewBreakingChangesModel(t *testing.T) {
	tests := []struct {
		name          string
		confirmPrompt string
		messagePrompt string
		wantModel     BreakingChangesModel
	}{
		{
			name:          "[正常系] デフォルト設定でのモデル作成",
			confirmPrompt: "",
			messagePrompt: "",
			wantModel: BreakingChangesModel{
				stage: BreakingStageConfirm,
				confirm: func() confirm.Model {
					m := confirm.New()
					m.Prompt = defaultConfirmPrompt
					return m
				}(),
				textinput: func() textinput.Model {
					m := textinput.New()
					m.Placeholder = defaultBreakingMessagePrompt
					return m
				}(),
			},
		},
		{
			name:          "[正常系] カスタム設定でのモデル作成",
			confirmPrompt: "カスタム確認メッセージ",
			messagePrompt: "カスタムプレースホルダー",
			wantModel: BreakingChangesModel{
				stage: BreakingStageConfirm,
				confirm: func() confirm.Model {
					m := confirm.New()
					m.Prompt = "カスタム確認メッセージ"
					return m
				}(),
				textinput: func() textinput.Model {
					m := textinput.New()
					m.Placeholder = "カスタムプレースホルダー"
					return m
				}(),
			},
		},
		{
			name:          "[正常系] 部分的なカスタム設定（確認メッセージのみ）",
			confirmPrompt: "部分カスタム確認",
			messagePrompt: "",
			wantModel: BreakingChangesModel{
				stage: BreakingStageConfirm,
				confirm: func() confirm.Model {
					m := confirm.New()
					m.Prompt = "部分カスタム確認"
					return m
				}(),
				textinput: func() textinput.Model {
					m := textinput.New()
					m.Placeholder = defaultBreakingMessagePrompt
					return m
				}(),
			},
		},
		{
			name:          "[正常系] 部分的なカスタム設定（プレースホルダーのみ）",
			confirmPrompt: "",
			messagePrompt: "部分カスタムプレースホルダー",
			wantModel: BreakingChangesModel{
				stage: BreakingStageConfirm,
				confirm: func() confirm.Model {
					m := confirm.New()
					m.Prompt = defaultConfirmPrompt
					return m
				}(),
				textinput: func() textinput.Model {
					m := textinput.New()
					m.Placeholder = "部分カスタムプレースホルダー"
					return m
				}(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBreakingChangesModel(tt.confirmPrompt, tt.messagePrompt)
			if diff := cmp.Diff(tt.wantModel, got, cmpopts.IgnoreUnexported(BreakingChangesModel{})); diff != "" {
				t.Errorf("NewBreakingChangesModel() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestInit(t *testing.T) {
	model := createTestBreakingChangesModel()
	cmd := model.Init()
	if cmd != nil {
		t.Errorf("Init() = %v, want nil", cmd)
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name         string
		initialModel func(BreakingChangesModel) BreakingChangesModel
		keyInput     tea.KeyMsg
		wantStage    BreakingStage
		wantQuit     bool
	}{
		// Confirm stage tests
		{
			name:      "[正常系] Confirm段階でY選択（肯定的回答）",
			keyInput:  tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}},
			wantStage: BreakingStageInput,
			wantQuit:  false,
		},
		{
			name:      "[正常系] Confirm段階でN選択（否定的回答）",
			keyInput:  tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}},
			wantStage: BreakingStageFinished,
			wantQuit:  false,
		},
		{
			name:      "[正常系] Confirm段階でEnter選択（false値で確定）",
			keyInput:  tea.KeyMsg{Type: tea.KeyEnter},
			wantStage: BreakingStageFinished,
			wantQuit:  false,
		},
		{
			name: "[正常系] Confirm段階でスペース選択（true値に変更後確定）",
			initialModel: func(m BreakingChangesModel) BreakingChangesModel {
				m.confirm, _ = m.confirm.Update(tea.KeyMsg{Type: tea.KeyTab})
				return m
			},
			keyInput:  tea.KeyMsg{Type: tea.KeyEnter},
			wantStage: BreakingStageInput,
			wantQuit:  false,
		},
		{
			name:      "[正常系] Confirm段階でEsc選択（終了）",
			keyInput:  tea.KeyMsg{Type: tea.KeyEsc},
			wantStage: BreakingStageConfirm,
			wantQuit:  true,
		},
		// Input stage tests
		{
			name: "[正常系] Input段階でAlt+Enter選択",
			initialModel: func(m BreakingChangesModel) BreakingChangesModel {
				m.stage = BreakingStageInput
				return m
			},
			keyInput:  tea.KeyMsg{Type: tea.KeyEnter, Alt: true},
			wantStage: BreakingStageFinished,
			wantQuit:  false,
		},
		{
			name: "[正常系] Input段階で通常文字入力",
			initialModel: func(m BreakingChangesModel) BreakingChangesModel {
				m.stage = BreakingStageInput
				return m
			},
			keyInput:  tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}},
			wantStage: BreakingStageInput,
			wantQuit:  false,
		},
		// Finished stage tests
		{
			name: "[正常系] Finished段階では何も変更されない",
			initialModel: func(m BreakingChangesModel) BreakingChangesModel {
				m.stage = BreakingStageFinished
				return m
			},
			keyInput:  tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}},
			wantStage: BreakingStageFinished,
			wantQuit:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestBreakingChangesModel()
			if tt.initialModel != nil {
				model = tt.initialModel(model)
			}

			updatedModel, cmd := model.Update(tt.keyInput)

			if updatedModel.stage != tt.wantStage {
				t.Errorf("stage = %v, want %v", updatedModel.stage, tt.wantStage)
			}
			if tt.wantQuit && cmd == nil {
				t.Error("expected quit command, got nil")
			}
		})
	}
}

func TestView(t *testing.T) {
	tests := []struct {
		name  string
		stage BreakingStage
	}{
		{
			name:  "[正常系] Confirm段階の表示",
			stage: BreakingStageConfirm,
		},
		{
			name:  "[正常系] Input段階の表示",
			stage: BreakingStageInput,
		},
		{
			name:  "[正常系] Finished段階の表示",
			stage: BreakingStageFinished,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestBreakingChangesModel()
			model.stage = tt.stage

			view := model.View()

			if view == "" {
				t.Error("expected non-empty view")
			}
		})
	}
}

func TestQuit(t *testing.T) {
	tests := []struct {
		name     string
		keyInput tea.KeyMsg
	}{
		{
			name:     "[正常系] エスケープキー終了",
			keyInput: tea.KeyMsg{Type: tea.KeyEsc},
		},
		{
			name:     "[正常系] Ctrl+C終了",
			keyInput: tea.KeyMsg{Type: tea.KeyCtrlC},
		},
		{
			name:     "[正常系] Q終了",
			keyInput: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestBreakingChangesModel()

			tm := teatest.NewTestModel(t, ModelWrapper{model}, teatest.WithInitialTermSize(80, 24))
			tm.Send(tt.keyInput)
			tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
		})
	}
}

func TestStageConstants(t *testing.T) {
	if string(BreakingStageConfirm) != "confirm" {
		t.Errorf("BreakingStageConfirm = %s, want confirm", string(BreakingStageConfirm))
	}
	if string(BreakingStageInput) != "input" {
		t.Errorf("BreakingStageInput = %s, want input", string(BreakingStageInput))
	}
	if string(BreakingStageFinished) != "finished" {
		t.Errorf("BreakingStageFinished = %s, want finished", string(BreakingStageFinished))
	}
}
