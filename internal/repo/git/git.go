package git

import (
	"errors"
	"fmt"

	gitIF "github.com/cffnpwr/git-cz-go/internal/interface/git"
	"github.com/cffnpwr/git-cz-go/internal/interface/repo"
	"github.com/go-git/go-git/v5"
)

type gitRepositoryImpl struct {
	client       gitIF.GitClient
	configReader gitIF.GitConfigReader
	repoPath     string
}

func NewGitRepository(repoPath string) repo.GitRepository {
	return NewGitRepositoryWithClient(repoPath, &GitClient{}, NewGitConfigReader())
}

func NewGitRepositoryWithClient(repoPath string, client gitIF.GitClient, configReader gitIF.GitConfigReader) repo.GitRepository {
	return &gitRepositoryImpl{
		client:       client,
		configReader: configReader,
		repoPath:     repoPath,
	}
}

func (r *gitRepositoryImpl) GetCurrentBranch() (string, error) {
	repo, err := r.client.PlainOpen(r.repoPath)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	if !head.Name().IsBranch() {
		return "", errors.New("HEAD is not pointing to a branch")
	}

	branchName := head.Name().Short()
	if branchName == "" || branchName == "HEAD" {
		return "", errors.New("invalid branch name")
	}
	return branchName, nil
}

func (r *gitRepositoryImpl) Commit(message string) error {
	// Open the repository
	repo, err := r.client.PlainOpen(r.repoPath)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the worktree
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Load git configuration and create signature
	err = r.configReader.LoadConfig(r.repoPath)
	if err != nil {
		return fmt.Errorf("failed to load git config: %w", err)
	}

	signature, err := r.configReader.CreateSignature()
	if err != nil {
		return fmt.Errorf("failed to create signature: %w", err)
	}

	// Create the commit
	_, err = worktree.Commit(message, &git.CommitOptions{
		Author:    signature,
		Committer: signature,
	})
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}
