package jobqueue

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pieterclaerhout/go-log"
)

// DBRepository defines a MySQL-based repository
type DBRepository struct {
	db        *sqlx.DB
	tableName string
}

// NewMySQLRepository returns a new MySQL-based repository
func NewMySQLRepository(db *sqlx.DB, tableName string) JobRepository {
	r := &DBRepository{
		db:        db,
		tableName: tableName,
	}
	r.Setup()
	return r
}

// Setup is used to perform the setup of the repository
func (r *DBRepository) Setup() {

	log.Info("Performing setup")

	statements := r.setupForDBType(r.db.DriverName())

	for _, stmt := range statements {
		log.Debug("executing:", stmt)
		if _, err := r.db.Exec(stmt); err != nil {
			if strings.HasPrefix(err.Error(), "Error 1061:") {
				continue
			}
			log.Error(err)
		}
	}

	log.Info("Performed setup")

}

// AddJob adds a job to the queue
func (r *DBRepository) AddJob(job *Job) (*Job, error) {

	job.markAsQueued()

	if r.db.DriverName() == "postgres" {

		stmt, err := r.db.PrepareNamed(
			`INSERT INTO "` + r.tableName + `" (
			"queue", "state", "error", "payload", "created_on", "started_on", "finished_on"
		) VALUES (
			:queue, :state, :error, :payload, :created_on, :started_on, :finished_on
		) RETURNING "id"`,
		)
		if err != nil {
			return nil, err
		}

		if err := stmt.Get(&job.ID, job); err != nil {
			return nil, err
		}

	}

	if r.db.DriverName() == "mysql" {

		result, err := r.db.NamedExec(
			`INSERT INTO "`+r.tableName+`" (
				"queue", "state", "error", "payload", "created_on", "started_on", "finished_on"
			) VALUES (
				:queue, :state, :error, :payload, :created_on, :started_on, :finished_on
			)`,
			job,
		)
		if err != nil {
			return nil, err
		}

		if job.ID, err = result.LastInsertId(); err != nil {
			return nil, err
		}

	}

	log.InfoDump(job, "Queued job:")

	return job, nil

}

// Process starts processing jobs from the given queue with the given interval
func (r *DBRepository) Process(queue string, interval time.Duration, processors map[string]JobProcessor) {

	queues := []string{}
	for queue := range processors {
		queues = append(queues, queue)
	}

	log.Debug("Processing jobs from queue(s):", strings.Join(queues, "|"), "interval:", interval)

	for {

		log.Debug("Checking for jobs in queue(s):", strings.Join(queues, "|"))

		job, err := r.dequeueJob(queues)
		if err != nil {
			log.Error(err)
			time.Sleep(interval)
		}

		if job == nil {
			time.Sleep(interval)
			continue
		}

		processor := processors[job.Queue]

		jobErr := processor.Process(job)

		if err := r.finishJob(job, jobErr); err != nil {
			log.Error(err)
			time.Sleep(interval)
		}

	}

}
