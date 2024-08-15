package database

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	database *sql.DB
}

func New(connString string) (*Database, error) {
	// Connect to database
	db, err := sql.Open("pgx", connString)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return &Database{
		database: db,
	}, nil
}

// CREATE TABLE IF NOT EXISTS "user_t" (
//   "mail" varchar PRIMARY KEY,
//   "role" varchar,
//   "phone" varchar,
//   "created_at" timestamp with time zone not null default now(),
//   "pass" varchar,
//   "point" integer not null default 0,
//   "unpaid" integer not null default 0
// );

type User struct {
	Mail      string
	Role      string
	Phone     string
	CreatedAt string
	Pass      string
	Point     int
	Unpaid    int
}

// CREATE TABLE IF NOT EXISTS "subuser_t" (
//   "user_mail" varchar,
//   "avatar" varchar,
//   "created_at" timestamp with time zone not null default now(),
//   "name" varchar,
//   "gender" varchar,
//   "birth" date,
//   PRIMARY KEY ("user_mail", "name")
// );

type Subuser struct {
	Avatar    sql.NullString
	CreatedAt string
	Name      string
	Gender    sql.NullString
	Birth     sql.NullString
}

type CreateUserReq struct {
	Mail   string
	Role   string
	Phone  string
	Pass   string
	Name   string
	Avatar string
	Gender string
	Birth  string
}
