// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package sqlc

import ()

type User struct {
	ID       int64 `json:"id"`
	Name     string `json:"name"`
	Password string	`json:"password"`
}
