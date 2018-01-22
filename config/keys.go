package config

const (
	ProjectVersionRoot          = "project.version"
	ProjectVersionFile          = ProjectVersionRoot + ".file"
	ProjectVersionProperty      = ProjectVersionRoot + ".property"
	ProjectVersionIncrementType = ProjectVersionRoot + ".increment"

	ProjectGitRoot               = "project.git"
	ProjectGitMessageRelease     = ProjectGitRoot + ".message.release"
	ProjectGitMessageDevelopment = ProjectGitRoot + ".message.development"
)
