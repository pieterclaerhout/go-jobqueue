package jobqueue

import (
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/pieterclaerhout/go-jobqueue/environ"
	"github.com/pieterclaerhout/go-log"
	"github.com/pkg/errors"
)

// JobRepository is the interface to which repositories should conform
type JobRepository interface {
	Setup()
	AddJob(job *Job) (*Job, error)
	Process(queue string, interval time.Duration, processor JobProcessor)
}

// DefaultRepository returns the default repository based on the env variables
func DefaultRepository() (JobRepository, error) {

	dsn := environ.String("DSN", "")
	dbType := environ.String("DB_TYPE", "mysql")

	if dsn == "" {
		return nil, errors.New("DSN env var is not set")
	}

	if dbType != "mysql" && dbType != "postgres" {
		return nil, errors.New("DB_TYPE " + dbType + " is not supported")
	}

	db, err := sqlx.Open(dbType, dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(10)

	log.Info("Connected to:", dsn, "type:", db.DriverName())

	return NewMySQLRepository(db, defaultTableName), nil

}
