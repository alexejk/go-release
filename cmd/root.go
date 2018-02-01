package cmd

import (
	"path"

	"github.com/alexejk/go-release/config"
	"github.com/alexejk/go-release/release"
	"github.com/spf13/cobra"
)

type releaseOptions struct {
	configFile string
	workDir    string
}

const defaultConfigFile = "release.yaml"

func NewReleaseCommand(appName, appVersion string) *cobra.Command {

	o := releaseOptions{
		configFile: defaultConfigFile,
	}
	cmd := &cobra.Command{
		Use:  appName,
		RunE: o.Run,
	}

	cmd.Flags().StringVarP(&o.configFile, "config", "c", "", "")
	cmd.Flags().StringVarP(&o.workDir, "dir", "d", "./", "")

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	cmd.Version = appVersion

	return cmd
}

func (o *releaseOptions) Run(cmd *cobra.Command, args []string) error {

	if o.configFile == "" {
		o.configFile = path.Join(o.workDir, defaultConfigFile)
	}

	if err := config.LoadConfig(o.configFile); err != nil {
		return err
	}

	manager := release.NewManager(o.workDir)

	if err := manager.PreRelease(); err != nil {

		return err
	}

	if err := manager.MakeRelease(); err != nil {

		return err
	}

	if err := manager.PostRelease(); err != nil {

		return err
	}

	return nil
}
