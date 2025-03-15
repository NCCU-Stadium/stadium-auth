package app

import (
	"auth-service/internal/helper"
	"auth-service/jwt"
	"encoding/json"
	"net/http"
)

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (s *Server) Refresh(w http.ResponseWriter, r *http.Request) {
	in := &RefreshRequest{}
	err := json.NewDecoder(r.Body).Decode(in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 1. Decode refresh token
	// 2. Get refresh meta from DB and verify refresh meta
	// 3. Generate new refresh token
	// 4. Generate new access token (and set cookie)
	// 5. Return the response

	// 1. Decode the refresh token
	content, err := jwt.Parse(in.RefreshToken, s.config.Secret, "Bearer ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if isExpired, err := jwt.IsExpired(content); isExpired || err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	refreshContent := &RefreshContent{} // refreshContent is a struct containing the decoded refresh token
	refreshContent.ToDomain(content)    // Convert the content to RefreshContent struct

	// 2. Get the refresh meta and verify refresh meta
	_, err = s.refreshHelper.GetRefreshMeta(refreshContent.TokenID)
	if err == helper.ErrorTokenUsed || err == helper.ErrorTokenNotFound {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// 3. Generate new refresh token
	tokenid, err := s.refreshHelper.SaveRefreshMeta(refreshContent.UserRole, RefreshTokenDuration)
	newRefreshContent := &RefreshContent{UserRole: refreshContent.UserRole, TokenID: tokenid, UserMail: refreshContent.UserMail}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newRefreshToken, err := jwt.Sign(newRefreshContent, s.config.Secret, "Bearer ", RefreshTokenDuration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Generate new access token and set cookie
	newAccessContent := &AccessContent{
		UserMail: refreshContent.UserMail,
		UserRole: refreshContent.UserRole,
	}
	newAccessToken, err := jwt.Sign(newAccessContent, s.config.Secret, "Bearer ", AccessTokenDuration)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.SetCookie(w, &http.Cookie{ // Set cookie
		Name:     "accessToken",
		Value:    newAccessToken,
		MaxAge:   int(AccessTokenDuration),
		Path:     "/",
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	// 5. Return the response
	user, err := s.helper.GetUserByMail(r.Context(), refreshContent.UserMail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res := &LoginResponse{Mail: user.Mail, Role: user.Role, RefreshToken: newRefreshToken}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}
