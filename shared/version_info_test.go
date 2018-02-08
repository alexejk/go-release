package shared

import (
	"testing"

	"github.com/coreos/go-semver/semver"
	"github.com/stretchr/testify/assert"
)

func TestVersionInformation_InterpolateReleaseVersionInString(t *testing.T) {

	data := []struct {
		input    string
		version  string
		expected string
	}{
		{"Release ${version}.", "1.2.3", "Release 1.2.3."},
		{"Release ${version", "1.2.3", "Release ${version"},
		{"Release $${version}}", "1.2.3", "Release $1.2.3}"},
	}

	for _, tt := range data {

		instance := &VersionInformation{
			ReleaseVersion: semver.New(tt.version),
		}
		result := instance.InterpolateReleaseVersionInString(tt.input)

		assert.Equal(t, tt.expected, result)
	}
}

func TestVersionInformation_InterpolateNextVersionInString(t *testing.T) {

	data := []struct {
		input    string
		version  string
		expected string
	}{
		{"Next Release ${version}.", "1.2.3", "Next Release 1.2.3."},
		{"Next Release ${version", "1.2.3", "Next Release ${version"},
		{"Next Release $${version}}", "1.2.3", "Next Release $1.2.3}"},
	}

	for _, tt := range data {

		instance := &VersionInformation{
			NextVersion: semver.New(tt.version),
		}
		result := instance.InterpolateNextVersionInString(tt.input)

		assert.Equal(t, tt.expected, result)
	}
}
