package service

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"demoservice/pkg/middleware"
)

const (
	timeFormat = "15:04:05"
	dateFormat = "2006.01.02"
)

// TimeHandler gives us current time on request
type TimeHandler struct {
	tzOffset uint8
}

// NewTimeHandler initializes TimeHandler object
func NewTimeHandler(tzOffset uint8) *TimeHandler {
	return &TimeHandler{tzOffset: tzOffset}
}

// TimeStruct just struct
type TimeStruct struct {
	Date string
	Time string
}

// AuthorizedUserID is the best user
const AuthorizedUserID = 42

func checkAuthOk(r *http.Request) bool {
	v := r.Context().Value(middleware.UserIDKey)
	userID, ok := v.(int)
	if !ok || userID < 1 {
		return false
	}

	if userID == AuthorizedUserID {
		return true
	}

	return false
}

// Handle handles http request
func (th *TimeHandler) Handle(w http.ResponseWriter, r *http.Request) {

	if !checkAuthOk(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	outTime := time.Now().UTC().Add(time.Hour * time.Duration(th.tzOffset))
	out := TimeStruct{
		Date: outTime.Format(dateFormat),
		Time: outTime.Format(timeFormat),
	}

	b, err := json.Marshal(&out)
	if err != nil {
		log.Printf("failed to marshal time: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}
