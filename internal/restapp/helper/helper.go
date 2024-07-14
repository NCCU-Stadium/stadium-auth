package restapp_helper

import (
	"auth-service/internal/config"
	"auth-service/internal/database"
	"context"
	"database/sql"
	"errors"
)

type RestHelper struct {
	SQLDatabase *database.Database
}

func NewRestHelper(c *config.Config) *RestHelper {
	db, err := database.New(c.DatabaseURI)
	if err != nil {
		panic(err)
	}
	return &RestHelper{
		SQLDatabase: db,
	}
}

var ErrorUserNotFound = errors.New("User not found")

func (rh *RestHelper) GetUserName(ctx context.Context, mail string) (string, error) {
	user, err := rh.SQLDatabase.GetUserByEmail(ctx, mail)
	if err == sql.ErrNoRows {
		return "", ErrorUserNotFound
	}
	if err != nil {
		return "", err
	}
	return user.Name, nil
}

func (rh *RestHelper) GetUserByEmail(ctx context.Context, mail string) (database.User, error) {
	user, err := rh.SQLDatabase.GetUserByEmail(ctx, mail)
	if err == sql.ErrNoRows {
		return user, ErrorUserNotFound
	}
	if err != nil {
		return user, err
	}
	return user, nil
}
