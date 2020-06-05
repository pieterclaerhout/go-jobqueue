package jobqueue

import (
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/jmoiron/sqlx"
	"github.com/pieterclaerhout/go-jobqueue/environ"
	"github.com/pieterclaerhout/go-log"
)

type Repository interface {
	Setup() error
	Queue(job *Job) (*Job, error)
}

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
