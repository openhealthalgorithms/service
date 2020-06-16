package database

import (
	"database/sql"
	"log"

	// this package provides sqlite3 support
	_ "github.com/mattn/go-sqlite3"
)

// SqliteDb object
type SqliteDb struct {
	DB     *sql.DB
	Closer func()
}

// InitDb initializes the db
func InitDb(dbFilePath string) (*SqliteDb, error) {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}

	closer := func() {
		db.Close()
	}

	sqlite := SqliteDb{DB: db, Closer: closer}

	return &sqlite, nil
}

// Migrate a db
func (sq *SqliteDb) Migrate() error {
	sqlStmt := "create table if not exists logs (id integer not null primary key autoincrement, request text, response text, access_time timestamp default current_timestamp);delete from logs;"

	_, err := sq.DB.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}

	return nil
}
