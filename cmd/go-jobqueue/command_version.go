package main

import (
	"github.com/pieterclaerhout/go-jobqueue/versioninfo"
	"github.com/pieterclaerhout/go-log"
	"github.com/urfave/cli/v2"
)

var commandVersion = &cli.Command{
	Name:  "version",
	Usage: "show the version info",
	Action: func(c *cli.Context) error {
		log.Info(versioninfo.ProjectName, versioninfo.Version, "("+versioninfo.Revision+")")
		return nil
	},
}
