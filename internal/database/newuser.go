package database

import (
	"context"
	"database/sql"
)

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func (db *Database) CreateUser(ctx context.Context, user CreateUserReq) error {

	qString := `
        with newuser as ( insert into user_t (mail, role, phone, pass) values ($1, $2, $3, $4) returning mail )
        insert into subuser_t (user_mail, name, avatar, gender, birth) values ($1, $5::varchar, $6::varchar, $7, $8::date) returning user_mail;
    `
	row := db.database.QueryRowContext(
		ctx, qString,
		user.Mail, user.Role, user.Phone, user.Pass, user.Name, nullString(user.Avatar), nullString(user.Gender), nullString(user.Birth),
	)
	if row.Err() == sql.ErrNoRows {
		return sql.ErrNoRows
	}
	if row.Err() != nil {
		return row.Err()
	}

	return nil
}
