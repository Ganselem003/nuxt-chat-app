package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ganselem/nuxt-chat-app/backend/internal/api"
	"github.com/ganselem/nuxt-chat-app/backend/internal/store"
	"github.com/ganselem/nuxt-chat-app/backend/internal/ws"
	"github.com/rs/cors"
)

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "file:chat.db?cache=shared&_fk=1"
	}

	db, err := store.NewSQLite(dsn)
	if err != nil {
		log.Fatal(err)
	}

	hub := ws.NewHub()
	go hub.Run()

	srv := api.NewServer(db)

	// token-iig ashiglan ws register
	srv.Router.HandleFunc("/ws", ws.WSHandler(hub, db))

	addr := ":8080"
	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler(srv.Router)

	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
