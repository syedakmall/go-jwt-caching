package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/syedakmall/malhttp/handler"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Header["Authorization"]
		if !ok {
			handler.ErrorResponse(w, "No Token", http.StatusUnauthorized)
			return
		}
		token := r.Header.Get("Authorization")[7:]
		claims := handler.Claims{}
		tkn, err := jwt.ParseWithClaims(token, &claims,
			func(t *jwt.Token) (interface{}, error) {
				return handler.JwtKey, nil
			})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				handler.ErrorResponse(w, "Invalid Token", http.StatusUnauthorized)
				return
			}
			handler.ErrorResponse(w, "Error", http.StatusInternalServerError)
			return
		}

		if !tkn.Valid {
			handler.ErrorResponse(w, "Invalid Token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
