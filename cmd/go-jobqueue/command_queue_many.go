package main

import (
	"github.com/pieterclaerhout/go-jobqueue"
	"github.com/pieterclaerhout/go-log"
	"github.com/urfave/cli/v2"
)

var commandQueueMany = &cli.Command{
	Name:  "queue-many",
	Usage: "add many jobs to the queue",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "queue",
			Usage: "The queue to process",
			Value: "email",
		},
		&cli.IntFlag{
			Name:  "count",
			Usage: "The number of jobs to add",
			Value: 10,
		},
	},
	Action: func(c *cli.Context) error {

		r, err := jobqueue.DefaultRepository()
		if err != nil {
			return err
		}

		count := c.Int("count")
		queue := c.String("queue")
		log.Info("Queueing", count, "jobs in queue:", queue)

		for i := 1; i <= count; i++ {

			j, err := jobqueue.NewJob(queue, jobqueue.JobParams{
				"from":     "pieter.claerhout@gmail.com",
				"to":       "pieter@yellowduck.be",
				"subject":  "hello world",
				"sequence": i,
			})
			if err != nil {
				return err
			}

			if _, err := r.AddJob(j); err != nil {
				return err
			}

		}

		return nil

	},
}
