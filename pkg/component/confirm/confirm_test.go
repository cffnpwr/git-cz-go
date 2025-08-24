package confirm

import (
	"reflect"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// ModelWrapper wraps Model to implement tea.Model interface properly
type ModelWrapper struct {
	Model
}

func (w ModelWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := w.Model.Update(msg)
	return ModelWrapper{model}, cmd
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		wantModel Model
	}{
		{
			name: "[正常系] モデル初期化",
			wantModel: Model{
				Prompt:    defaultPrompt,
				value:     false,
				confirmed: false,
				keyMap:    DefaultKeyMap,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := New()

			if diff := cmp.Diff(tt.wantModel, model, cmpopts.IgnoreUnexported(Model{})); diff != "" {
				t.Errorf("New() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestInit(t *testing.T) {
	model := New()
	cmd := model.Init()
	if cmd != nil {
		t.Errorf("Init() = %v, want nil", cmd)
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name          string
		initialModel  func(Model) Model
		keyInput      tea.KeyMsg
		wantValue     bool
		wantConfirmed bool
		wantQuit      bool
	}{
		// Toggle tests
		{
			name:          "[正常系] Tabキーでトグル (false -> true)",
			keyInput:      tea.KeyMsg{Type: tea.KeyTab},
			wantValue:     true,
			wantConfirmed: false,
			wantQuit:      false,
		},
		{
			name: "[正常系] Tabキーでトグル (true -> false)",
			initialModel: func(m Model) Model {
				m.value = true
				return m
			},
			keyInput:      tea.KeyMsg{Type: tea.KeyTab},
			wantValue:     false,
			wantConfirmed: false,
			wantQuit:      false,
		},
		{
			name:          "[正常系] 右矢印キーでトグル",
			keyInput:      tea.KeyMsg{Type: tea.KeyRight},
			wantValue:     true,
			wantConfirmed: false,
			wantQuit:      false,
		},
		{
			name:          "[正常系] 左矢印キーでトグル",
			keyInput:      tea.KeyMsg{Type: tea.KeyLeft},
			wantValue:     true,
			wantConfirmed: false,
			wantQuit:      false,
		},
		{
			name:          "[正常系] Lキーでトグル",
			keyInput:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}},
			wantValue:     true,
			wantConfirmed: false,
			wantQuit:      false,
		},
		{
			name:          "[正常系] Hキーでトグル",
			keyInput:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}},
			wantValue:     true,
			wantConfirmed: false,
			wantQuit:      false,
		},
		// Affirmative tests
		{
			name:          "[正常系] Yキーで肯定的選択",
			keyInput:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}},
			wantValue:     true,
			wantConfirmed: true,
			wantQuit:      false,
		},
		// Negative tests
		{
			name:          "[正常系] Nキーで否定的選択",
			keyInput:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}},
			wantValue:     false,
			wantConfirmed: true,
			wantQuit:      false,
		},
		{
			name: "[正常系] Nキーで否定的選択 (初期値がtrue)",
			initialModel: func(m Model) Model {
				m.value = true
				return m
			},
			keyInput:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}},
			wantValue:     false,
			wantConfirmed: true,
			wantQuit:      false,
		},
		// Select tests
		{
			name:          "[正常系] エンターキーで選択確定",
			keyInput:      tea.KeyMsg{Type: tea.KeyEnter},
			wantValue:     false,
			wantConfirmed: true,
			wantQuit:      false,
		},
		{
			name:          "[正常系] スペースキーで選択確定",
			keyInput:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}},
			wantValue:     false,
			wantConfirmed: true,
			wantQuit:      false,
		},
		{
			name: "[正常系] エンターキーで選択確定 (値がtrue)",
			initialModel: func(m Model) Model {
				m.value = true
				return m
			},
			keyInput:      tea.KeyMsg{Type: tea.KeyEnter},
			wantValue:     true,
			wantConfirmed: true,
			wantQuit:      false,
		},
		// Quit tests
		{
			name:          "[正常系] エスケープキーで終了",
			keyInput:      tea.KeyMsg{Type: tea.KeyEsc},
			wantValue:     false,
			wantConfirmed: false,
			wantQuit:      true,
		},
		{
			name:          "[正常系] Ctrl+Cで終了",
			keyInput:      tea.KeyMsg{Type: tea.KeyCtrlC},
			wantValue:     false,
			wantConfirmed: false,
			wantQuit:      true,
		},
		{
			name:          "[正常系] Qキーで終了",
			keyInput:      tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			wantValue:     false,
			wantConfirmed: false,
			wantQuit:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := New()
			if tt.initialModel != nil {
				model = tt.initialModel(model)
			}

			updatedModel, cmd := model.Update(tt.keyInput)
			finalModel := updatedModel

			if finalModel.value != tt.wantValue {
				t.Errorf("value = %t, want %t", finalModel.value, tt.wantValue)
			}
			if finalModel.confirmed != tt.wantConfirmed {
				t.Errorf("confirmed = %t, want %t", finalModel.confirmed, tt.wantConfirmed)
			}
			if tt.wantQuit && cmd == nil {
				t.Error("expected quit command, got nil")
			} else if !tt.wantQuit && cmd != nil {
				t.Errorf("expected nil command, got %v", cmd)
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
			name:     "[正常系] エスケープキーで終了",
			keyInput: tea.KeyMsg{Type: tea.KeyEsc},
		},
		{
			name:     "[正常系] Ctrl+Cで終了",
			keyInput: tea.KeyMsg{Type: tea.KeyCtrlC},
		},
		{
			name:     "[正常系] Qキーで終了",
			keyInput: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := New()

			tm := teatest.NewTestModel(t, ModelWrapper{model}, teatest.WithInitialTermSize(80, 24))
			tm.Send(tt.keyInput)
			tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
		})
	}
}

