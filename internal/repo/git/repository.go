package git

import (
	gitIF "github.com/cffnpwr/git-cz-go/internal/interface/git"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

type GitRepository struct {
	repo *git.Repository
}

func (r *GitRepository) Head() (*plumbing.Reference, error) {
	return r.repo.Head()
}

func (r *GitRepository) Worktree() (gitIF.GitWorktree, error) {
	return r.repo.Worktree()
}
