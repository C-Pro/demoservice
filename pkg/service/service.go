package service

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"demoservice/pkg/middleware"
)

// Run starts our cool service
func Run(ctx context.Context, db *sql.DB) error {

	mux := http.NewServeMux()

	h := NewTimeHandler(10)
	mux.HandleFunc("/time", h.Handle)

	u := NewUsersService(db)
	mux.HandleFunc("/users/", u.UsersHandler)
	mux.HandleFunc("/users", u.CreateUser)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, webservice!"))
	})

	s := &http.Server{Addr: "0.0.0.0:8080",
		Handler: middleware.UserIDMiddleware(mux.ServeHTTP),
	}
	ch := make(chan error)

	go func(s *http.Server, ch chan error) {
		ch <- s.ListenAndServe()
	}(s, ch)

	select {
	case err := <-ch:
		log.Printf("server returned and error: %v", err)
		return err
	case <-ctx.Done():
		log.Println("context is canceled")
		return s.Close()
	}
}
