package publishing

import (
	"github.com/alexejk/go-release/shared"
)

type PublishWorker interface {
	Configured() bool
	Publish(workDir string, versionInfo *shared.VersionInformation) error
}

var publishers []PublishWorker

func init() {

	publishers = []PublishWorker{
		NewGitHubPublisher(),
		NewS3Publisher(),
	}
}

func ExecPublishers(workDir string, versionInfo *shared.VersionInformation) error {

	for _, pub := range publishers {
		if pub.Configured() {
			if err := pub.Publish(workDir, versionInfo); err != nil {
				return err
			}
		}
	}

	return nil
}
