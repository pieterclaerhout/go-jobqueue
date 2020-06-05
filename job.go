package jobqueue

import (
	"encoding/json"
	"time"
)

type Job struct {
	ID         int64
	JobType    string
	State      string
	Payload    string
	CreatedOn  int64 `db:"created_on"`
	StartedOn  int64 `db:"started_on"`
	FinishedOn int64 `db:"finished_on"`
}

func NewJob(jobType string, payload JobParams) (*Job, error) {

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Job{
		JobType:    jobType,
		Payload:    string(payloadBytes),
		CreatedOn:  time.Now().Unix(),
		StartedOn:  0,
		FinishedOn: 0,
	}, nil

}
