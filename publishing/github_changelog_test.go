package publishing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitHubChangeLog_GetForVersion(t *testing.T) {

	data := []struct {
		contents string
		format   string
		version  string
		expected string
	}{
		{
			contents: `
-- UNRELEASED --
WIP Stuff

--- Release 1.2.3 ---
This is my release notes
Something else

- fix 1
- fix 2

--- Release 1.2.2 ---
Not what I want
`,
			format:  "--- Release ${version} ---",
			version: "1.2.3",
			expected: `This is my release notes
Something else

- fix 1
- fix 2`,
		},
	}

	for _, tt := range data {

		conf := &GitHubChangeLogConfig{
			Format: tt.format,
		}

		instance := NewGitHubChangeLog([]byte(tt.contents), conf)
		result := instance.GetForVersion(tt.version)
		assert.Equal(t, tt.expected, result)
	}

}
