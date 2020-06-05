package main

import (
	"os"

	"github.com/pieterclaerhout/go-jobqueue"
	"github.com/pieterclaerhout/go-jobqueue/environ"
	"github.com/pieterclaerhout/go-jobqueue/versioninfo"
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
			{
				Name:  "version",
				Usage: "show the version info",
				Action: func(c *cli.Context) error {
					log.Info(versioninfo.ProjectName, versioninfo.Version, "("+versioninfo.Revision+")")
					return nil
				},
			},
			{
				Name:  "setup",
				Usage: "creates the database tables",
				Action: func(c *cli.Context) error {

					r, err := jobqueue.DefaultRepository()
					if err != nil {
						return err
					}

					return r.Setup()

				},
			},
			{
				Name:  "queue",
				Usage: "adds a job to the queue",
				Action: func(c *cli.Context) error {

					r, err := jobqueue.DefaultRepository()
					if err != nil {
						return err
					}

					j, err := jobqueue.NewJob("email", jobqueue.JobParams{
						"from":    "pieter.claerhout@gmail.com",
						"to":      "pieter@yellowduck.be",
						"subject": "hello world",
					})
					if err != nil {
						return err
					}

					if _, err := r.Queue(j); err != nil {
						return err
					}

					return nil

				},
			},
			{
				Name:  "queue-many",
				Usage: "add many jobs to the queue",
				Action: func(c *cli.Context) error {

					r, err := jobqueue.DefaultRepository()
					if err != nil {
						return err
					}

					for i := 1; i <= 10; i++ {

						j, err := jobqueue.NewJob("email", jobqueue.JobParams{
							"from":     "pieter.claerhout@gmail.com",
							"to":       "pieter@yellowduck.be",
							"subject":  "hello world",
							"sequence": i,
						})
						if err != nil {
							return err
						}

						if _, err := r.Queue(j); err != nil {
							return err
						}

					}

					return nil

				},
			},
		},
	}

	err := app.Run(os.Args)
	log.CheckError(err)

}
