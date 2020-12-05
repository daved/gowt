package main

import (
	"database/sql"
	"fmt"
	"os"
)

func dbConn() (*sql.DB, error) {
	efmt := "dbconn: %w"

	dbDriver := "mysql"
	dbUser := os.Getenv("DATABASE_USERNAME")
	dbPass := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_NAME")
	dbServer := os.Getenv("DATABASE_SERVER")
	dbPort := os.Getenv("DATABASE_PORT")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbServer+":"+dbPort+")/"+dbName)
	if err != nil {
		return nil, fmt.Errorf(efmt, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf(efmt, err)
	}

	return db, nil
}
