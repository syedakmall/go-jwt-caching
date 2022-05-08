package handler

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	db "github.com/syedakmall/malhttp/config"
	"github.com/syedakmall/malhttp/models"
)

var ctx context.Context = context.Background()

var UserHandler userHandler = userHandler{}

type userHandler struct {
}

func (us *userHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := db.Pg.Queries.ListUsers(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	u := models.Users{Users: users}
	j, _ := json.Marshal(u)
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func (us *userHandler) DeleteUserById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, "Invalid Id", http.StatusBadRequest)
		return
	}

	err = db.Pg.Queries.DeleteUser(ctx, int64(id))
	if err != nil {
		ErrorResponse(w, "User Not Found", http.StatusNotFound)
		return
	}
	err = db.MyRedis.DeleteCache(chi.URLParam(r, "id"), context.TODO())
	if err != nil {
		ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Deleted user with id " + strconv.Itoa(id)))

}

func (us *userHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		ErrorResponse(w, "Invalid Id", http.StatusBadRequest)
		return
	}

	u, err := db.Pg.Queries.GetUser(ctx, int64(id))

	if err != nil {
		ErrorResponse(w, "User Not Found", http.StatusNotFound)
		return
	}
	ctx := context.TODO()

	if err := db.MyRedis.SetCache(strconv.Itoa(id), ctx, u, time.Hour); err != nil {
		ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(u)
	w.Write(j)
}

func (us *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Header["Authorization"]
	if !ok {
		ErrorResponse(w, "No Token", http.StatusUnauthorized)
		return
	}
	token := r.Header.Get("Authorization")[7:]
	claims := Claims{}

	tkn, err := jwt.ParseWithClaims(token, &claims,
		func(t *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			ErrorResponse(w, "Invalid Token", http.StatusUnauthorized)
			return
		}
		ErrorResponse(w, "Error", http.StatusInternalServerError)
		return
	}

	if !tkn.Valid {
		ErrorResponse(w, "Invalid Token", http.StatusUnauthorized)
		return
	}

	u, err := db.Pg.Queries.GetUserByName(ctx, claims.Name)
	if err != nil {
		ErrorResponse(w, "User Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(u)
	w.Write(j)
}

func ErrorResponse(w http.ResponseWriter, message string, status int) {

	error := map[string]any{
		"message": message,
		"status":  status,
	}
	j, _ := json.Marshal(error)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}
