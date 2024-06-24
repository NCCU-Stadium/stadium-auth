package database

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	database *sql.DB
}

func Connect(connString string) (Database, error) {
	// Connect to database
	db, err := sql.Open("pgx", connString)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return Database{
		database: db,
	}, nil
}
