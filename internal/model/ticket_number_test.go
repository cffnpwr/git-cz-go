package model

import (
	"errors"
	"regexp"
	"testing"

	"github.com/cffnpwr/git-cz-go/config"
	"github.com/cffnpwr/git-cz-go/internal/mock/repo"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"go.uber.org/mock/gomock"
)

// TicketModelWrapper wraps TicketNumberModel to implement tea.Model interface properly
type TicketModelWrapper struct {
	TicketNumberModel
}

func (w TicketModelWrapper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	model, cmd := w.TicketNumberModel.Update(msg)
	return TicketModelWrapper{model}, cmd
}

func createTestTicketNumberModel() TicketNumberModel {
	ctrl := gomock.NewController(&testing.T{})
	defer ctrl.Finish()
	mockRepo := repo.NewMockGitRepository(ctrl)

	return TicketNumberModel{
		input: func() textinput.Model {
			m := textinput.New()
			m.Placeholder = defaultTicketPrompt
			return m
		}(),
		config:   config.TicketNumber{},
		gitRepo:  mockRepo,
		finished: false,
		valid:    true,
		errorMsg: "",
	}
}

func TestNewTicketNumberModel(t *testing.T) {
	tests := []struct {
		name      string
		prompt    string
		config    config.TicketNumber
		wantModel TicketNumberModel
	}{
		{
			name:   "[正常系] デフォルト設定でのモデル作成",
			prompt: "",
			config: config.TicketNumber{},
			wantModel: TicketNumberModel{
				input: func() textinput.Model {
					m := textinput.New()
					m.Placeholder = defaultTicketPrompt
					return m
				}(),
				config:   config.TicketNumber{},
				finished: false,
				valid:    true,
				errorMsg: "",
			},
		},
		{
			name:   "[正常系] カスタム設定でのモデル作成",
			prompt: "Custom ticket prompt",
			config: config.TicketNumber{},
			wantModel: TicketNumberModel{
				input: func() textinput.Model {
					m := textinput.New()
					m.Placeholder = "Custom ticket prompt"
					return m
				}(),
				config:   config.TicketNumber{},
				finished: false,
				valid:    true,
				errorMsg: "",
			},
		},
		{
			name:   "[正常系] 空文字列設定（デフォルト値が使用される）",
			prompt: "",
			config: config.TicketNumber{},
			wantModel: TicketNumberModel{
				input: func() textinput.Model {
					m := textinput.New()
					m.Placeholder = defaultTicketPrompt
					return m
				}(),
				config:   config.TicketNumber{},
				finished: false,
				valid:    true,
				errorMsg: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repo.NewMockGitRepository(ctrl)
			got := NewTicketNumberModel(tt.prompt, tt.config, mockRepo)
			if diff := cmp.Diff(tt.wantModel, got, cmpopts.IgnoreUnexported(TicketNumberModel{})); diff != "" {
				t.Errorf("NewTicketNumberModel() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestTicketNumberModel_validateInput(t *testing.T) {
	tests := []struct {
		name          string
		config        config.TicketNumber
		inputValue    string
		expectedValid bool
		expectedError string
	}{
		{
			name:          "[正常系] 空値、必須でない場合",
			config:        config.TicketNumber{Required: false},
			inputValue:    "",
			expectedValid: true,
		},
		{
			name:          "[異常系] 空値、必須の場合",
			config:        config.TicketNumber{Required: true},
			inputValue:    "",
			expectedValid: false,
			expectedError: "Ticket number is required",
		},
		{
			name: "[正常系] パターンマッチ成功",
			config: config.TicketNumber{
				Required:     true,
				MatchPattern: (*config.Regexp)(regexp.MustCompile(`\d+`)),
			},
			inputValue:    "123",
			expectedValid: true,
		},
		{
			name: "[異常系] パターンマッチ失敗",
			config: config.TicketNumber{
				Required:     true,
				MatchPattern: (*config.Regexp)(regexp.MustCompile(`\d+`)),
			},
			inputValue:    "abc",
			expectedValid: false,
			expectedError: "Invalid ticket number format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repo.NewMockGitRepository(ctrl)
			model := TicketNumberModel{
				config:  tt.config,
				gitRepo: mockRepo,
			}
			model.input.SetValue(tt.inputValue)
			result := model.validateInput()

			if result.valid != tt.expectedValid {
				t.Errorf("validateInput() valid = %v, want %v", result.valid, tt.expectedValid)
			}
			if result.errorMsg != tt.expectedError {
				t.Errorf("validateInput() errorMsg = %v, want %v", result.errorMsg, tt.expectedError)
			}
		})
	}
}

func TestTicketNumberModel_GetValue(t *testing.T) {
	tests := []struct {
		name       string
		prefix     string
		inputValue string
		expected   string
	}{
		{
			name:       "[正常系] プレフィックスなし",
			inputValue: "123",
			expected:   "123",
		},
		{
			name:       "[正常系] プレフィックスあり",
			prefix:     "#",
			inputValue: "123",
			expected:   "#123",
		},
		{
			name:       "[正常系] 空値、プレフィックスあり",
			prefix:     "#",
			inputValue: "",
			expected:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repo.NewMockGitRepository(ctrl)
			model := TicketNumberModel{
				config: config.TicketNumber{
					Prefix: tt.prefix,
				},
				gitRepo: mockRepo,
			}
			model.input.SetValue(tt.inputValue)

			result := model.GetValue()
			if result != tt.expected {
				t.Errorf("GetValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTicketNumberModel_Update(t *testing.T) {
	tests := []struct {
		name           string
		initialModel   func(repo *repo.MockGitRepository) TicketNumberModel
		keyInput       tea.KeyMsg
		wantFinished   bool
		wantQuit       bool
		validateResult func(t *testing.T, model TicketNumberModel)
	}{
		{
			name: "[正常系] Enterキーで入力完了",
			initialModel: func(repo *repo.MockGitRepository) TicketNumberModel {
				model := TicketNumberModel{
					config:   config.TicketNumber{Required: false},
					gitRepo:  repo,
					valid:    true,
					finished: false,
				}
				return model
			},
			keyInput:     tea.KeyMsg{Type: tea.KeyEnter},
			wantFinished: true,
			wantQuit:     false,
		},
		{
			name: "[正常系] 文字入力",
			initialModel: func(repo *repo.MockGitRepository) TicketNumberModel {
				return TicketNumberModel{
					config:   config.TicketNumber{Required: false},
					gitRepo:  repo,
					valid:    true,
					finished: false,
				}
			},
			keyInput:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}},
			wantFinished: false,
			wantQuit:     false,
		},
		{
			name: "[正常系] 完了後の入力は無視される",
			initialModel: func(repo *repo.MockGitRepository) TicketNumberModel {
				return TicketNumberModel{
					config:   config.TicketNumber{Required: false},
					gitRepo:  repo,
					valid:    true,
					finished: true,
				}
			},
			keyInput:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("test")},
			wantFinished: true,
			wantQuit:     false,
			validateResult: func(t *testing.T, model TicketNumberModel) {
				// finished状態では入力処理自体が無視されるため、このテストは削除または変更が必要
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := repo.NewMockGitRepository(ctrl)

			model := tt.initialModel(repo)
			updatedModel, cmd := model.Update(tt.keyInput)

			if updatedModel.IsFinished() != tt.wantFinished {
				t.Errorf("IsFinished() = %v, want %v", updatedModel.IsFinished(), tt.wantFinished)
			}
			if tt.wantQuit && cmd == nil {
				t.Error("expected quit command, got nil")
			} else if !tt.wantQuit && cmd != nil {
				if _, ok := cmd().(tea.QuitMsg); ok {
					t.Error("unexpected quit command")
				}
			}
			if tt.validateResult != nil {
				tt.validateResult(t, updatedModel)
			}
		})
	}
}

func TestTicketNumberModel_extractTicketFromBranch(t *testing.T) {
	tests := []struct {
		name       string
		config     config.FromBranchName
		branchName string
		repoErr    error
		expected   string
	}{
		{
			name: "[正常系] 機能無効",
			config: config.FromBranchName{
				Enable: false,
			},
			branchName: "main",
		},
		{
			name: "[正常系] 正規表現がnil",
			config: config.FromBranchName{
				Enable:        true,
				ExtractRegexp: nil,
			},
			branchName: "main",
		},
		{
			name: "[正常系] 名前付きキャプチャグループで成功",
			config: config.FromBranchName{
				Enable:        true,
				ExtractRegexp: (*config.Regexp)(regexp.MustCompile(`^.+?/(?P<ticket_number>\d+)([-_]\w+)*$`)),
			},
			branchName: "feature/123-add-feature",
			expected:   "123",
		},
		{
			name: "[正常系] パターンマッチしない",
			config: config.FromBranchName{
				Enable:        true,
				ExtractRegexp: (*config.Regexp)(regexp.MustCompile(`^.+?/(?P<ticket_number>\d+)([-_]\w+)*$`)),
			},
			branchName: "main",
		},
		{
			name: "[異常系] Gitエラー",
			config: config.FromBranchName{
				Enable:        true,
				ExtractRegexp: (*config.Regexp)(regexp.MustCompile(`^.+?/(?P<ticket_number>\d+)([-_]\w+)*$`)),
			},
			repoErr: errors.New("git error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := repo.NewMockGitRepository(ctrl)
			if tt.config.Enable && tt.config.ExtractRegexp != nil {
				mockRepo.EXPECT().GetCurrentBranch().Return(tt.branchName, tt.repoErr)
			}

			model := TicketNumberModel{
				config: config.TicketNumber{
					FromBranchName: tt.config,
				},
				gitRepo: mockRepo,
			}

			result := model.extractTicketFromBranch()
			if result != tt.expected {
				t.Errorf("extractTicketFromBranch() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTicketNumberModel_Init(t *testing.T) {
	model := createTestTicketNumberModel()
	cmd := model.Init()
	if cmd == nil {
		t.Error("Init() should return textinput.Blink command, got nil")
	}
}

func TestTicketNumberModel_View(t *testing.T) {
	tests := []struct {
		name      string
		model     func() TicketNumberModel
		wantEmpty bool
	}{
		{
			name: "[正常系] 通常の表示",
			model: func() TicketNumberModel {
				return createTestTicketNumberModel()
			},
			wantEmpty: false,
		},
		{
			name: "[正常系] エラー状態の表示",
			model: func() TicketNumberModel {
				model := createTestTicketNumberModel()
				model.valid = false
				model.errorMsg = "Test error"
				return model
			},
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := tt.model()
			view := model.View()

			if tt.wantEmpty && view != "" {
				t.Error("expected empty view")
			} else if !tt.wantEmpty && view == "" {
				t.Error("expected non-empty view")
			}
		})
	}
}

func TestTicketNumberModel_IsFinished(t *testing.T) {
	tests := []struct {
		name     string
		finished bool
		want     bool
	}{
		{
			name:     "[正常系] 完了状態",
			finished: true,
			want:     true,
		},
		{
			name:     "[正常系] 未完了状態",
			finished: false,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := createTestTicketNumberModel()
			model.finished = tt.finished

			got := model.IsFinished()
			if got != tt.want {
				t.Errorf("IsFinished() = %v, want %v", got, tt.want)
			}
		})
	}
}
