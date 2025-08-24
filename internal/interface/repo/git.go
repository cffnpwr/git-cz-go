package repo

//go:generate mockgen -source=git.go -destination=../../mock/repo/git.go -package=repo

type GitRepository interface {
	GetCurrentBranch() (string, error)
	Commit(message string) error
}
