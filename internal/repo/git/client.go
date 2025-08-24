package git

import (
	gitIF "github.com/cffnpwr/git-cz-go/internal/interface/git"
	"github.com/go-git/go-git/v5"
)

type GitClient struct{}

func (c *GitClient) PlainOpen(path string) (gitIF.GitRepository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	return &GitRepository{repo: repo}, nil
}
