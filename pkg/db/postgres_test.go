package db

import (
	"context"
	"testing"
)

func TestConnect(t *testing.T) {

	conn, err := Connect(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error in Connect: %v", err)
	}

	if conn == nil {
		t.Fatal("Returned connection is nil")
	}
}
