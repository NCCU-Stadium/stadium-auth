package database

import (
	"context"
	"database/sql"
)

func (db *Database) GetUserByMail(ctx context.Context, email string) (User, error) {
	row := db.database.QueryRowContext(ctx, `select mail, role, phone, pass from user_t where mail = $1`, email)

	if row.Err() == sql.ErrNoRows {
		return User{}, sql.ErrNoRows
	}
	if row.Err() != nil {
		return User{}, row.Err()
	}

	var u User
	err := row.Scan(&u.Mail, &u.Role, &u.Phone, &u.Pass)
	if err != nil {
		return User{}, err
	}

	return u, nil
}
