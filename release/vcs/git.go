package vcs

import (
	"github.com/alexejk/go-release-tools/config"
	"github.com/alexejk/go-release-tools/log"
	"gopkg.in/src-d/go-git.v4"
)

type GitHandler struct {
	workDir string

	isGit bool
	repo  *git.Repository
}

func NewGitHandler(workDir string) *GitHandler {
	g := &GitHandler{
		workDir: workDir,
	}

	repo, err := git.PlainOpen(g.workDir)
	if err != nil {
		log.Warnf("Working directory is not a valid git repository.")
	} else {
		g.repo = repo
	}

	return g
}

func (g *GitHandler) IsGitRepository() bool {
	return g.repo != nil
}

func (g *GitHandler) Commit(message string) error {

	worktree, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	versionFile := config.GetString(config.ProjectVersionFile)

	log.Infof("-> Git add: %s", versionFile)
	worktree.Add(versionFile)

	log.Infof("-> Git commit: %s", message)
	if _, err = worktree.Commit(message, &git.CommitOptions{}); err != nil {
		return err
	}

	return nil
}

func (g *GitHandler) Tag(name string) error {

	return nil
}
func (g *GitHandler) Push() error {

	return nil
}
