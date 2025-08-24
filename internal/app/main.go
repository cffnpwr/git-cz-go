package app

import (
	"os"

	"github.com/cffnpwr/git-cz-go/config"
	"github.com/cffnpwr/git-cz-go/internal/model"
	"github.com/cffnpwr/git-cz-go/internal/repo/git"
	tea "github.com/charmbracelet/bubbletea"
)

func Run(cfg *config.Config) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	gitRepo := git.NewGitRepository(wd)
	m, err := model.NewModel(cfg, gitRepo)
	if err != nil {
		return err
	}

	p := tea.NewProgram(&m)
	p.Run()
	return nil
}
