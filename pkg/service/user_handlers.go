package service

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	database "demoservice/pkg/db"
)

// UsersService has users handlers
type UsersService struct {
	db *sql.DB
}

// NewUsersService initializes UsersService object
func NewUsersService(db *sql.DB) *UsersService {
	return &UsersService{db: db}
}

// UsersHandler is a service-level demuxer
func (u *UsersService) UsersHandler(w http.ResponseWriter, r *http.Request) {
	ok, err := regexp.MatchString("/users/[0-9]+", r.URL.Path)
	if err != nil {
		log.Printf("failed to parse path: %s, %v", r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if ok && r.Method == "GET" {
		u.GetUser(w, r)
	}
}

// GetUser handles GET /users/{id} http request
func (u *UsersService) GetUser(w http.ResponseWriter, r *http.Request) {

	if !checkAuthOk(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		log.Printf("wrong path: %s", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		log.Printf("wrong user id format: %s", parts[2])
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := database.GetShortTimeoutContext(r.Context())
	defer cancel()

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("failed to create transaction: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer tx.Rollback()

	user, err := database.GetUser(tx, id)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		if ctx.Err() != nil {
			err = ctx.Err()
		}
		log.Printf("failed to get user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(user)
	if err != nil {
		log.Printf("failed to marshal user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

// CreateUser handles POST /users http request
func (u *UsersService) CreateUser(w http.ResponseWriter, r *http.Request) {

	if !checkAuthOk(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user := database.User{}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("wrong user json format: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx, cancel := database.GetShortTimeoutContext(r.Context())
	defer cancel()

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("failed to create transaction: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer tx.Rollback()

	userResult, err := database.CreateUser(tx, &user)
	if err != nil {
		if ctx.Err() != nil {
			err = ctx.Err()
		}
		log.Printf("failed to create user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(userResult)
	if err != nil {
		log.Printf("failed to marshal user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("failed to commit transaction: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
