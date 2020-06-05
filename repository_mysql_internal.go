package jobqueue

import (
	"github.com/jmoiron/sqlx"
)

func (r *MySQLRepository) setJobStatus(trx *sqlx.Tx, job *Job) error {

	if _, err := trx.Exec(
		`UPDATE "jobqueue" SET "state" = ?, "finished_on" = ? WHERE "id" = ?`,
		job.State, job.FinishedOn, job.ID,
	); err != nil {
		return err
	}

	return trx.Commit()

}
