package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	db "github.com/syedakmall/malhttp/config"
	"github.com/syedakmall/malhttp/sqlc"
	"golang.org/x/crypto/bcrypt"
)

type SignInDto struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

var AuthHandler authHandler = authHandler{}

type authHandler struct {
}

func (h *authHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var createUserDto sqlc.CreateUserParams
	err := json.NewDecoder(r.Body).Decode(&createUserDto)
	defer r.Body.Close()
	if err != nil {
		ErrorResponse(w, "Error Decoding", http.StatusInternalServerError)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserDto.Password), bcrypt.DefaultCost)
	if err != nil {
		ErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}
	createUserDto.Password = string(hashedPassword)
	user, err := db.Pg.Queries.CreateUser(ctx, createUserDto)

	if err != nil {
		ErrorResponse(w, "Error Creating User", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(user)
	w.Write(j)
}

func (h *authHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var dto SignInDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		ErrorResponse(w, "Bad Request", http.StatusBadRequest)
		return
	}

	u, err := db.Pg.Queries.GetUserByName(ctx, dto.Name)
	if err != nil {
		ErrorResponse(w, "User Not Found", http.StatusNotFound)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(dto.Password))
	if err != nil {
		ErrorResponse(w, "Wrong Name Or Password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Minute.Round(10))
	claims := Claims{
		Name: dto.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    r.URL.Host,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		ErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Write([]byte(tokenString))
}

func (h *authHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
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
		ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if !tkn.Valid {
		ErrorResponse(w, "Invalid Token", http.StatusUnauthorized)
		return
	}

	if time.Until(time.Unix(claims.ExpiresAt, 0)) > 30*time.Second {
		ErrorResponse(w, "Bad request", http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := refreshToken.SignedString(JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(tokenString))
}
