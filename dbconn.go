package main

import (
	"database/sql"
)

func dbConnect() (*sql.DB, error) {

	dbConn, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}

	if err = dbConn.Ping(); err != nil {
		return nil, err
	}

	return dbConn, nil
}
