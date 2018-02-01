package release

import (
	"errors"
	"os"
	"time"

	"github.com/alexejk/go-release/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	gitconfig "gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type GitHandler struct {
	workDir string
	version *VersionHandler

	repo *git.Repository
	conf *GitConfiguration
}

const (
	defaultCommitMessageRelease     = "Release version: ${version}"
	defaultCommitMessageDevelopment = "Next development version"
)

func NewGitHandler(workDir string, versionHandler *VersionHandler) *GitHandler {

	g := &GitHandler{
		workDir: workDir,
		conf:    &GitConfiguration{},
		version: versionHandler,
	}

	repo, err := git.PlainOpen(g.workDir)
	if err != nil {
		log.Warnf("Working directory is not a valid git repository.")
	} else {
		g.repo = repo
	}

	if err = config.Unmarshal(config.ProjectGitRoot, g.conf); err != nil {
		log.Fatalf("Unable to read Git configuration for project. %s", err)
	}

	if g.conf == nil {
		log.Fatalf("Must have git configuration for project")
	}

	return g
}

func (g *GitHandler) IsGitRepository() bool {
	return g.repo != nil
}

func (g *GitHandler) ReleaseCommit() (string, error) {

	if g.conf.Commit == nil {
		return "", errors.New("commits are not enabled in config")
	}

	message := defaultCommitMessageRelease
	if g.conf.Commit.Format != nil && g.conf.Commit.Format.Release != "" {

		message = g.conf.Commit.Format.Release
	}

	message = g.version.InterpolateVersionInString(message)

	return g.commit(message)
}

func (g *GitHandler) ReleaseTag() error {

	// Tagging not enabled
	if g.conf.Tag == nil {
		return nil
	}

	currVersion, _ := g.version.GetVersion()
	tagName := currVersion.String()
	if g.conf.Tag.Format != "" {
		tagName = g.version.InterpolateVersionInString(g.conf.Tag.Format)
	}

	return g.tag(tagName)
}

func (g *GitHandler) NextDevelopmentCommit() (string, error) {

	if g.conf.Commit == nil {
		return "", errors.New("commits are not enabled in config")
	}

	message := defaultCommitMessageDevelopment
	if g.conf.Commit.Format != nil && g.conf.Commit.Format.Development != "" {

		message = g.conf.Commit.Format.Development
	}

	message = g.version.InterpolateVersionInString(message)

	return g.commit(message)
}

func (g *GitHandler) commit(message string) (string, error) {

	workTree, err := g.repo.Worktree()
	if err != nil {
		return "", err
	}

	log.Debugf("-> Git commit with message: %s", message)
	hash, err := workTree.Commit(message, &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  g.conf.Author.Name,
			Email: g.conf.Author.Email,
			When:  time.Now(),
		},
	})

	if err != nil {
		return "", err
	}

	return hash.String(), nil
}
func (g *GitHandler) tag(tagName string) error {

	ref, err := g.repo.Head()
	if err != nil {
		return err
	}

	refName := plumbing.ReferenceName("refs/tags/" + tagName)
	tag := plumbing.NewHashReference(refName, ref.Hash())

	return g.repo.Storer.SetReference(tag)
}
func (g *GitHandler) Push() error {

	// Push Commits
	if g.conf.Commit != nil && g.conf.Commit.Push {

		log.Debug("-> Git Pushing commits")
		if err := g.repo.Push(&git.PushOptions{
			Progress: os.Stdout,
		}); err != nil {
			return err
		}
	}

	// Push Tags
	if g.conf.Tag != nil && g.conf.Tag.Push {

		log.Debug("-> Git Pushing tags")

		rs := gitconfig.RefSpec("refs/tags/*:refs/tags/*")
		if err := g.repo.Push(&git.PushOptions{
			Progress: os.Stdout,
			RefSpecs: []gitconfig.RefSpec{rs},
		}); err != nil {
			return err
		}

	}

	return nil
}
