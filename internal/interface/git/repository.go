package git

import "github.com/go-git/go-git/v5/plumbing"

//go:generate mockgen -source=repository.go -destination=../../mock/git/repository.go -package=git

type GitRepository interface {
	Head() (*plumbing.Reference, error)
	Worktree() (GitWorktree, error)
}
