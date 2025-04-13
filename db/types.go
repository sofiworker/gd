package db

type Database interface {
	Migrate(source string) error
	Close() error
	Backend()
	Tx() error
	Commit() error
}

// DTM use dtm as distributed transaction framework
type DTM interface {
}
