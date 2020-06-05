package processors

import (
	"errors"
	"time"

	"github.com/pieterclaerhout/go-jobqueue"
	"github.com/pieterclaerhout/go-log"
)

type EmailProcessor struct{}

func NewEmailProcessor() *EmailProcessor {
	return &EmailProcessor{}
}

func (p *EmailProcessor) Process(job *jobqueue.Job) error {

	log.Info("Processing job:", job.ID, job.Payload)

	if job.ID%3 == 0 {
		log.Error("Failing job:", job.ID)
		return errors.New("job error message")
	}

	from := job.StringArg("from")
	sequence := job.IntArg("sequence")
	unknown := job.StringArg("unknown")

	log.Info(from, sequence, unknown)

	time.Sleep(500 * time.Millisecond)

	log.Info("Processed job:", job.ID, job.Payload)

	return nil

}
