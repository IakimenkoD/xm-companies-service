package api

import (
	"encoding/json"
	"github.com/IakimenkoD/xm-companies-service/internal/model"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

func (srv *Server) signIn(w http.ResponseWriter, r *http.Request) {

	creds := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		respondError(w, err)
		return
	}

	expectedPassword, ok := users[creds.Login]

	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(256 * time.Minute)
	claims := &model.Claims{
		Username: creds.Login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(srv.cfg.API.JWTKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//TODO just write token?
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
