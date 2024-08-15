package app

import (
	"auth-service/bcrypt"
	"auth-service/internal/helper"
	"auth-service/jwt"
	"encoding/json"
	"log"
	"net/http"
)

const BcryptCost = 12

type LoginResponse struct {
	Mail         string `json:"email"`
	Role         string `json:"role"`
	RefreshToken string `json:"refreshToken"`
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	s.credentials(w, r)
}

type LoginRequest struct {
	Mail string `json:"email"`
	Pass string `json:"password"`
}

func (s *Server) credentials(w http.ResponseWriter, r *http.Request) {
	in := &LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 1. Check if user exists and check password matches
	// 2. If exists, generate access token and set cookies
	// 3. Generate refresh token and put into refresh token db
	// 4. Return user info

	// 1. Check if user exists and check password matches
	user, err := s.helper.GetUserByMail(r.Context(), in.Mail)
	if err == helper.ErrorUserNotFound {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ok, err := bcrypt.Compare(in.Pass, user.Pass); !ok {
		log.Print(in.Pass)
		log.Print(user.Pass)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2. Generate access token and set cookies
	accessContent := &AccessContent{UserMail: in.Mail, UserRole: user.Role}
	accessToken, err := jwt.Sign(accessContent, s.config.Secret, "Bearer ", AccessTokenDuration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{ // Set a cookie
		Name:     "accessToken",
		Value:    accessToken,
		MaxAge:   int(AccessTokenDuration),
		Path:     "/",
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	// 3. Generate refresh token and put into refresh token db
	tokenid, err := s.refreshHelper.SaveRefreshMeta(user.Mail, RefreshTokenDuration)
	refreshContent := &RefreshContent{UserRole: user.Mail, TokenID: tokenid, UserMail: in.Mail}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	refreshToken, err := jwt.Sign(refreshContent, s.config.Secret, "Bearer ", RefreshTokenDuration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Return user info
	loginResponse := &LoginResponse{Mail: in.Mail, Role: user.Role, RefreshToken: refreshToken}

	err = json.NewEncoder(w).Encode(loginResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}
