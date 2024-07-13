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

func _TestSign(t *testing.T) {
	user := &User{
		Username: "admin",
		Password: "password",
	}
	key := "secret"
	token, err := Sign(user, key, "Bearer ", time.Hour*24) // One day expiration
	if err != nil {
		t.Log(err)
		panic(err)
	}
	t.Log(token)
}

func TestParse(t *testing.T) {
	jwtUser := &User{
		Username: "admin",
		Password: "password",
	}
	key := "secret"
	token, err := Sign(jwtUser, key, "Bearer ", time.Hour*24) // One day expiration
	t.Log("token: ", token)

	tokenString := token
	claims, err := Parse(tokenString, key, "Bearer ")
	if err != nil {
		t.Log(err)
	}

	var user User
	user.ToDomain(claims)
	t.Log(user)
}
