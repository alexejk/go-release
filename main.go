package main

import "github.com/alexejk/go-release-tools/cmd"

func main() {

	if err := cmd.NewReleaseCommand().Execute(); err != nil {

	}

}
