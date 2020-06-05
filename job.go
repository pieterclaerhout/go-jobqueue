package jobqueue

import (
	"encoding/json"
	"time"

	"github.com/tidwall/gjson"
)

// Job defines a job which can be queued
type Job struct {
	ID         int64  `db:"id"`
	Queue      string `db:"queue"`
	State      string `db:"state"`
	Error      string `db:"error"`
	Payload    string `db:"payload"`
	CreatedOn  int64  `db:"created_on"`
	StartedOn  int64  `db:"started_on"`
	FinishedOn int64  `db:"finished_on"`
}

// NewJob returns a new job for the given queue with an optional payload
func NewJob(queue string, payload JobPayload) (*Job, error) {

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Job{
		Queue:     queue,
		Payload:   string(payloadBytes),
		CreatedOn: time.Now().Unix(),
	}, nil

}

// StringArg retrieves a string argument from the payload
func (job *Job) StringArg(name string) string {
	return gjson.Get(job.Payload, name).String()
}

// IntArg retrieves an integer argument from the payload
func (job *Job) IntArg(name string) int64 {
	return gjson.Get(job.Payload, name).Int()
}

// FloatArg retrieves an integer argument from the payload
func (job *Job) FloatArg(name string) float64 {
	return gjson.Get(job.Payload, name).Float()
}

// BoolArg retrieves a bool argument from the payload
func (job *Job) BoolArg(name string) bool {
	return gjson.Get(job.Payload, name).Bool()
}
