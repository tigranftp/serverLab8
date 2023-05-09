package API

import (
	"net/http"
)

func (a *API) handleSignOut() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		c := &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(writer, c)

		c = &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(writer, c)

		writer.WriteHeader(http.StatusOK)
	}
}
