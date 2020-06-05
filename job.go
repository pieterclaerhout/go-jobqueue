package jobqueue

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID        int64
	UUID      string
	JobType   string
	State     string
	Payload   string
	CreatedOn int64 `db:"created_on"`
}

func NewJob(jobType string, payload JobParams) (*Job, error) {

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Job{
		UUID:      uuid.New().String(),
		JobType:   jobType,
		State:     "queued",
		Payload:   string(payloadBytes),
		CreatedOn: time.Now().Unix(),
	}, nil

}
