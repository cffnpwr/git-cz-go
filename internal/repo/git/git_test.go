package git

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	gitMock "github.com/cffnpwr/git-cz-go/internal/mock/git"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/google/go-cmp/cmp"
	"go.uber.org/mock/gomock"
)

func TestGetCurrentBranch(t *testing.T) {
	tests := []struct {
		name       string
		mockSetup  func(*gitMock.MockGitClient, *gitMock.MockGitRepository)
		wantBranch string
		wantError  error
	}{
		{
			name: "[正常系] 有効なリポジトリパスでの処理",
			mockSetup: func(mockClient *gitMock.MockGitClient, mockRepo *gitMock.MockGitRepository) {
				mockClient.EXPECT().PlainOpen("/test/path").Return(mockRepo, nil)
				ref := plumbing.NewBranchReferenceName("main")
				mockRepo.EXPECT().Head().Return(plumbing.NewReferenceFromStrings(ref.String(), "dummy-hash"), nil)
			},
			wantBranch: "main",
			wantError:  nil,
		},
		{
			name: "[異常系] 無効なリポジトリパスでのエラー",
			mockSetup: func(mockClient *gitMock.MockGitClient, mockRepo *gitMock.MockGitRepository) {
				mockClient.EXPECT().PlainOpen("/test/path").Return(nil, fmt.Errorf("failed to open repository"))
			},
			wantError: fmt.Errorf("failed to open repository: %w", fmt.Errorf("")),
		},
		{
			name: "[異常系] HEADの取得エラー",
			mockSetup: func(mockClient *gitMock.MockGitClient, mockRepo *gitMock.MockGitRepository) {
				mockClient.EXPECT().PlainOpen("/test/path").Return(mockRepo, nil)
				mockRepo.EXPECT().Head().Return(nil, fmt.Errorf("failed to get HEAD"))
			},
			wantError: fmt.Errorf("failed to get HEAD: %w", fmt.Errorf("")),
		},
		{
			name: "[異常系] HEADがブランチを指していない",
			mockSetup: func(mockClient *gitMock.MockGitClient, mockRepo *gitMock.MockGitRepository) {
				mockClient.EXPECT().PlainOpen("/test/path").Return(mockRepo, nil)
				ref := plumbing.NewHashReference("HEAD", plumbing.NewHash("dummy-hash"))
				mockRepo.EXPECT().Head().Return(ref, nil)
			},
			wantError: errors.New("HEAD is not pointing to a branch"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := gitMock.NewMockGitClient(ctrl)
			mockRepo := gitMock.NewMockGitRepository(ctrl)

			tt.mockSetup(mockClient, mockRepo)

			mockConfigReader := gitMock.NewMockGitConfigReader(ctrl)
			gitRepo := NewGitRepositoryWithClient("/test/path", mockClient, mockConfigReader)
			branch, err := gitRepo.GetCurrentBranch()

			if diff := cmp.Diff(tt.wantBranch, branch); diff != "" {
				t.Errorf("GetCurrentBranch() branch mismatch (-want +got):\n%s", diff)
			}
			if err != nil || tt.wantError != nil {
				if reflect.TypeOf(err) != reflect.TypeOf(tt.wantError) {
					t.Errorf("GetCurrentBranch() error type mismatch: got %T, want %T", err, tt.wantError)
				}
				return
			}
		})
	}
}

