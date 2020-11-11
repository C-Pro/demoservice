package service

import (
	"context"
	"log"
	"net/http"
)

// Run starts our cool service
func Run(ctx context.Context) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, webservice!"))
	})

	h := NewTimeHandler(10)
	http.HandleFunc("/time", h.Handle)

	s := &http.Server{Addr: "0.0.0.0:8080"}
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
