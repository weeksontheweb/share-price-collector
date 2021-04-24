package database

import (
	"database/sql"
	"fmt"
)

type (
	ShareDB struct {
		db *sql.DB
	}
)

func (db *ShareDB) ConnectToDatabase(host string, port int, user string, password string, dbname string) (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	newdb, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, err
	}

	return newdb, nil
}

func ConnectToDatabase(host string, port int, user string, password string, dbname string) (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	newdb, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, err
	}

	return newdb, nil
}
