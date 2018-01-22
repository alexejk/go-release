package release

import (
	"testing"

	"github.com/alexejk/go-release-tools/config"
	"github.com/stretchr/testify/assert"
)

func TestManager_getVersionFile(t *testing.T) {

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

		instance := &Manager{
			workDir: tt.workDir,
		}

		result := instance.getVersionFile()
		assert.Equal(t, tt.expected, result)
	}
}
