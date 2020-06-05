package main

import (
	"os"
	"time"

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
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "jobtype",
						Usage: "The jobtype to process",
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
					jobType := c.String("jobtype")
					log.Info("Queueing", count, "jobs of type:", jobType)

					for i := 1; i <= count; i++ {

						j, err := jobqueue.NewJob(jobType, jobqueue.JobParams{
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
			{
				Name:  "worker",
				Usage: "runs the worker process for a given job type",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "jobtype",
						Usage: "The jobtype to process",
					},
					&cli.DurationFlag{
						Name:  "interval",
						Usage: "how often we should check for new jobs",
						Value: 250 * time.Millisecond,
					},
				},
				Action: func(c *cli.Context) error {

					r, err := jobqueue.DefaultRepository()
					if err != nil {
						return err
					}

					jobType := c.String("jobtype")
					interval := c.Duration("interval")
					log.Info("Processing jobs with type:", jobType, "interval:", interval)

					for {

						log.Debug("Checking for jobs:", jobType)

						job, err := r.Dequeue(jobType)
						if err != nil {
							log.Error(err)
						}

						if job == nil {
							time.Sleep(interval)
							continue
						}

						log.Info("Processing job:", job.ID)
						time.Sleep(500 * time.Millisecond)
						if err := r.FinishJob(job); err != nil {
							log.Error("Failed job:", err)
						} else {
							log.Info("Processed job:", job.ID)
						}

					}

					return nil

				},
			}},
	}

	err := app.Run(os.Args)
	log.CheckError(err)

}
