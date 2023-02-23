package data

import (
	"database/sql"
	"errors"
	"time"
)

// Define a custom ErrDuplicateEmail error
var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

// Define a string type for the UserRole.
type UserRoleType string

// Define a custom enum for the User Role fields. The user can either be a "USER" or an "ADMIN".
const (
	ADMIN UserRoleType = "admin"
	USER               = "user"
)

type User struct {
	ID               int64        `json:"id"`
	Role             UserRoleType `json:"role"`
	Username         string       `json:"username"`
	Email            string       `json:"email"`
	PhoneNumber      string       `json:"phone_number"`
	Password         password     `json:"-"`
	RegistrationDate time.Time    `json:"registration_date"`
	Profile          Profile      `json:"profile"`
	Activated        bool         `json:"activated"`
	Version          int          `json:"-"`
}

// Create a UserModel struct which wraps the connection pool.
type UserModel struct {
	DB *sql.DB
}

// Create a custom password type which is a struct containing the plaintext and hashed
// versions of the password for a user. The plaintext field is a *pointer* to a string,
// so that we're able to distinguish between a plaintext password not being present in
// the struct at all, versus a plaintext password which is the empty string "".
type password struct {
	plaintext *string
	hash      []byte
}
