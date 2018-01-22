package release

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/alexejk/go-release-tools/config"
	"github.com/alexejk/go-release-tools/log"
	"github.com/alexejk/go-release-tools/release/vcs"
	"github.com/alexejk/go-release-tools/release/version"
	"github.com/coreos/go-semver/semver"
)

type Manager struct {
	workDir string

	currentVersion *semver.Version
	versionHandler *version.Handler
	vcsHandler     *vcs.GitHandler
}

func NewManager(workDir string) *Manager {

	absWorkDir, err := filepath.Abs(workDir)
	if err != nil {
		log.Fatalf("Working directory cannot be accessed. %s", err.Error())
	}

	m := &Manager{
		workDir: absWorkDir,
	}

	m.versionHandler = version.NewVersionHandler(m.getVersionFile())
	m.vcsHandler = vcs.NewGitHandler(m.workDir)

	return m
}

func (m *Manager) PreRelease() error {

	// Get Current version
	if err := m.readCurrentVersion(); err != nil {
		return err
	}

	log.Infof("Current project version is '%s'", m.currentVersion)

	releaseVersion := m.versionHandler.ReleaseVersion(m.currentVersion)

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
	if m.vcsHandler.IsGitRepository() {
		log.Info("Creating a release commit")

		msg := m.replaceVersionPlaceholder(config.GetString(config.ProjectGitMessageRelease))
		if err := m.vcsHandler.Commit(msg); err != nil {
			return err
		}
	}

	// Build Project
	log.Info("Building project")

	// Create tag
	if m.vcsHandler.IsGitRepository() {
		log.Info("Creating a release tag")
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

	// Push

	// Publish

	return nil
}

func (m *Manager) Revert() error {

	// What to do here?

	return nil
}

func (m *Manager) readCurrentVersion() error {

	projectVersion, err := m.versionHandler.GetVersion()
	m.currentVersion = projectVersion

	return err
}

func (m *Manager) replaceVersionPlaceholder(input string) string {

	res := strings.Replace(input, "${version}", m.currentVersion.String(), -1)

	log.Infof("Initial string: %s", input)
	log.Infof("Resulting string: %s", res)

	return res
}

func (m *Manager) getVersionFile() string {

	return path.Join(m.workDir, config.GetString(config.ProjectVersionFile))
}
