package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

//go:generate mockgen -source=worktree.go -destination=../../mock/git/worktree.go -package=git

type GitWorktree interface {
	Commit(msg string, opts *git.CommitOptions) (plumbing.Hash, error)
}
