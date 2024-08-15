package helper

import (
	"auth-service/internal/config"
	"auth-service/internal/database"
	"context"
	"errors"
)

type Helper struct {
	SQLDatabase *database.Database
}

func NewRestHelper(c *config.Config) *Helper {
	db, err := database.New(c.DatabaseURI)
	if err != nil {
		panic(err)
	}
	return &Helper{
		SQLDatabase: db,
	}
}

var ErrorUserNotFound = errors.New("User not found")

type User struct {
	Mail string
	Role string
	Pass string
}

type Subuser struct{}

func (h *Helper) GetUserByMail(ctx context.Context, mail string) (User, error) {
	user, err := h.SQLDatabase.GetUserByMail(ctx, mail)
	if err != nil {
		return User{}, err
	}
	if user.Mail == "" {
		return User{}, ErrorUserNotFound
	}
	return User{Mail: user.Mail, Role: user.Role, Pass: user.Pass}, nil
}

func (h *Helper) IsUserExist(ctx context.Context, mail string) (bool, error) {
	user, err := h.SQLDatabase.GetUserByMail(ctx, mail)
	if err != nil {
		return false, err
	}
	if user.Mail == "" {
		return false, nil
	}
	return true, nil
}

type RUserReq struct {
	Mail   string
	Pass   string
	Role   string
	Phone  string
	Name   string
	Avatar string
	Gender string
	Birth  string
}

func (r RUserReq) ToDBUser() database.CreateUserReq {
	return database.CreateUserReq{
		Mail:   r.Mail,
		Role:   r.Role,
		Phone:  r.Phone,
		Pass:   r.Pass,
		Name:   r.Name,
		Avatar: r.Avatar,
		Gender: r.Gender,
		Birth:  r.Birth,
	}
}

func (h *Helper) RegisterUser(ctx context.Context, user RUserReq) error {
	err := h.SQLDatabase.CreateUser(ctx, user.ToDBUser())
	if err != nil {
		return err
	}
	return nil
}
