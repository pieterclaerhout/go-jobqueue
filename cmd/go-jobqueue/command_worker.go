package main

import (
	"errors"
	"time"

	"github.com/pieterclaerhout/go-jobqueue"
	"github.com/pieterclaerhout/go-log"
	"github.com/urfave/cli/v2"
)

var commandWorker = &cli.Command{
	Name:  "worker",
	Usage: "runs the worker process for a given job type",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "queue",
			Usage: "The queue to process",
		},
		&cli.DurationFlag{
			Name:  "interval",
			Usage: "how often we should check for new jobs",
			Value: 1 * time.Second,
		},
	},
	Action: func(c *cli.Context) error {

		r, err := jobqueue.DefaultRepository()
		if err != nil {
			return err
		}

		queue := c.String("queue")
		interval := c.Duration("interval")

		return r.Process(
			queue,
			interval,
			func(job *jobqueue.Job) error {

				log.Info("Processing job:", job.ID)

				if job.ID%3 == 0 {
					log.Error("Failing job:", job.ID)
					return errors.New("job error message")
				}

				time.Sleep(500 * time.Millisecond)

				log.Info("Processed job:", job.ID)

				return nil

			},
		)

	},
}
