package API

import (
	"db_lab8/db"
	"db_lab8/types"
	"encoding/json"
	"io"
	"net/http"
)

func (a *API) handleCreateUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, role, err := a.GetIDAndRoleFromTokenAndRefreshTokenIfNeeded(writer, request)
		if err != nil {
			http.Error(writer, "You are not logged in. Sign In please", http.StatusBadRequest)
			return
		}
		if role != "admin" {
			http.Error(writer, "You are not admin and you have no right for this act.", http.StatusBadRequest)
			return
		}
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "can't read body", http.StatusBadRequest)
			return
		}
		err = request.Body.Close()
		if err != nil {
			http.Error(writer, "can't close body", http.StatusInternalServerError)
			return
		}
		var usr types.User
		err = json.Unmarshal(body, &usr)
		if err != nil {
			http.Error(writer, "can't close body", http.StatusInternalServerError)
			return
		}
		_, err = a.store.Exec(db.CreateUserQuery, usr.Name, usr.Username, generatePasswordHash(usr.Password), usr.Role)
		if err != nil {
			if err.Error() == "UNIQUE constraint failed: users.Username" {
				http.Error(writer, "Username is already in use. Try to use another one.", http.StatusBadGateway)
				return
			}
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}
