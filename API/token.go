package API

import (
	"crypto/sha1"
	"db_lab8/db"
	"db_lab8/types"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"net/http"
	"time"
)

const (
	salt            = "asjhdjahsdjahsdas"
	signingKey      = "%*FG67G%f786^G%&()(&J*H)(_I*K{76534d5D"
	tokenTTL        = 5 * time.Second
	refreshTokenTTL = 5 * 24 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int64  `json:"user_id"`
	Role   string `json:"role"`
}

//func (a *API) handleGetToken() http.HandlerFunc {
//	return func(writer http.ResponseWriter, request *http.Request) {
//		ckc, err := request.Cookie("session_token")
//		if err != nil && !errors.Is(err, http.ErrNoCookie) {
//			http.Error(writer, "no cookie", http.StatusInternalServerError)
//			return
//		}
//		if err == nil {
//			_, err := a.ParseToken(ckc.Value)
//			if err == nil {
//				writer.WriteHeader(http.StatusOK)
//				return
//			}
//		}
//
//		body, err := io.ReadAll(request.Body)
//		if err != nil {
//			http.Error(writer, "can't read body", http.StatusBadRequest)
//			return
//		}
//		err = request.Body.Close()
//		if err != nil {
//			http.Error(writer, "can't close body", http.StatusInternalServerError)
//			return
//		}
//		var usr types.User
//		err = json.Unmarshal(body, &usr)
//		if err != nil {
//			http.Error(writer, "can't close body", http.StatusInternalServerError)
//			return
//		}
//		token, err := a.generateToken(usr.Username, usr.Password)
//		if err != nil {
//			http.Error(writer, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		http.SetCookie(writer, &http.Cookie{
//			Name:    "session_token",
//			Value:   token,
//			Expires: time.Now().Add(tokenTTL),
//		})
//		writer.WriteHeader(http.StatusOK)
//	}
//}

func (a *API) handleParseToken() http.HandlerFunc {
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
		ckc, err := request.Cookie("session_token")
		if err != nil {
			http.Error(writer, "no cookie", http.StatusInternalServerError)
			return
		}
		userID, role, err := a.ParseToken(ckc.Value)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(fmt.Sprintf("%v, %v", userID, role)))
	}
}

func (a *API) ParseToken(accessToken string) (int64, string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, claims.Role, nil
}

func (a *API) ParseRefreshToken(refreshToken string) (*types.User, error) {
	rows, err := a.store.Query(db.GetUserByRefreshTokenQuery, refreshToken)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var user types.User
	isThereAnyRow := rows.Next()
	if !isThereAnyRow {
		return nil, errors.New("no such refresh token")
	}
	err = rows.Scan(&user.Id, &user.Name, &user.Username, &user.Password, &user.RefreshToken, &user.RefreshTokenEAT, &user.Role)
	if err != nil {
		return nil, err
	}

	if time.Now().After(time.Unix(user.RefreshTokenEAT.Int64, 0)) {
		return nil, errors.New("refresh token expired")
	}
	return &user, nil
}

func (a *API) getUserByUserNameAndPassword(username, password string) (int64, string, error) {
	rows, err := a.store.Query(db.GetUserQuery, username, generatePasswordHash(password))
	if err != nil {
		return 0, "", err
	}
	defer rows.Close()
	var user types.User
	isThereAnyRow := rows.Next()
	if !isThereAnyRow {
		rows.Close()
		return 0, "", errors.New("login or password is incorrect")
	}
	err = rows.Scan(&user.Id, &user.Name, &user.Username, &user.Password, &user.RefreshToken, &user.RefreshTokenEAT, &user.Role)
	return user.Id, user.Role, err
}

func (a *API) generateTokensByCred(username, password string) (string, string, error) {
	userID, role, err := a.getUserByUserNameAndPassword(username, password)
	if err != nil {
		return "", "", err
	}
	return a.generateTokens(userID, role)
}

func (a *API) generateTokens(userID int64, role string) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userID,
		role,
	})

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}
	_, err = a.store.Exec(db.UpdateRefreshQuery, refreshToken, time.Now().Add(refreshTokenTTL).Unix(), userID)
	if err != nil {
		return "", "", err
	}
	ttk, err := token.SignedString([]byte(signingKey))
	return ttk, refreshToken, err
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
