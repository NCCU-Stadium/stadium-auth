package restapp

import (
	"auth-service/bcrypt"
	"auth-service/internal/restapp/helper"
	"auth-service/jwt"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

const BcryptCost = 12

type LoginResponse struct {
	Name         string `json:"name"`
	Mail         string `json:"email"`
	Avatar       string `json:"avatar"`
	RefreshToken string `json:"refreshToken"`
}

func (s *RestServer) Login(w http.ResponseWriter, r *http.Request) {
	// Login logic here
	if strings.HasPrefix("/login/", r.URL.Path) {
		http.Error(w, "Invalid login option", http.StatusBadRequest)
		return
	}
	option := r.URL.Path[len("/login/"):]

	switch option {
	case "credentials":
		s.credentials(w, r)
	case "google":
		s.google(w, r)
	default:
		http.Error(w, "Invalid login option", http.StatusBadRequest)
	}
}

type GoogleLoginRequest struct {
	IdToken string `json:"idToken"`
}

func getResponse(reqString string) (map[string]interface{}, error) {
	client := &http.Client{} // Create HTTP request
	req, err := http.NewRequest("GET", reqString, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	json.Unmarshal([]byte(body), &data)
	return data, nil
}

func (s *RestServer) google(w http.ResponseWriter, r *http.Request) {
	in := &GoogleLoginRequest{}
	err := json.NewDecoder(r.Body).Decode(in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reqString := "https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + in.IdToken
	data, err := getResponse(reqString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 1. Check if email is verified by google
	// 2, If yes, check if user exists
	// 3. If yes, generate access token and set cookies
	// 4, Generate refresh token and put into refresh token db
	// 5. Return user info

	// 1. Check if email is verified by google
	if data["email_verified"] != "true" {
		http.Error(w, "Email not verified", http.StatusBadRequest)
		return
	}
	// (Set up constants)
	email := data["email"].(string)

	// 2. If yes, check if user exists
	user, err := s.restHelper.GetUserByEmail(r.Context(), email)
	if err == restapp_helper.ErrorUserNotFound {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Generate access token and set cookies
	accessContent := &AccessContent{UserMail: email, UserName: user.Name}
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
		SameSite: http.SameSiteLaxMode,
	})

	// 4. Generate refresh token and put into refresh token db
	tokenid, err := s.refreshHelper.SaveRefreshMeta(user.Id, RefreshTokenDuration)
	refreshContent := &RefreshContent{UserID: user.Id, TokenID: tokenid, UserMail: email}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	refreshToken, err := jwt.Sign(refreshContent, s.config.Secret, "Bearer ", RefreshTokenDuration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Return user info
	loginResponse := &LoginResponse{Name: user.Name, Mail: email, Avatar: user.Avatar, RefreshToken: refreshToken}

	err = json.NewEncoder(w).Encode(loginResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

type CredentialLoginRequest struct {
	Mail string `json:"mail"`
	Pass string `json:"pass"`
}

func (s *RestServer) credentials(w http.ResponseWriter, r *http.Request) {
	in := &CredentialLoginRequest{}
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
	user, err := s.restHelper.GetUserByEmail(r.Context(), in.Mail)
	if err == restapp_helper.ErrorUserNotFound {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if ok, err := bcrypt.Compare(in.Pass, user.Password.String); !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// 2. Generate access token and set cookies
	accessContent := &AccessContent{UserMail: in.Mail, UserName: user.Name}
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
		SameSite: http.SameSiteLaxMode,
	})

	// 3. Generate refresh token and put into refresh token db
	tokenid, err := s.refreshHelper.SaveRefreshMeta(user.Id, RefreshTokenDuration)
	refreshContent := &RefreshContent{UserID: user.Id, TokenID: tokenid, UserMail: in.Mail}
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
	loginResponse := &LoginResponse{Name: user.Name, Mail: in.Mail, Avatar: user.Avatar, RefreshToken: refreshToken}

	err = json.NewEncoder(w).Encode(loginResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}
