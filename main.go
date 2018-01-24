package main

import (
	"github.com/alexejk/go-release-tools/cmd"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&PrefixedTextFormatter{
		DisableTimestamp:  true,
		CapitalizeMessage: true,
		LogLevelFormat:    "[%5s]",
	})
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {

	if err := cmd.NewReleaseCommand().Execute(); err != nil {

	}

}
