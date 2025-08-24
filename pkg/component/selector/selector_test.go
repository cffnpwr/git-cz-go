package selector

import (
	"errors"
	"reflect"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type testItem struct {
	value string
}

func (t testItem) String() string {
	return t.value
}

// ModelWrapper wraps Model to implement tea.Model interface properly
type ModelWrapper struct {
	Model
}

func (w ModelWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := w.Model.Update(msg)
	return ModelWrapper{model}, cmd
}

func createTestItems(count int) []SelectItem {
	items := make([]SelectItem, count)
	for i := range count {
		items[i] = testItem{value: string(rune('A' + i))}
	}
	return items
}

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		items       []SelectItem
		displaySize int
		wantModel   Model
		wantError   error
	}{
		{
			name:        "[正常系] モデル初期化",
			items:       createTestItems(5),
			displaySize: 3,
			wantModel: Model{
				Prompt:       defaultPrompt,
				items:        createTestItems(5),
				displayRange: [2]int{0, 3},
				displaySize:  3,
				keyMap:       DefaultKeyMap,
			},
		},
		{
			name:        "[正常系] 単一アイテムでのモデル初期化",
			items:       createTestItems(1),
			displaySize: 3,
			wantModel: Model{
				Prompt:       defaultPrompt,
				items:        createTestItems(1),
				displayRange: [2]int{0, 1},
				displaySize:  3,
				keyMap:       DefaultKeyMap,
			},
		},
		{
			name:        "[正常系] 空のアイテムリストでのモデル初期化",
			items:       []SelectItem{},
			displaySize: 3,
			wantModel: Model{
				Prompt:       defaultPrompt,
				items:        []SelectItem{},
				displayRange: [2]int{0, 0},
				displaySize:  3,
				keyMap:       DefaultKeyMap,
			},
		},
		{
			name:        "[異常系] 無効な表示サイズ (0)",
			items:       createTestItems(3),
			displaySize: 0,
			wantError:   errors.New("invalid display size, must be positive integer"),
		},
		{
			name:        "[異常系] 無効な表示サイズ (負の値)",
			items:       createTestItems(3),
			displaySize: -1,
			wantError:   errors.New("invalid display size, must be positive integer"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := New(tt.items, tt.displaySize)

			if err != nil || tt.wantError != nil {
				if reflect.TypeOf(err) != reflect.TypeOf(tt.wantError) {
					t.Errorf("New() error type mismatch: got %T, want %T", err, tt.wantError)
				}
				return
			}

			if diff := cmp.Diff(tt.wantModel, model, cmpopts.IgnoreUnexported(Model{})); diff != "" {
				t.Errorf("New() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestInit(t *testing.T) {
	model, err := New(createTestItems(3), 3)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	cmd := model.Init()
	if cmd != nil {
		t.Errorf("Init() = %v, want nil", cmd)
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name         string
		items        []SelectItem
		displaySize  int
		initialModel func(Model) Model
		keyInput     tea.KeyMsg
		wantCursor   int
		wantSelected bool
		wantQuit     bool
	}{
		// Navigation tests
		{
			name:         "[正常系] 下方向に移動",
			items:        createTestItems(5),
			displaySize:  3,
			keyInput:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			wantCursor:   1,
			wantSelected: false,
			wantQuit:     false,
		},
		{
			name:        "[正常系] 上方向に移動",
			items:       createTestItems(5),
			displaySize: 3,
			initialModel: func(m Model) Model {
				m.cursor = 2
				return m
			},
			keyInput:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}},
			wantCursor:   1,
			wantSelected: false,
			wantQuit:     false,
		},
		{
			name:         "[正常系] 非循環モードでの上方向境界",
			items:        createTestItems(3),
			displaySize:  3,
			keyInput:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}},
			wantCursor:   0,
			wantSelected: false,
			wantQuit:     false,
		},
		{
			name:        "[正常系] 非循環モードでの下方向境界",
			items:       createTestItems(3),
			displaySize: 3,
			initialModel: func(m Model) Model {
				m.cursor = 2
				return m
			},
			keyInput:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			wantCursor:   2,
			wantSelected: false,
			wantQuit:     false,
		},
		{
			name:        "[正常系] 循環モードでの上方向ラップ",
			items:       createTestItems(3),
			displaySize: 3,
			initialModel: func(m Model) Model {
				return m.SetCyclic(true)
			},
			keyInput:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}},
			wantCursor:   2,
			wantSelected: false,
			wantQuit:     false,
		},
		{
			name:        "[正常系] 循環モードでの下方向ラップ",
			items:       createTestItems(3),
			displaySize: 3,
			initialModel: func(m Model) Model {
				m = m.SetCyclic(true)
				m.cursor = 2
				return m
			},
			keyInput:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			wantCursor:   0,
			wantSelected: false,
			wantQuit:     false,
		},
		{
			name:         "[正常系] 単一アイテムでのナビゲーション",
			items:        createTestItems(1),
			displaySize:  3,
			keyInput:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			wantCursor:   0,
			wantSelected: false,
			wantQuit:     false,
		},
		// Selection tests
		{
			name:         "[正常系] スペースキーで選択",
			items:        createTestItems(3),
			displaySize:  3,
			keyInput:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}},
			wantCursor:   0,
			wantSelected: true,
			wantQuit:     false,
		},
		{
			name:         "[正常系] エンターキーで選択",
			items:        createTestItems(3),
			displaySize:  3,
			keyInput:     tea.KeyMsg{Type: tea.KeyEnter},
			wantCursor:   0,
			wantSelected: true,
			wantQuit:     false,
		},
		// Selection state behavior
		{
			name:        "[正常系] 選択後は移動できない",
			items:       createTestItems(3),
			displaySize: 3,
			initialModel: func(m Model) Model {
				m.selected = true
				return m
			},
			keyInput:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			wantCursor:   0,
			wantSelected: true,
			wantQuit:     false,
		},
		// Quit tests
		{
			name:         "[正常系] エスケープキーで終了",
			items:        createTestItems(3),
			displaySize:  3,
			keyInput:     tea.KeyMsg{Type: tea.KeyEsc},
			wantCursor:   0,
			wantSelected: false,
			wantQuit:     true,
		},
		{
			name:         "[正常系] Ctrl+Cで終了",
			items:        createTestItems(3),
			displaySize:  3,
			keyInput:     tea.KeyMsg{Type: tea.KeyCtrlC},
			wantCursor:   0,
			wantSelected: false,
			wantQuit:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := New(tt.items, tt.displaySize)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			if tt.initialModel != nil {
				model = tt.initialModel(model)
			}

			updatedModel, cmd := model.Update(tt.keyInput)
			finalModel := updatedModel

			if finalModel.cursor != tt.wantCursor {
				t.Errorf("cursor = %d, want %d", finalModel.cursor, tt.wantCursor)
			}
			if finalModel.selected != tt.wantSelected {
				t.Errorf("selected = %t, want %t", finalModel.selected, tt.wantSelected)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := createTestItems(3)
			model, err := New(items, 3)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}

			tm := teatest.NewTestModel(t, ModelWrapper{model}, teatest.WithInitialTermSize(80, 24))
			tm.Send(tt.keyInput)
			tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
		})
	}
}

