package main

import (
	"net/http"

	"github.com/chtozamm/annynotes-go/internal/auth"
	"github.com/chtozamm/annynotes-go/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (app *application) withAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := auth.ValidateJWT(r.Header)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		user, err := app.DB.GetUserByID(r.Context(), claims.UserID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		handler(w, r, user)
	}
}

// TODO: add logger middleware

// func (app *application) logger(handler http.Handler) http.HandlerFunc {
// }
