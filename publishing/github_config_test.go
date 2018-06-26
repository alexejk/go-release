package publishing

import "testing"

func TestGitHubPublishConfig_parseRemoteUrl(t *testing.T) {

	testData := []struct {
		input       string
		outputOwner string
		outputRepo  string
		error       bool
	}{
		{"ssh://git@github.com/OwnerUser/target-repo", "OwnerUser", "target-repo", false},
		{"git@github.com:OwnerUser/target-repo.git", "OwnerUser", "target-repo", false},
		{"git@customgit.com:OwnerUser/target-repo.git", "", "", true},
	}

	cfg := &GitHubPublishConfig{
		workDir: ".",
	}
	for _, tt := range testData {
		o, r, err := cfg.parseRemoteUrl(tt.input)

		if (err != nil) != tt.error {
			t.Error("error state and expectation of error did not match", err)
		}

		if tt.outputOwner != o {
			t.Error("parsed owner and expectation did not match")
		}

		if tt.outputRepo != r {
			t.Error("parsed repository and expectation did not match")
		}
	}
}
