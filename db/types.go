package db


type Database interface {
	Open(source string) error
	Migrate(source string) error
	Close() error
	Backend()
	Tx() error
	Commit() error
}


// use dtm as distributed transaction framework
type DTM interface {


}
