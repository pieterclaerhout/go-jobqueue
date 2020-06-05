package jobqueue

import (
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/jmoiron/sqlx"
	"github.com/pieterclaerhout/go-jobqueue/environ"
	"github.com/pieterclaerhout/go-log"
)

// Repository is the interface to which repositories should conform
type Repository interface {
	Setup() error
	AddJob(job *Job) (*Job, error)
	Process(queue string, interval time.Duration, processor JobProcessor) error
}

// DefaultRepository returns the default repository based on the env variables
func DefaultRepository() (Repository, error) {

	dsn := environ.String("DSN", "")
	dbType := environ.String("DB_TYPE", "mysql")

	db, err := sqlx.Open(dbType, dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(10)

	log.Info("Connected to:", dsn)

	return NewMySQLRepository(db), nil

}
