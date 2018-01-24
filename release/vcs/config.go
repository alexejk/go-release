package vcs

type GitConfiguration struct {
	Tag    *GitConfigurationTag
	Commit *GitConfigurationCommit

	Author *struct {
		Name  string
		Email string
	}
}

type GitConfigurationTag struct {
	Push   bool
	Format string
}

type GitConfigurationCommit struct {
	Push   bool
	Format *struct {
		Release     string
		Development string
	}
}

type GitAuthor struct {
	Name  string
	Email string
}
