package api

import (
	"github.com/ganselem/nuxt-chat-app/backend/internal/store"
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	DB     *store.SQLiteStore
}

func NewServer(db *store.SQLiteStore) *Server {
	s := &Server{
		Router: mux.NewRouter(),
		DB:     db,
	}
	registerHandlers(s)
	return s
}
