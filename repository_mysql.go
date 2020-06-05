package jobqueue

import (
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

	if _, err := r.db.Exec(
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
	); err != nil {
		return err
	}

	return nil

}

func (r *MySQLRepository) Queue(job *Job) (*Job, error) {

	log.InfoDump(job, "Queueing job:")

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
