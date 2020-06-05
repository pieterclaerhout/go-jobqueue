package main

import (
	"github.com/pieterclaerhout/go-jobqueue"
	"github.com/urfave/cli/v2"
)

var commandQueue = &cli.Command{
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

		if _, err := r.AddJob(j); err != nil {
			return err
		}

		return nil

	},
}
