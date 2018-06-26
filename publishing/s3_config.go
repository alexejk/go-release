package publishing

import "github.com/alexejk/go-release/config"

type S3PublishConfig struct {
	Bucket    string
	Prefix    string
	Region    string
	Artifacts []string

	workDir string
}

func NewS3PublishConfig(workDir string) *S3PublishConfig {

	cfg := &S3PublishConfig{
		workDir: workDir,
	}

	config.Unmarshal(s3uploadConfigKey, cfg)
	//	cfg.ensureOwnerAndRepo()

	return cfg
}
