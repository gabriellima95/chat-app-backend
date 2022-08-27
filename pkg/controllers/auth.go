package controllers

import (
	"net/http"

	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(handlerFN func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Authentication error: token not found", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			http.Error(w, "Authentication error: "+err.Error(), http.StatusUnauthorized)
			return
		}
		handlerFN(w, r)
	}
}
