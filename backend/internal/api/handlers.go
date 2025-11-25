package api

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"

	"github.com/ganselem/nuxt-chat-app/backend/internal/auth"
	"github.com/ganselem/nuxt-chat-app/backend/internal/store"
)

type registerReq struct{ Username, Password string }
type loginReq struct{ Username, Password string }

func registerHandlers(s *Server) {
	s.Router.HandleFunc("/api/auth/register", func(w http.ResponseWriter, r *http.Request) {
		var req registerReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		hash, err := auth.HashPassword(req.Password)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		u := store.User{Username: req.Username, Password: hash}
		if err := s.DB.DB.Create(&u).Error; err != nil {
			http.Error(w, "could not create", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}).Methods("POST")

	s.Router.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		var req loginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		var u store.User
		if err := s.DB.DB.Where("username = ?", req.Username).First(&u).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, "invalid credentials", http.StatusUnauthorized)
				return
			}
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		if err := auth.CheckPassword(u.Password, req.Password); err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		tok, err := auth.CreateToken(u.ID)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		res := map[string]string{"token": tok}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}).Methods("POST")
}
