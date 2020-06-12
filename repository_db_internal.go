package jobqueue

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func (r *DBRepository) dequeueJob(queue string) (*Job, error) {

	trx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	job := &Job{}

	if err := trx.Get(
		job,
		r.db.Rebind(`SELECT * FROM "`+r.tableName+`"
		WHERE "state" = ? AND "queue" = ?
		ORDER BY "created_on"
		LIMIT 1
		FOR UPDATE SKIP LOCKED`),
		statusQueued, queue,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, trx.Commit()
		}
		return nil, err
	}

	if err := r.updateJob(trx, job.markAsStarted()); err != nil {
		return job, err
	}

	return job, nil

}

func (r *DBRepository) finishJob(job *Job, err error) error {
	return r.updateJob(nil, job.markAsFinished(err))
}

func (r *DBRepository) updateJob(trx *sqlx.Tx, job *Job) error {

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

func (r *DBRepository) setupForDBType(dbType string) []string {

	if dbType == "mysql" {
		return []string{
			`CREATE TABLE IF NOT EXISTS "` + r.tableName + `" (
				"id" bigint unsigned NOT NULL AUTO_INCREMENT,
				"queue" varchar(255) NOT NULL DEFAULT '',
				"state" varchar(255) NOT NULL DEFAULT '',
				"payload" json,
				"error" longtext,
				"created_on" int NOT NULL DEFAULT '0',
				"started_on" int NOT NULL DEFAULT '0',
				"finished_on" int NOT NULL DEFAULT '0',
				PRIMARY KEY ("id")
			)`,
			`CREATE INDEX "` + r.tableName + `_queue" ON "` + r.tableName + `" ("queue")`,
			`CREATE INDEX "` + r.tableName + `_state" ON "` + r.tableName + `" ("state")`,
			`CREATE INDEX "` + r.tableName + `_created_on" ON "` + r.tableName + `" ("created_on")`,
		}
	}

	if dbType == "postgres" {
		return []string{
			`CREATE TABLE IF NOT EXISTS "` + r.tableName + `" (
				"id" bigserial,
				"queue" varchar NOT NULL DEFAULT '',
				"state" varchar NOT NULL DEFAULT '',
				"payload" json,
				"error" text NOT NULL DEFAULT '',
				"created_on" int NOT NULL DEFAULT 0,
				"started_on" int NOT NULL DEFAULT 0,
				"finished_on" int NOT NULL DEFAULT 0,
				PRIMARY KEY ("id")
			)`,
			`CREATE INDEX IF NOT EXISTS "` + r.tableName + `_queue" ON "public"."jobs" USING BTREE ("queue")`,
			`CREATE INDEX IF NOT EXISTS "` + r.tableName + `_state" ON "public"."jobs" USING BTREE ("state")`,
			`CREATE INDEX IF NOT EXISTS "` + r.tableName + `_created_on" ON "public"."jobs" USING BTREE ("created_on")`,
		}

	}

	return []string{}

}
