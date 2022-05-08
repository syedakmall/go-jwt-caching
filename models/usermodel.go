package models

import (
	"github.com/syedakmall/malhttp/sqlc"
)

type Users struct {
	Users []sqlc.User `json:"users"`
}
