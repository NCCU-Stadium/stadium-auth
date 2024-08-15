package app

import (
	"auth-service/bcrypt"
	"auth-service/internal/helper"
	"encoding/json"
	"net/http"
)

type RegisterRequest struct {
	Mail   string `json:"email"`
	Pass   string `json:"password"`
	Phone  string `json:"phone"`
	Role   string `json:"role"` // "user", "admin" or "coach"
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
	Gender string `json:"gender,omitempty"` // "F" or "M"
	Birth  string `json:"birth,omitempty"`  // "YYYY-MM-DD"
}

type RegisterResponse struct {
	Mail string `json:"mail"`
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	in := &RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if in.Mail == "" || in.Pass == "" || in.Phone == "" || in.Role == "" || in.Name == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
	}

	exist, err := s.helper.IsUserExist(r.Context(), in.Mail)
	if exist {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	if in.Role != "user" && in.Role != "admin" && in.Role != "coach" {
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	encrypted, err := bcrypt.Encrypt(in.Pass, BcryptCost)
	regu := helper.RUserReq{
		Mail:   in.Mail,
		Pass:   encrypted,
		Role:   in.Role,
		Phone:  in.Phone,
		Name:   in.Name,
		Avatar: in.Avatar,
		Gender: in.Gender,
		Birth:  in.Birth,
	}

	err = s.helper.RegisterUser(r.Context(), regu)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := RegisterResponse{Mail: in.Mail}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}
