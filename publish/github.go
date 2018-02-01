package publish

import "github.com/alexejk/go-release/config"

type GitHubPublisher struct {
}

func (p *GitHubPublisher) Configured() bool {
	return config.IsSet("publish.github")
}

func (p *GitHubPublisher) Publish() error {
	return nil
}
