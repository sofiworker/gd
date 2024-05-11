package db

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
)

type Sqlx struct {
	SqlxDB *sqlx.DB
}

func (s *Sqlx) Migrate(source string) error {
	err := s.Open(source)
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
        "file:///migrations",
        "postgres", nil)
    return m.Up()
}

func (s *Sqlx) Open(source string) error {
	db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	if err != nil {
		return err
	}
	s.SqlxDB = db
	return db.Ping()
}

func (s *Sqlx) Close() error {
	//TODO implement me
	panic("implement me")
}

func (s *Sqlx) Backend() {
	//TODO implement me
	panic("implement me")
}

func (s *Sqlx) Tx() error {
	//TODO implement me
	panic("implement me")
}

func (s *Sqlx) Commit() error {
	//TODO implement me
	panic("implement me")
}
