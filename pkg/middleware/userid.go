package middleware

import (
	"context"
	"log"
	"net/http"
	"strconv"
)

// UserIDHeader is a header name for incoming user id
const UserIDHeader = "X-User-ID"

// ContextKey is a custom type for context key
type ContextKey string

// UserIDKey blah blah
const UserIDKey ContextKey = UserIDHeader

// UserIDMiddleware adds middleware to extract user id from request
func UserIDMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userIDStr := r.Header.Get(UserIDHeader)
			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				log.Printf("no valid %s from %s", UserIDHeader, r.RemoteAddr)
				return
			}
			ctx := context.WithValue(
				r.Context(),
				UserIDKey,
				userID,
			)

			f(w, r.WithContext(ctx))
		},
	)
}
