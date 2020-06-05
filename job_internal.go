package jobqueue

import (
	"time"
)

func (job *Job) markAsQueued() *Job {
	job.CreatedOn = time.Now().Unix()
	job.StartedOn = 0
	job.FinishedOn = 0
	job.Error = ""
	job.State = statusQueued
	return job
}

func (job *Job) markAsStarted() *Job {
	job.StartedOn = time.Now().Unix()
	job.FinishedOn = 0
	job.Error = ""
	job.State = statusRunning
	return job
}

func (job *Job) markAsFinished(err error) *Job {
	if err != nil {
		job.State = statusError
		job.Error = err.Error()
	} else {
		job.State = statusFinished
		job.Error = ""
	}
	job.FinishedOn = time.Now().Unix()
	return job
}
