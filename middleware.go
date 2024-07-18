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

		user, err := app.DB.GetUserByEmail(r.Context(), claims.Email)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		handler(w, r, user)
	}
}

// func (app *application) logger(handler http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// logger.Logger.Info(r.Method, r.URL.Path, r.Header.Get("Referer"), r.RemoteAddr, r.Header.Get(""))
// 		log.Println(r.Method, r.URL.Path, r.Header.Get("Referer"), r.RemoteAddr, r.Header.Get(""))
// 		handler(w, r)
// 	}
// }
