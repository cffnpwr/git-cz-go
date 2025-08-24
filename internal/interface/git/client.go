package git

//go:generate mockgen -source=client.go -destination=../../mock/git/client.go -package=git

type GitClient interface {
	PlainOpen(path string) (GitRepository, error)
}
