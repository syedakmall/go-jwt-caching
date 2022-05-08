package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	db "github.com/syedakmall/malhttp/config"
	route "github.com/syedakmall/malhttp/routes"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}

	db.Pg.Init()
	db.MyRedis.InitRedis()

	r.Route("/user", route.UserRoute)
	r.Route("/auth", route.AuthRoute)

	log.Println("🚀 port 8081")
	http.ListenAndServe(":8081", r)

}
