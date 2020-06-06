package jobqueue

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pieterclaerhout/go-log"
)

// MySQLRepository defines a MySQL-based repository
type MySQLRepository struct {
	db        *sqlx.DB
	tableName string
}

// NewMySQLRepository returns a new MySQL-based repository
func NewMySQLRepository(db *sqlx.DB, tableName string) Repository {
	return &MySQLRepository{
		db:        db,
		tableName: tableName,
	}
}

// Setup is used to perform the setup of the repository
func (r *MySQLRepository) Setup() error {

	log.Info("Performing setup")

	statements := []string{
		`CREATE TABLE IF NOT EXISTS "` + r.tableName + `" (
			"id" bigint unsigned NOT NULL AUTO_INCREMENT,
			"queue" varchar(255) NOT NULL DEFAULT '',
			"state" varchar(255) NOT NULL DEFAULT '',
			"payload" json,
			"error" longtext,
			"created_on" int NOT NULL DEFAULT '0',
			"started_on" int NOT NULL DEFAULT '0',
			"finished_on" int NOT NULL DEFAULT '0',
			PRIMARY KEY ("id"),
			UNIQUE KEY "id" ("id"),
			KEY "jobqueue_queue" ("queue"),
			KEY "jobqueue_state" ("state"),
			KEY "jobqueue_created_on" ("created_on")
		)`,
	}

	for _, stmt := range statements {
		if _, err := r.db.Exec(stmt); err != nil {
			log.Error(err)
		}
	}

	log.Info("Performed setup")

	return nil

}

// AddJob adds a job to the queue
func (r *MySQLRepository) AddJob(job *Job) (*Job, error) {

	job.markAsQueued()

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

	job.ID, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	log.InfoDump(job, "Queued job:")

	return job, nil

}

// Process starts processing jobs from the given queue with the given interval
func (r *MySQLRepository) Process(queue string, interval time.Duration, processor JobProcessor) error {

	log.Debug("Processing jobs from queue:", queue, "interval:", interval)

	for {

		log.Debug("Checking for jobs in queue:", queue)

		job, err := r.dequeueJob(queue)
		if err != nil {
			log.Error(err)
			time.Sleep(interval)
		}

		if job == nil {
			time.Sleep(interval)
			continue
		}

		jobErr := processor.Process(job)

		if err := r.finishJob(job, jobErr); err != nil {
			log.Error(err)
			time.Sleep(interval)
		}

	}

	return nil

}
