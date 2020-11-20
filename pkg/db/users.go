package db

import "database/sql"

// User holds basic user info
type User struct {
	ID           int64
	Name         string
	PasswordHash string
}

// CreateUser creates new user and returns User object with
// id field populated
func CreateUser(tx *sql.Tx, u *User) (*User, error) {
	return u, tx.QueryRow(
		`insert into users(name, password_hash)
	values($1, $2)
	returning id`,
		u.Name,
		u.PasswordHash,
	).Scan(&u.ID)
}

// GetUser returns user by its id
func GetUser(tx *sql.Tx, id int64) (*User, error) {
	u := User{}
	return &u, tx.QueryRow(
		"select id, name, password_hash from users where id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Name,
		&u.PasswordHash,
	)
}
