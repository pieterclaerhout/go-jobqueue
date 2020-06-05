package jobqueue

import (
	"encoding/json"
	"time"
)

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

func NewJob(queue string, payload JobParams) (*Job, error) {

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
