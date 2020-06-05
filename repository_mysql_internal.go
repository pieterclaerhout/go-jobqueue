package jobqueue

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

func (r *MySQLRepository) dequeueJob(queue string) (*Job, error) {

	trx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	job := &Job{}

	if err := trx.Get(
		job,
		`SELECT *
		FROM "`+r.tableName+`"
		WHERE "state" = ? AND "queue" = ?
		ORDER BY "created_on"
		LIMIT 1
		FOR UPDATE SKIP LOCKED`,
		statusQueued, queue,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, trx.Commit()
		}
		return nil, err
	}

	job.StartedOn = time.Now().Unix()
	job.State = statusRunning

	if err := r.updateJob(trx, job); err != nil {
		return job, err
	}

	return job, nil

}

func (r *MySQLRepository) finishJob(job *Job, err error) error {
	job.State = statusFinished
	job.Error = ""
	if err != nil {
		job.State = statusError
		job.Error = err.Error()
	}
	job.FinishedOn = time.Now().Unix()
	return r.updateJob(nil, job)
}

func (r *MySQLRepository) updateJob(trx *sqlx.Tx, job *Job) error {

	if trx == nil {
		var err error
		if trx, err = r.db.Beginx(); err != nil {
			return err
		}
	}

	if _, err := trx.NamedExec(
		`UPDATE "`+r.tableName+`" SET "state" = :state, "error" = :error, "started_on" = :started_on, "finished_on" = :finished_on WHERE "id" = :id`,
		job,
	); err != nil {
		return err
	}

	return trx.Commit()

}
