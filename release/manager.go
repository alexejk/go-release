package release

import (
	"path/filepath"

	"github.com/alexejk/go-release/release/build"
	"github.com/alexejk/go-release/release/vcs"
	"github.com/alexejk/go-release/release/version"
	"github.com/coreos/go-semver/semver"
	log "github.com/sirupsen/logrus"
)

type Manager struct {
	workDir string

	currentVersion *semver.Version
	versionHandler *version.Handler
	vcs            *vcs.GitHandler
	builder        *build.Builder
}

func NewManager(workDir string) *Manager {

	absWorkDir, err := filepath.Abs(workDir)
	if err != nil {
		log.Fatalf("Working directory cannot be accessed. %s", err.Error())
	}

	m := &Manager{
		workDir: absWorkDir,
	}

	m.versionHandler = version.NewVersionHandler(m.workDir)
	m.vcs = vcs.NewGitHandler(m.workDir, m.versionHandler)
	m.builder = build.NewBuilder(m.workDir)

	return m
}

func (m *Manager) PreRelease() error {

	// Get Current version
	currentVersion, err := m.versionHandler.GetVersion()
	if err != nil {
		return err
	}

	log.Infof("Current project version is '%s'", currentVersion)

	releaseVersion := m.versionHandler.ReleaseVersion(currentVersion)

	// Overwrite with the release version
	log.Infof("Updating project version to '%s'", releaseVersion)
	if err := m.versionHandler.SetVersion(releaseVersion); err != nil {
		return err
	}

	// Update current state
	m.currentVersion = releaseVersion

	return nil
}

func (m *Manager) MakeRelease() error {

	// Commit release version
	if m.vcs.IsGitRepository() {
		log.Info("Creating a release commit")
		hash, err := m.vcs.ReleaseCommit()
		if err != nil {
			return err
		}

		log.Infof("Commit SHA: %s", hash)
	}

	// Build Project
	log.Info("Building project")
	if err := m.builder.Build(); err != nil {
		return err
	}

	// Create tag
	if m.vcs.IsGitRepository() {
		log.Info("Creating a release tag")
		if err := m.vcs.ReleaseTag(); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) PostRelease() error {

	// Set next development version
	nextVersion := m.versionHandler.NextDevelopmentVersion(m.currentVersion)

	// Overwrite with the development version
	log.Infof("Updating project version to '%s'", nextVersion)
	if err := m.versionHandler.SetVersion(nextVersion); err != nil {
		return err
	}
	m.currentVersion = nextVersion

	// Commit development version
	if m.vcs.IsGitRepository() {
		log.Info("Creating a development version commit")
		hash, err := m.vcs.NextDevelopmentCommit()
		if err != nil {
			return err
		}

		log.Infof("Commit SHA: %s", hash)
	}

	// Push
	if m.vcs.IsGitRepository() {
		log.Info("Pushing changes to remote")
		if err := m.vcs.Push(); err != nil {
			return err
		}
	}

	// Publish

	return nil
}

func (m *Manager) Revert() error {

	// What to do here?

	return nil
}
