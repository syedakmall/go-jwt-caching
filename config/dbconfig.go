package db

import (
	"database/sql"
	"fmt"
	"log"

	"os"

	_ "github.com/lib/pq"

	"github.com/syedakmall/malhttp/sqlc"
)

var Pg database = database{}

type database struct {
	Db      *sql.DB
	Queries *sqlc.Queries
}

func (db *database) Init() {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DBNAME"))

	Db, err := sql.Open("postgres", psqlInfo)
	db.Db = Db
	if err != nil {
		log.Fatal(err)
	}

	db.Queries = sqlc.New(Db)
}
