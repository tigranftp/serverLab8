package API

import (
	"db_lab8/types"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

func (a *API) handleSignIn() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _, err := a.GetIDAndRoleFromTokenAndRefreshTokenIfNeeded(writer, request)
		if err == nil {
			writer.WriteHeader(http.StatusOK)
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
		token, refreshToken, err := a.generateTokensByCred(usr.Username, usr.Password)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		setTokenCookies(writer, token, refreshToken)
		writer.WriteHeader(http.StatusOK)
	}
}

func (a *API) GetIDAndRoleFromTokenAndRefreshTokenIfNeeded(writer http.ResponseWriter, request *http.Request) (int64, string, error) {
	// пытаемся спарсить из кук токен сессии
	ckc, err := request.Cookie("session_token")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		return 0, "", err
	}
	if err == nil {
		userID, role, err := a.ParseToken(ckc.Value)
		if err == nil {
			return userID, role, nil
		}
	}

	// пытаемся спарсить из кук токен рефреша, если с токеном сессии плохо
	ckc, err = request.Cookie("refresh_token")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		return 0, "", err
	}
	if err != nil {
		return 0, "", err
	}
	usr, err := a.ParseRefreshToken(ckc.Value)
	if err != nil {
		return 0, "", err
	}
	token, refreshToken, err := a.generateTokens(usr.Id, usr.Role)
	if err != nil {
		return 0, "", err
	}
	setTokenCookies(writer, token, refreshToken)
	return usr.Id, usr.Role, nil
}

func setTokenCookies(writer http.ResponseWriter, token, refreshToken string) {
	http.SetCookie(writer, &http.Cookie{
		Name:    "session_token",
		Value:   token,
		Expires: time.Now().Add(tokenTTL),
	})
	http.SetCookie(writer, &http.Cookie{
		Name:    "refresh_token",
		Value:   refreshToken,
		Expires: time.Now().Add(refreshTokenTTL),
	})
}
