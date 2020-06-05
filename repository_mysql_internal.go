package jobqueue

import (
	"github.com/jmoiron/sqlx"
)

func (r *MySQLRepository) setJobStatus(trx *sqlx.Tx, job *Job) error {

	if trx == nil {
		var err error
		if trx, err = r.db.Beginx(); err != nil {
			return err
		}
	}

	if _, err := trx.Exec(
		`UPDATE "`+r.tableName+`" SET "state" = ?, "started_on" = ?, "finished_on" = ? WHERE "id" = ?`,
		job.State, job.StartedOn, job.FinishedOn, job.ID,
	); err != nil {
		return err
	}

	return trx.Commit()

}
