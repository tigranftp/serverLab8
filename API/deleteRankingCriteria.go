package API

import (
	"db_lab8/db"
	"db_lab8/types"
	"encoding/json"
	"io"
	"net/http"
)

func (a *API) handleDeleteRankingCriteria() http.HandlerFunc {
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
		var rc types.RankingCriteria
		err = json.Unmarshal(body, &rc)
		if err != nil {
			http.Error(writer, "can't close body", http.StatusInternalServerError)
			return
		}
		if rc.CriteriaName == "" {
			http.Error(writer, "can't add country empty with empty countryName", http.StatusInternalServerError)
			return
		}
		_, err = a.store.Exec(db.DeleteRankingCriteriaQuery, rc.SystemID, rc.CriteriaName, rc.SystemID, rc.CriteriaName)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}
