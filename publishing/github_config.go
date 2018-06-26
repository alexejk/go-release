package publishing

import (
	"errors"
	"regexp"

	"github.com/alexejk/go-release/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
)

const defaultGitRemote = "origin"

type GitHubPublishConfig struct {
	Draft       bool
	ReleaseName string
	Artifacts   []string

	Changelog *GitHubChangeLogConfig

	Owner      string
	Repository string

	// Working directory
	workDir string
}

func NewGitHubPublishConfig(workDir string) *GitHubPublishConfig {

	cfg := &GitHubPublishConfig{
		workDir: workDir,
	}

	config.Unmarshal(githubConfigKey, cfg)
	cfg.ensureOwnerAndRepo()

	return cfg
}

func (c *GitHubPublishConfig) ensureOwnerAndRepo() error {

	if c.Owner == "" || c.Repository == "" {

		// Figure out from workdir origin
		repo, _ := git.PlainOpen(c.workDir)
		remote, _ := repo.Remote(defaultGitRemote)
		urls := remote.Config().URLs

		remoteUrl := urls[0]

		owner, repoName, err := c.parseRemoteUrl(remoteUrl)

		if err != nil {
			return err
		}

		c.Owner = owner
		c.Repository = repoName
	}

	log.Debugf("Owner: %s, Repo: %s", c.Owner, c.Repository)

	return nil
}

func (c *GitHubPublishConfig) parseRemoteUrl(url string) (string, string, error) {

	var owner, repo string

	remoteRegexp := regexp.MustCompile(`(.*)github\.com(:|/)([a-zA-Z0-9-_]+)/([a-zA-Z0-9-_]+)(\.git)?$`)
	matches := remoteRegexp.FindStringSubmatch(url)

	if len(matches) >= 5 {
		owner = matches[3]
		repo = matches[4]
	} else {
		return owner, repo, errors.New("non github remote")
	}

	return owner, repo, nil
}
