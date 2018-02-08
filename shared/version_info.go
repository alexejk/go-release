package shared

import (
	"strings"

	"github.com/coreos/go-semver/semver"
)

type VersionInformation struct {
	ReleaseVersion *semver.Version
	NextVersion    *semver.Version

	GitTagName string
}

func (v *VersionInformation) InterpolateReleaseVersionInString(input string) string {

	// TODO: nil-check

	res := strings.Replace(input, "${version}", v.ReleaseVersion.String(), -1)

	return res
}

func (v *VersionInformation) InterpolateNextVersionInString(input string) string {

	// TODO: nil-check

	res := strings.Replace(input, "${version}", v.NextVersion.String(), -1)

	return res
}
