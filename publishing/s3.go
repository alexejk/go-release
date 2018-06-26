package publishing

import (
	"github.com/alexejk/go-release/config"
	"github.com/alexejk/go-release/shared"
)

const s3uploadConfigKey = config.PublishingRoot + "s3upload"

type S3Publisher struct {
}

func NewS3Publisher() *S3Publisher {

	p := &S3Publisher{}

	return p
}

func (p *S3Publisher) Configured() bool {
	return config.IsSet(s3uploadConfigKey)
}

func (p *S3Publisher) Publish(workDir string, versionInfo *shared.VersionInformation) error {

	//cfg := NewS3PublishConfig(workDir)
	//
	//// Create S3 connection
	//sess := session.Must(session.NewSessionWithOptions(session.Options{
	//	Config: aws.Config{
	//		CredentialsChainVerboseErrors: aws.Bool(true),
	//	},
	//	SharedConfigState:       session.SharedConfigEnable,
	//	AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	//}))
	//
	//artifacts := GetPublishableArtifacts(cfg.Artifacts, workDir, versionInfo)
	//for _, artifact := range artifacts {
	//
	//}

	return nil
}