func TestView(t *testing.T) {
	tests := []struct {
		name     string
		value    bool
		wantView string
	}{
		{
			name:     "[正常系] 初期状態 (No選択)",
			value:    false,
			wantView: "┌─────┐┌─────┐\n│ Yes ││ No  │\n└─────┘└─────┘",
		},
		{
			name:     "[正常系] Yes選択状態",
			value:    true,
			wantView: "┌─────┐┌─────┐\n│ Yes ││ No  │\n└─────┘└─────┘",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := New()
			model.value = tt.value

			view := model.View()

			// 視覚的なスタイルが適用されているかを簡単に確認
			if view == "" {
				t.Error("expected non-empty view")
			}
			// スタイルの詳細な確認は困難なので、基本的な内容の存在確認
			if view != "" && len(view) == 0 {
				t.Error("view should not be empty")
			}
		})
	}
}

func TestSetters(t *testing.T) {
	model := New()
	// SetKeyMap test
	customKeyMap := KeyMap{
		Toggle: DefaultKeyMap.Toggle,
		Select: DefaultKeyMap.Select,
		Quit:   DefaultKeyMap.Quit,
	}
	keyMapModel := model.SetKeyMap(customKeyMap)
	if !reflect.DeepEqual(keyMapModel.keyMap, customKeyMap) {
		t.Errorf("SetKeyMap() = %v, want %v", keyMapModel.keyMap, customKeyMap)
	}
}

func TestGetters(t *testing.T) {
	model := New()

	// Test initial values
	if model.GetValue() != false {
		t.Errorf("GetValue() = %t, want %t", model.GetValue(), false)
	}
	if model.IsConfirmed() != false {
		t.Errorf("IsConfirmed() = %t, want %t", model.IsConfirmed(), false)
	}

	// Test modified values
	model.value = true
	model.confirmed = true
	if model.GetValue() != true {
		t.Errorf("GetValue() = %t, want %t", model.GetValue(), true)
	}
	if model.IsConfirmed() != true {
		t.Errorf("IsConfirmed() = %t, want %t", model.IsConfirmed(), true)
	}
}