func TestView(t *testing.T) {
	tests := []struct {
		name                    string
		items                   []SelectItem
		showSelectedItem        bool
		selected                bool
		cursor                  int
		wantEmptyAfterSelection bool
	}{
		{
			name:                    "[正常系] 初期状態での非空表示",
			items:                   createTestItems(3),
			wantEmptyAfterSelection: false,
		},
		{
			name:                    "[正常系] 選択済みアイテムを非表示",
			items:                   createTestItems(3),
			showSelectedItem:        false,
			selected:                true,
			wantEmptyAfterSelection: true,
		},
		{
			name:                    "[正常系] 選択済みアイテムを表示",
			items:                   createTestItems(3),
			showSelectedItem:        true,
			selected:                true,
			cursor:                  1,
			wantEmptyAfterSelection: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := New(tt.items, 3)
			if err != nil {
				t.Fatalf("New() error = %v", err)
			}
			model = model.SetShowSelectedItem(tt.showSelectedItem)
			if tt.selected {
				model.selected = tt.selected
				model.cursor = tt.cursor
			}

			view := model.View()

			if tt.wantEmptyAfterSelection && view != "" {
				t.Errorf("expected empty view, got: %q", view)
			} else if !tt.wantEmptyAfterSelection && view == "" {
				t.Error("expected non-empty view")
			}

			if tt.selected && tt.showSelectedItem {
				expectedItem := tt.items[tt.cursor].String()
				if view != questionStyle.Render(model.Prompt+"> ")+expectedItem {
					t.Errorf("expected selected item in view, got: %q", view)
				}
			}
		})
	}
}
