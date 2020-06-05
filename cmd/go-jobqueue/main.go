package main

import (
	"os"

	"github.com/pieterclaerhout/go-jobqueue/environ"
	"github.com/pieterclaerhout/go-log"
	"github.com/urfave/cli/v2"
)

func configure() {

	log.PrintColors = true
	log.PrintTimestamp = true

	environ.LoadFromPath()

}

func main() {

	configure()

	app := &cli.App{
		Commands: []*cli.Command{
			commandVersion,
			commandSetup,
			commandQueue,
			commandQueueMany,
			commandWorker,
		},
	}

	err := app.Run(os.Args)
	log.CheckError(err)

}
