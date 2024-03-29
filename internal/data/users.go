package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Skaifai/gophers-online-store/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	AnonymousUser     = &User{}
)

type RoleType string

const (
	ADMIN RoleType = "ADMIN"
	USER           = "USER"
	OWNER          = "OWNER"
)

type User struct {
	ID               int64     `json:"id"`
	Role             RoleType  `json:"role"`
	Username         string    `json:"username"`
	Email            string    `json:"email"`
	PhoneNumber      string    `json:"phone_number"`
	Password         password  `json:"-"`
	RegistrationDate time.Time `json:"registration_date"`
	Name             string    `json:"name"`
	Surname          string    `json:"surname"`
	DOB              time.Time `json:"date_of_birth"`
	Address          string    `json:"address"`
	AboutMe          string    `json:"about_me"`
	PictureURL       string    `json:"picture_url"`
	Activated        bool      `json:"activated"`
	Version          int       `json:"-"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type password struct {
	plaintext *string
	hash      []byte
}

type UserModel struct {
	DB *sql.DB
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 20, "name", "must not be more than 20 bytes long")
	ValidateEmail(v, user.Email)
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

func (u UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (username, email, phone_number, password_hash, name, surname, date_of_birth, address, about_me)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id, role, registration_date, picture_url, activated, version`

	args := []any{user.Username, user.Email, user.PhoneNumber, user.Password.hash,
		user.Name, user.Surname, user.DOB, user.Address, user.AboutMe}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Role, &user.RegistrationDate,
		&user.PictureURL, &user.Activated, &user.Version)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (u UserModel) GetById(id int64) (*User, error) {
	query := `
	SELECT id, role, username, email, phone_number, password_hash, registration_date, 
	       name, surname, date_of_birth, address, about_me, picture_url, 
	       activated
	FROM users
	WHERE id = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Role,
		&user.Username,
		&user.Email,
		&user.PhoneNumber,
		&user.Password.hash,
		&user.RegistrationDate,
		&user.Name,
		&user.Surname,
		&user.DOB,
		&user.Address,
		&user.AboutMe,
		&user.PictureURL,
		&user.Activated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u UserModel) GetByEmail(email string) (*User, error) {
	query := `
	SELECT id, role, username, email, phone_number, password_hash, registration_date, 
	       name, surname, date_of_birth, address, about_me, picture_url, 
	       activated
	FROM users
	WHERE email = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Role,
		&user.Username,
		&user.Email,
		&user.PhoneNumber,
		&user.Password.hash,
		&user.RegistrationDate,
		&user.Name,
		&user.Surname,
		&user.DOB,
		&user.Address,
		&user.AboutMe,
		&user.PictureURL,
		&user.Activated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u UserModel) Update(user *User) error {
	query := `
	UPDATE users
	SET role = $1, username = $2, email = $3, phone_number = $4, 
	    password_hash = $5, registration_date = $6, name = $7, surname = $8, date_of_birth = $9, 
	    address = $10, about_me = $11, picture_url = $12, activated = $13, 
	    version = version + 1
	WHERE id = $14
	RETURNING version`

	args := []any{
		user.Role,
		user.Username,
		user.Email,
		user.PhoneNumber,
		user.Password.hash,
		user.RegistrationDate,
		user.Name,
		user.Surname,
		user.DOB,
		user.Address,
		user.AboutMe,
		user.PictureURL,
		user.Activated,
		user.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (u UserModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM users
		WHERE id = $1`

	result, err := u.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
