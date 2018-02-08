package publishing

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/alexejk/go-release/config"
	"github.com/alexejk/go-release/shared"
	"github.com/google/go-github/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const githubConfigKey = "publishing.github"

type GitHubChangeLogConfig struct {
	File   string
	Format string
}

type GitHubPublisher struct {
	api *github.Client
}

func NewGitHubPublisher() PublishWorker {

	ctx := context.Background()

	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: os.Getenv("GO_RELEASE_GITHUB_TOKEN"),
		},
	)
	authClient := oauth2.NewClient(ctx, tokenSource)
	p := &GitHubPublisher{
		api: github.NewClient(authClient),
	}

	return p
}

func (p *GitHubPublisher) Configured() bool {
	return config.IsSet(githubConfigKey)
}

func (p *GitHubPublisher) Publish(workDir string, versionInfo *shared.VersionInformation) error {

	cfg := NewGitHubPublishConfig(workDir)

	releaseBody := p.getChangelog(cfg, versionInfo.ReleaseVersion.String())
	releaseName := versionInfo.InterpolateReleaseVersionInString(cfg.ReleaseName)

	release := &github.RepositoryRelease{
		Name:    &releaseName,
		Body:    &releaseBody,
		TagName: &versionInfo.GitTagName,
		Draft:   &cfg.Draft,
	}

	rel, resp, err := p.api.Repositories.CreateRelease(context.Background(), cfg.Owner, cfg.Repository, release)
	if err != nil {
		return err
	}

	log.Debugf("Rel: %v", rel)
	log.Debugf("Resp: %v", resp)

	log.Infof("Uploading release assets (%d)", len(cfg.Artifacts))

	for _, artifactPath := range cfg.Artifacts {

		log.Debugf("Processing artifact definition: '%s'", artifactPath)

		interpolatedPath := versionInfo.InterpolateReleaseVersionInString(artifactPath)
		absPath, err := filepath.Abs(path.Join(workDir, interpolatedPath))
		if err != nil {
			log.Debugf("Unable to find artifact '%s'. %s", absPath, err)
			continue
		}

		_, fileName := filepath.Split(absPath)

		file, err := os.Open(absPath)
		if err != nil {
			log.Debugf("Unable to read artifact '%s'. %s", absPath, err)
		}
		uploadOps := &github.UploadOptions{
			Name: fileName,
		}
		_, _, err = p.api.Repositories.UploadReleaseAsset(context.Background(), cfg.Owner, cfg.Repository, *rel.ID, uploadOps, file)
		if err != nil {
			log.Warnf("Failed to upload artifact '%s'. %s", absPath, err)
		}
		file.Close()
	}

	return err
}

func (p *GitHubPublisher) getChangelog(config *GitHubPublishConfig, version string) string {
	if config.Changelog == nil || config.Changelog.File == "" {
		return ""
	}

	changeLogFile := path.Join(config.workDir, config.Changelog.File)
	fileBytes, err := ioutil.ReadFile(changeLogFile)
	if err != nil {
		log.Warnf("Cannot open a changelog file '%s'. %s", changeLogFile, err.Error())
		return ""
	}

	changelog := NewGitHubChangeLog(fileBytes, config.Changelog)

	return changelog.GetForVersion(version)
}
