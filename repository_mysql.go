package jobqueue

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pieterclaerhout/go-log"
)

type MySQLRepository struct {
	db        *sqlx.DB
	tableName string
}

func NewMySQLRepository(db *sqlx.DB) Repository {
	return &MySQLRepository{
		db:        db,
		tableName: defaultTableName,
	}
}

func (r *MySQLRepository) Setup() error {

	log.Info("Performing setup")

	statements := []string{
		`CREATE TABLE IF NOT EXISTS "` + r.tableName + `" (
			"id" bigint unsigned NOT NULL AUTO_INCREMENT,
			"jobtype" varchar(255) NOT NULL DEFAULT '',
			"state" varchar(255) NOT NULL DEFAULT '',
			"payload" json DEFAULT NULL,
			"created_on" int NOT NULL DEFAULT '0',
			"started_on" int NOT NULL DEFAULT '0',
			"finished_on" int NOT NULL DEFAULT '0',
			PRIMARY KEY ("id"),
			UNIQUE KEY "id" ("id"),
			KEY "jobqueue_state" ("state"),
			KEY "jobqueue_created_on" ("created_on")
		)`,
	}

	for _, stmt := range statements {
		if _, err := r.db.Exec(stmt); err != nil {
			log.Error(err)
		}
	}

	return nil

}

func (r *MySQLRepository) Queue(job *Job) (*Job, error) {

	log.InfoDump(job, "Queueing job:")

	job.State = statusQueued

	result, err := r.db.NamedExec(
		`INSERT INTO "`+r.tableName+`" (
			"jobtype", "state", "payload", "created_on"
		) VALUES (
			:jobtype, :state, :payload, :created_on
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
		FROM "`+r.tableName+`"
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

	job.StartedOn = time.Now().Unix()
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
