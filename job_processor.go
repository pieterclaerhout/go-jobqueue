package jobqueue

type JobProcessor interface {
	Process(job *Job) error
}
