package main

import (
	"github.com/pieterclaerhout/go-jobqueue"
	"github.com/urfave/cli/v2"
)

var commandSetup = &cli.Command{
	Name:  "setup",
	Usage: "creates the database tables",
	Action: func(c *cli.Context) error {

		r, err := jobqueue.DefaultRepository()
		if err != nil {
			return err
		}

		return r.Setup()

	},
}
