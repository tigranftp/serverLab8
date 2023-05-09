package API

import (
	"db_lab8/db"
	"db_lab8/types"
	"encoding/json"
	"io"
	"net/http"
)

func (a *API) handleAddUniversityRankingYear() http.HandlerFunc {
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
		var aury types.AddUniversityRankingYear
		err = json.Unmarshal(body, &aury)
		if err != nil {
			http.Error(writer, "error during unmarshal", http.StatusBadRequest)
			return
		}
		_, err = a.store.Exec(db.AddUniversityRankingYear, aury.UniversityName, aury.CriteriaName, aury.Year, aury.Score)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
	}
}
