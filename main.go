package main

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/bluemir/pw/cmd"
	"github.com/bluemir/pw/internal/buildinfo"
)

var AppName string
var Version string

func main() {
	buildinfo.AppName = AppName
	buildinfo.Version = Version

	if err := cmd.Run(); err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
}
