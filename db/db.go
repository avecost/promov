package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

type Tx struct {
	*sql.Tx
}

// Open returns a DB reference for a data source
func Open(connStr string) (*DB, error) {
	appDb, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	if err = appDb.Ping(); err != nil {
		return nil, err
	}

	return &DB{appDb}, nil
}

// Begin starts and return a new transaction
func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}

	return &Tx{tx}, nil
}
