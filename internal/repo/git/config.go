package git

import (
	"errors"
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gopasspw/gitconfig"
)

type GitConfigReaderImpl struct {
	config *gitconfig.Configs
}

func NewGitConfigReader() *GitConfigReaderImpl {
	return &GitConfigReaderImpl{
		config: gitconfig.New(),
	}
}

func (g *GitConfigReaderImpl) LoadConfig(repoPath string) error {
	g.config.LoadAll(repoPath)
	return nil
}

func (g *GitConfigReaderImpl) GetUserName() (string, error) {
	name := g.config.Get("user.name")
	if name == "" {
		return "", errors.New("user.name is not configured")
	}
	return name, nil
}

func (g *GitConfigReaderImpl) GetUserEmail() (string, error) {
	email := g.config.Get("user.email")
	if email == "" {
		return "", errors.New("user.email is not configured")
	}
	return email, nil
}

func (g *GitConfigReaderImpl) CreateSignature() (*object.Signature, error) {
	name, err := g.GetUserName()
	if err != nil {
		return nil, err
	}

	email, err := g.GetUserEmail()
	if err != nil {
		return nil, err
	}

	return &object.Signature{
		Name:  name,
		Email: email,
		When:  time.Now(),
	}, nil
}
