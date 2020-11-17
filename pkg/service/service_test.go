package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"demoservice/pkg/middleware"
)

func TestTimeHandler(t *testing.T) {
	h := NewTimeHandler(10)
	srv := httptest.NewServer(http.HandlerFunc(middleware.UserIDMiddleware(h.Handle)))

	client := http.Client{Timeout: time.Millisecond * 10}
	req, err := http.NewRequest("get", srv.URL+"/time", nil)
	if err != nil {
		t.Fatalf("unexpected error in NewRequest: %v", err)
	}

	req.Header.Add(middleware.UserIDHeader, "42")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("unexpected error in client.Do: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status while calling /time: %d", resp.StatusCode)
	}

	value := TimeStruct{}
	if err := json.NewDecoder(resp.Body).Decode(&value); err != nil {
		t.Fatalf("failed to read json response from /time: %v", err)
	}

	if value.Date == "" {
		t.Error("date value is empty")
	}

	if value.Time == "" {
		t.Error("time value is empty")
	}

}

func BenchmarkTimeHandler(b *testing.B) {
	h := NewTimeHandler(10)
	srv := httptest.NewServer(http.HandlerFunc(middleware.UserIDMiddleware(h.Handle)))

	client := http.Client{Timeout: time.Millisecond * 100}

	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest("get", srv.URL+"/time", nil)
		if err != nil {
			b.Fatalf("unexpected error in NewRequest: %v", err)
		}

		req.Header.Add(middleware.UserIDHeader, "42")

		resp, err := client.Do(req)
		if err != nil {
			b.Fatalf("unexpected error in client.Do: %v", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			b.Fatalf("unexpected status while calling /time: %d", resp.StatusCode)
		}
	}
}
