package main

import (
	"fmt"

	"github.com/alexejk/go-release/cmd"
	"github.com/sirupsen/logrus"
)

// AppVersion denotes version number of this application
var AppVersion = "x.y.z-dev"

// AppName denotes name of this application
var AppName = "go-release"

// AppBuild denotes git commit that this binary was built from
var AppBuild = "~unknown~"

// AppBuildDate denotes the date when this binary was built
var AppBuildDate = "~unknown~"

// AppBuildTime denotes the time when this binary was built
var AppBuildTime = "~unknown~"

func init() {
	logrus.SetFormatter(&PrefixedTextFormatter{
		DisableTimestamp:  true,
		CapitalizeMessage: true,
		LogLevelFormat:    "[%5s]",
	})
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {

	versionString := fmt.Sprintf("Version: %s (Build: %s, %s @ %s)",
		AppVersion,
		AppBuild,
		AppBuildDate,
		AppBuildTime)

	if err := cmd.NewReleaseCommand(AppName, versionString).Execute(); err != nil {
		logrus.Fatal(err.Error())
	}

}
