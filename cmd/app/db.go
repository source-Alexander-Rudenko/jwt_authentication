package app

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
)

func NewMySQLAStorage(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("Unable to connect to DB", err)
	}
	return db, nil
}
