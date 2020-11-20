package db

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreateGetUser(t *testing.T) {
	conn, err := Connect(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error in Connect: %v", err)
	}

	tx, err := conn.Begin()
	if err != nil {
		t.Fatalf("Unexpected error in Begin: %v", err)
	}

	defer tx.Rollback()

	u := User{
		Name:         "Васисуалий",
		PasswordHash: "DEADBEEF",
	}

	_, err = CreateUser(tx, &u)
	if err != nil {
		t.Fatalf("Unexpected error in CreateUser: %v", err)
	}

	if u.ID == 0 {
		t.Error("Expected user id to be not zero")
	}

	u2, err := GetUser(tx, u.ID)
	if err != nil {
		t.Fatalf("Unexpected error in GetUser: %v", err)
	}

	if diff := cmp.Diff(&u, u2); diff != "" {
		t.Fatalf("Saved and loaded users do not match:\n%s", diff)
	}
}
