package jobqueue

// JobProcessor defines the interface for a job processor
type JobProcessor interface {
	Process(job *Job) error
}
