package jobqueue

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pieterclaerhout/go-log"
)

type MySQLRepository struct {
	db *sqlx.DB
}

func NewMySQLRepository(db *sqlx.DB) Repository {
	return &MySQLRepository{
		db: db,
	}
}

func (r *MySQLRepository) Setup() error {

	log.Info("Performing setup")

	statements := []string{
		`CREATE TABLE IF NOT EXISTS "jobqueue" (
			"id" bigint unsigned NOT NULL AUTO_INCREMENT,
			"uuid" varchar(36) NOT NULL DEFAULT '',
			"jobtype" varchar(255) NOT NULL DEFAULT '',
			"state" varchar(255) NOT NULL DEFAULT '',
			"payload" json DEFAULT NULL,
			"created_on" int NOT NULL DEFAULT '0',
			PRIMARY KEY ("id"),
			UNIQUE KEY "id" ("id"),
			UNIQUE KEY "jobqueue_uuid" ("uuid"),
			KEY "jobqueue_state" ("state"),
			KEY "jobqueue_created_on" ("created_on")
		)`,
		`ALTER TABLE "jobqueue" ADD COLUMN "finished_on" int NOT NULL DEFAULT '0'`,
	}

	for _, stmt := range statements {
		if _, err := r.db.Exec(stmt); err != nil {
			return err
		}
	}

	return nil

}

func (r *MySQLRepository) Queue(job *Job) (*Job, error) {

	log.InfoDump(job, "Queueing job:")

	job.State = statusQueued

	result, err := r.db.NamedExec(
		`INSERT INTO "jobqueue" (
			"uuid",
			"jobtype",
			"state",
			"payload",
			"created_on"
		) VALUES (
			:uuid,
			:jobtype,
			:state,
			:payload,
			:created_on
		)`,
		job,
	)
	if err != nil {
		return nil, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	job.ID = lastInsertID

	log.InfoDump(job, "Queued job:")

	return job, nil

}

func (r *MySQLRepository) Dequeue(jobType string) (*Job, error) {

	trx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	job := &Job{}

	if err := trx.Get(
		job,
		`SELECT *
		FROM "jobqueue"
		WHERE "state" = ? AND "jobtype" = ?
		ORDER BY "created_on"
		LIMIT 1
		FOR UPDATE SKIP LOCKED`,
		statusQueued, jobType,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, trx.Commit()
		}
		return nil, err
	}

	job.State = statusRunning

	if err := r.setJobStatus(trx, job); err != nil {
		return job, err
	}

	return job, nil

}

func (r *MySQLRepository) FailJob(job *Job) error {
	job.State = statusError
	job.FinishedOn = time.Now().Unix()
	return r.setJobStatus(nil, job)
}

func (r *MySQLRepository) FinishJob(job *Job) error {
	job.State = statusFinished
	job.FinishedOn = time.Now().Unix()
	return r.setJobStatus(nil, job)
}
