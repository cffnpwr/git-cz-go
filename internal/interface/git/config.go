package git

import "github.com/go-git/go-git/v5/plumbing/object"

//go:generate mockgen -source=config.go -destination=../../mock/git/config.go -package=git

type GitConfigReader interface {
	LoadConfig(repoPath string) error
	GetUserName() (string, error)
	GetUserEmail() (string, error)
	CreateSignature() (*object.Signature, error)
}
