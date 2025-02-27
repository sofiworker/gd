package db

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
)

type SqlxClient struct {
	SqlxDB *sqlx.DB
}

func NewSqlxClient() *SqlxClient {
	return &SqlxClient{}
}

func (s *SqlxClient) Migrate(source string) error {
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", nil)
	if err != nil {
		return err
	}
	return m.Up()
}

//func (s *Sqlx) Open(source string) error {
//	db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
//	if err != nil {
//		return err
//	}
//	s.SqlxDB = db
//	return db.Ping()
//}

func (s *SqlxClient) Close() error {
	//TODO implement me
	panic("implement me")
}

func (s *SqlxClient) Backend() {
	//TODO implement me
	panic("implement me")
}

func (s *SqlxClient) Tx() error {
	//TODO implement me
	panic("implement me")
}

func (s *SqlxClient) Commit() error {
	//TODO implement me
	panic("implement me")
}
