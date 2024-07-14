package database

import (
	"context"
	"database/sql"
)

type User struct {
	Name     string
	Email    string
	Password string
	Avatar   string
	Id       string
}

func (db *Database) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := db.database.QueryRowContext(ctx, `select name, mail, pass, avatar, id from user_t where mail = $1`, email)
	if row.Err() == sql.ErrNoRows {
		return User{}, sql.ErrNoRows
	}
	if row.Err() != nil {
		return User{}, row.Err()
	}

	user := User{}
	err := row.Scan(&user.Name, &user.Email, &user.Password, &user.Avatar, &user.Id)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
