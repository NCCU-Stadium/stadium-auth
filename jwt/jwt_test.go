package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Username string
	Password string
	jwt.RegisteredClaims
}

func (u *User) ToDomain(m map[string]interface{}) {
	u.Username = m["Username"].(string)
	u.Password = m["Password"].(string)
	return
}

func (u *User) SetClaims(claims jwt.RegisteredClaims) {
	u.RegisteredClaims = claims
}

func TestSign(t *testing.T) {
	user := &User{
		Username: "admin",
		Password: "password",
	}
	key := "secret"
	token, err := Sign(user, key, time.Hour*24) // One day expiration
	if err != nil {
		t.Log(err)
		panic(err)
	}
	t.Log(token)
}

func _TestParse(t *testing.T) {
	jwtUser := &User{
		Username: "admin",
		Password: "password",
	}
	key := "secret"
	token, err := Sign(jwtUser, key, time.Hour*24) // One day expiration
	tokenString := token
	claims, err := Parse(tokenString, key)
	if err != nil {
		t.Log(err)
	}

	var user User
	user.ToDomain(claims)
	t.Log(user)
}
