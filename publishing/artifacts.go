package publishing

import (
	"os"
	"path"
	"path/filepath"

	"github.com/alexejk/go-release/shared"
	log "github.com/sirupsen/logrus"
)

type PublishableArtifact struct {
	Path string
	Name string
	Size int64
}

func GetPublishableArtifacts(paths []string, workDir string, versionInfo *shared.VersionInformation) []*PublishableArtifact {

	var artifacts []*PublishableArtifact

	for _, p := range paths {
		log.Debugf("Processing artifact definition: '%s'", p)

		interpolatedPath := versionInfo.InterpolateReleaseVersionInString(p)
		absPath, err := filepath.Abs(path.Join(workDir, interpolatedPath))
		if err != nil {
			log.Errorf("Unable to find artifact '%s'. %s", absPath, err)
			continue
		}

		file, err := os.Stat(absPath)
		if err != nil {
			log.Errorf("Unable to read artifact '%s'. %s", absPath, err)
			continue
		}

		if file.IsDir() {
			log.Errorf("Unable to use directory as publishable artifact: '%s'", absPath)
			continue
		}

		_, fileName := filepath.Split(absPath)
		artifact := &PublishableArtifact{
			Path: absPath,
			Name: fileName,
			Size: file.Size(),
		}

		artifacts = append(artifacts, artifact)
	}

	return artifacts
}