func TestCommit(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*gitMock.MockGitClient, *gitMock.MockGitRepository, *gitMock.MockGitWorktree, *gitMock.MockGitConfigReader)
		message   string
		wantError error
	}{
		{
			name: "[正常系] 正常なコミットメッセージでのコミット成功",
			mockSetup: func(mockClient *gitMock.MockGitClient, mockRepo *gitMock.MockGitRepository, mockWorktree *gitMock.MockGitWorktree, mockConfigReader *gitMock.MockGitConfigReader) {
				mockClient.EXPECT().PlainOpen("/test/path").Return(mockRepo, nil)
				mockRepo.EXPECT().Worktree().Return(mockWorktree, nil)
				mockConfigReader.EXPECT().LoadConfig("/test/path").Return(nil)
				mockConfigReader.EXPECT().CreateSignature().Return(nil, nil)
				mockWorktree.EXPECT().Commit(gomock.Any(), gomock.Any()).Return(plumbing.NewHash("dummy-hash"), nil)
			},
			message:   "feat: add new feature",
			wantError: nil,
		},
		{
			name: "[正常系] 空のコミットメッセージでのコミット成功",
			mockSetup: func(mockClient *gitMock.MockGitClient, mockRepo *gitMock.MockGitRepository, mockWorktree *gitMock.MockGitWorktree, mockConfigReader *gitMock.MockGitConfigReader) {
				mockClient.EXPECT().PlainOpen("/test/path").Return(mockRepo, nil)
				mockRepo.EXPECT().Worktree().Return(mockWorktree, nil)
				mockConfigReader.EXPECT().LoadConfig("/test/path").Return(nil)
				mockConfigReader.EXPECT().CreateSignature().Return(nil, nil)
				mockWorktree.EXPECT().Commit(gomock.Any(), gomock.Any()).Return(plumbing.NewHash("dummy-hash"), nil)
			},
			message:   "",
			wantError: nil,
		},
		{
			name: "[異常系] 無効なリポジトリパスでのエラー",
			mockSetup: func(mockClient *gitMock.MockGitClient, mockRepo *gitMock.MockGitRepository, mockWorktree *gitMock.MockGitWorktree, mockConfigReader *gitMock.MockGitConfigReader) {
				mockClient.EXPECT().PlainOpen("/test/path").Return(nil, fmt.Errorf("failed to open repository"))
			},
			message:   "test message",
			wantError: fmt.Errorf("failed to open repository: %w", fmt.Errorf("")),
		},
		{
			name: "[異常系] Worktree取得エラー",
			mockSetup: func(mockClient *gitMock.MockGitClient, mockRepo *gitMock.MockGitRepository, mockWorktree *gitMock.MockGitWorktree, mockConfigReader *gitMock.MockGitConfigReader) {
				mockClient.EXPECT().PlainOpen("/test/path").Return(mockRepo, nil)
				mockRepo.EXPECT().Worktree().Return(nil, fmt.Errorf("failed to get worktree"))
			},
			message:   "test message",
			wantError: fmt.Errorf("failed to get worktree: %w", fmt.Errorf("")),
		},
		{
			name: "[異常系] コミットエラー",
			mockSetup: func(mockClient *gitMock.MockGitClient, mockRepo *gitMock.MockGitRepository, mockWorktree *gitMock.MockGitWorktree, mockConfigReader *gitMock.MockGitConfigReader) {
				mockClient.EXPECT().PlainOpen("/test/path").Return(mockRepo, nil)
				mockRepo.EXPECT().Worktree().Return(mockWorktree, nil)
				mockConfigReader.EXPECT().LoadConfig("/test/path").Return(nil)
				mockConfigReader.EXPECT().CreateSignature().Return(nil, nil)
				mockWorktree.EXPECT().Commit(gomock.Any(), gomock.Any()).Return(plumbing.Hash{}, fmt.Errorf("failed to commit"))
			},
			message:   "test message",
			wantError: fmt.Errorf("failed to commit: %w", fmt.Errorf("")),
		},
		{
			name: "[異常系] gitconfig未設定でのエラー",
			mockSetup: func(mockClient *gitMock.MockGitClient, mockRepo *gitMock.MockGitRepository, mockWorktree *gitMock.MockGitWorktree, mockConfigReader *gitMock.MockGitConfigReader) {
				mockClient.EXPECT().PlainOpen("/test/path").Return(mockRepo, nil)
				mockRepo.EXPECT().Worktree().Return(mockWorktree, nil)
				mockConfigReader.EXPECT().LoadConfig("/test/path").Return(nil)
				mockConfigReader.EXPECT().CreateSignature().Return(nil, errors.New("user.name is not configured"))
			},
			message:   "test message",
			wantError: fmt.Errorf("failed to create signature: %w", errors.New("user.name is not configured")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := gitMock.NewMockGitClient(ctrl)
			mockRepo := gitMock.NewMockGitRepository(ctrl)
			mockWorktree := gitMock.NewMockGitWorktree(ctrl)
			mockConfigReader := gitMock.NewMockGitConfigReader(ctrl)

			tt.mockSetup(mockClient, mockRepo, mockWorktree, mockConfigReader)

			gitRepo := NewGitRepositoryWithClient("/test/path", mockClient, mockConfigReader)
			err := gitRepo.Commit(tt.message)
			if err != nil || tt.wantError != nil {
				if reflect.TypeOf(err) != reflect.TypeOf(tt.wantError) {
					t.Errorf("Commit() error type mismatch: got %T, want %T", err, tt.wantError)
				}
			}
		})
	}
}
