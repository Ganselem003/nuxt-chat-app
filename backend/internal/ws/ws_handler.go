package ws

import (
	"net/http"

	auth "github.com/ganselem/nuxt-chat-app/backend/internal/auth"
	storepkg "github.com/ganselem/nuxt-chat-app/backend/internal/store"
)

func WSHandler(h *Hub, store *storepkg.SQLiteStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		uid, err := auth.ParseToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		var u storepkg.User
		if err := store.DB.First(&u, uid).Error; err != nil {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}

		ServeWS(h, w, r, u.Username, u.ID)
	}
}
