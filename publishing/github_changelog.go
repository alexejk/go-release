package publishing

import (
	"regexp"
	"strings"
)

const semverRegexp = `(\d+\.\d+\.\d+(-\w+(\+[\w\.\d]*)?)?)`

type GitHubChangeLog struct {
	versions map[string]string
}

func NewGitHubChangeLog(data []byte, conf *GitHubChangeLogConfig) *GitHubChangeLog {

	changeLogs := parseChangelog(data, conf.Format)

	cl := &GitHubChangeLog{
		versions: changeLogs,
	}

	return cl
}

func parseChangelog(data []byte, boundaryFormat string) map[string]string {

	formatWithoutVersionPlaceholder := strings.Replace(boundaryFormat, "${version}", semverRegexp, -1)

	// Build boundary Regexp
	boundaryRegexp := regexp.MustCompile(formatWithoutVersionPlaceholder)

	// Get positions of all boundaries
	boundaryIndices := boundaryRegexp.FindAllIndex(data, -1)

	totalSlices := len(boundaryIndices)
	totalBytes := len(data)

	changeLog := make(map[string]string, totalSlices)

	for idx, loc := range boundaryIndices {

		boundaryBytes := data[loc[0]:loc[1]]
		boundaryParts := boundaryRegexp.FindSubmatch(boundaryBytes)
		if len(boundaryParts) == 0 {
			// This is broken somehow
			continue
		}
		version := string(boundaryParts[1])

		sliceEnd := totalBytes
		if idx < totalSlices-1 {
			// Peak forward
			sliceEnd = boundaryIndices[idx+1][0]
		}

		chunkWithoutBoundary := string(data[loc[1]:sliceEnd])
		changeLog[version] = strings.TrimSpace(chunkWithoutBoundary)
	}

	return changeLog
}

func (c *GitHubChangeLog) GetForVersion(version string) string {

	return c.versions[version]
}
