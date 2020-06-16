package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteDb struct {
	DB     *sql.DB
	Closer func()
}

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

func (sq *SqliteDb) Migrate() error {
	sqlStmt := "create table if not exists logs (id integer not null primary key autoincrement, request text, response text, access_time timestamp default current_timestamp);delete from logs;"

	_, err := sq.DB.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}

	return nil
}
