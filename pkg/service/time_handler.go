package service

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
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

// Handle handles http request
func (th *TimeHandler) Handle(w http.ResponseWriter, r *http.Request) {
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
