package cmd

import (
	"github.com/alexejk/go-release-tools/config"
	"github.com/alexejk/go-release-tools/release"
	"github.com/spf13/cobra"
)

type releaseOptions struct {
	configFile string
	workDir    string
}

const defaultConfigFile = "release.yaml"

func NewReleaseCommand() *cobra.Command {

	o := releaseOptions{
		configFile: defaultConfigFile,
	}
	cmd := &cobra.Command{
		Use:  "release",
		RunE: o.Run,
	}

	cmd.Flags().StringVarP(&o.configFile, "config", "c", defaultConfigFile, "")
	cmd.Flags().StringVarP(&o.workDir, "dir", "d", "./", "")
	cmd.SilenceUsage = true

	return cmd
}

func (o *releaseOptions) Run(cmd *cobra.Command, args []string) error {

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
