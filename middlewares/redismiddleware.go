package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	db "github.com/syedakmall/malhttp/config"
	"github.com/syedakmall/malhttp/sqlc"
)

func VerifyRedisCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		ctx := context.TODO()
		var wanted sqlc.User
		if err := db.MyRedis.GetCache(id, ctx, &wanted); err == nil {
			w.Header().Set("Content-Type", "application/json")
			j, _ := json.Marshal(wanted)
			w.Write(j)
			return
		}

		next.ServeHTTP(w, r)

	})
}
