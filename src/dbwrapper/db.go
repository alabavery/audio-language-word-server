package dbwrapper

import (
	"database/sql"
	"log"
)

type DBWrapper struct {
	db *sql.DB
}

func Open() *DBWrapper {
	db, err := sql.Open(
		"postgres",
		"user=al dbname=words sslmode=veriy-full",
	)
	if err != nil {
		log.Fatal(err)
	}
	return &DBWrapper{db: db}
}

func (db *DBWrapper) Insert(table string, columns []string, rows *[][]interface{}) {

}
