package version

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/alexejk/go-release-tools/config"
	"github.com/coreos/go-semver/semver"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetVersion(t *testing.T) {

	data := []struct {
		contents string
		property string
		expected string
	}{
		{
			contents: `
APP_VERSION=0.3.8-alpha
`,
			property: "APP_VERSION",
			expected: "0.3.8-alpha",
		},
		{
			contents: `
SOME_VERSION = whatever
VERSION = 0.3.0
`,
			property: "VERSION",
			expected: "0.3.0",
		},

		{
			contents: `
prop=something
version  =32.0.0
`,
			property: "version",
			expected: "32.0.0",
		},
	}

	for _, tt := range data {

		config.Reset()
		config.Set(config.ProjectVersionProperty, tt.property)

		versionFile, err := ioutil.TempFile(os.TempDir(), "")
		assert.NoError(t, err)

		_, err = versionFile.WriteString(tt.contents)
		assert.NoError(t, err)

		instance := &Handler{
			versionFile:     versionFile.Name(),
			versionProperty: tt.property,
		}

		version, err := instance.GetVersion()
		assert.NoError(t, err)
		if version != nil {
			assert.Equal(t, tt.expected, version.String())
		}

		os.Remove(versionFile.Name())

	}
}

func TestHandler_ReleaseVersion(t *testing.T) {

	data := []struct {
		version  string
		expected string
	}{
		{
			version:  "1.2.3-dev",
			expected: "1.2.3",
		},
		{
			version:  "1.2.3-alpha",
			expected: "1.2.3",
		},

		{
			version:  "0.0.1-SNAPSHOT",
			expected: "0.0.1",
		},

		{
			version:  "0.0.1-SNAPSHOT+funnyBuild",
			expected: "0.0.1",
		},
	}

	instance := &Handler{}

	for _, tt := range data {

		input := semver.New(tt.version)
		result := instance.ReleaseVersion(input)

		assert.Equal(t, tt.expected, result.String())
		assert.NotEqual(t, input.String(), result.String()) // Sanity-check: did not modify other version
	}
}

func TestHandler_NextDevelopmentVersion(t *testing.T) {

	data := []struct {
		version   string
		increment string
		expected  string
	}{
		{
			version:   "1.2.3-dev",
			increment: "patch",
			expected:  "1.2.4-" + DevelopmentVersionPreRelease,
		},
		{
			version:   "1.2.3-alpha",
			increment: "major",
			expected:  "2.0.0-" + DevelopmentVersionPreRelease,
		},

		{
			version:   "0.0.4",
			increment: "minor",
			expected:  "0.1.0-" + DevelopmentVersionPreRelease,
		},
	}

	instance := &Handler{}

	for _, tt := range data {

		config.Reset()
		config.Set(config.ProjectVersionIncrementType, tt.increment)

		input := semver.New(tt.version)
		result := instance.NextDevelopmentVersion(input)

		assert.Equal(t, tt.expected, result.String())
		assert.NotEqual(t, input.String(), result.String()) // Sanity-check: did not modify other version
	}
}

func TestHandler_SetVersion(t *testing.T) {
	data := []struct {
		contents   string
		property   string
		newVersion string
		expected   string
	}{
		{
			contents: `
APP_VERSION=0.3.8-alpha

PROP=true
`,
			property:   "APP_VERSION",
			newVersion: "0.4.0",
			expected: `
APP_VERSION=0.4.0

PROP=true
`,
		},
		{
			contents: `
SOME_VERSION = whatever
VERSION = 0.3.0
VERSION2 =   0.3.1
`,
			property:   "VERSION",
			newVersion: "1.0.0-dev",
			expected: `
SOME_VERSION = whatever
VERSION = 1.0.0-dev
VERSION2 =   0.3.1
`,
		},

		{
			contents: `
prop=something
`,
			property:   "version",
			newVersion: "1.0.0",
			expected: `
prop=something
`,
		},
	}

	for _, tt := range data {

		config.Reset()
		config.Set(config.ProjectVersionProperty, tt.property)

		versionFile, err := ioutil.TempFile(os.TempDir(), "")
		assert.NoError(t, err)

		_, err = versionFile.WriteString(tt.contents)
		assert.NoError(t, err)

		instance := &Handler{
			versionFile:     versionFile.Name(),
			versionProperty: tt.property,
		}

		err = instance.SetVersion(semver.New(tt.newVersion))
		assert.NoError(t, err)

		savedFileBytes, err := ioutil.ReadFile(versionFile.Name())
		assert.NoError(t, err)

		assert.Equal(t, tt.expected, string(savedFileBytes))

		os.Remove(versionFile.Name())

	}
}

func Test_getVersionFile(t *testing.T) {

	data := []struct {
		workDir        string
		cfgVersionFile string
		expected       string
	}{
		{"", "version.properties", "version.properties"},
		{"./", "version.props", "version.props"},
		{"/opt/proj", "build.properties", "/opt/proj/build.properties"},
		{"../a/proj", "build.properties", "../a/proj/build.properties"},
		{"../a/proj", "build/version.properties", "../a/proj/build/version.properties"},
	}

	for _, tt := range data {

		config.Reset()
		config.Set(config.ProjectVersionFile, tt.cfgVersionFile)

		result := getVersionFile(tt.workDir)

		assert.Equal(t, tt.expected, result)
	}
}

func TestHandler_InterpolateVersionInString(t *testing.T) {

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

		instance := &Handler{
			versionStringCache: tt.version, // This is hacky, don't like it
		}
		result := instance.InterpolateVersionInString(tt.input)

		assert.Equal(t, tt.expected, result)
	}
}
