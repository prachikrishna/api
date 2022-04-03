package main

import (
	"database/sql"
	"log"
)

var connString string

func init() {
	connString = "server=IT1806139;database=Contacts"
}

func dbGetConn() *sql.DB {
	_, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}

	return nil
}
