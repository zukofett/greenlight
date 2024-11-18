package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/zukofett/greenlight/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
    ErrDuplicateEmail = errors.New("duplicate email")
)

type User struct {
    ID        int64     `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Password  password  `json:"-"`
    Activated bool      `json:"activated"`
    Version   int32     `json:"-"`
}

type password struct {
    plaintext *string
    hash      []byte
}


func (p *password) Set(plaintextPassword string) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
    if err != nil {
        return err
    }

    p.hash = hash
    p.plaintext = &plaintextPassword

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
    v.Check(len(password) <= 72, "password", "must not be more that 72 bytes long")

}

func ValidateUser(v *validator.Validator, user *User) {
    v.Check(user.Name != "", "name", "must be provided")
    v.Check(len(user.Name) <= 500, "password", "must not be more that 500 bytes long")

    ValidateEmail(v, user.Email)

    if user.Password.plaintext != nil {
        ValidatePasswordPlaintext(v, *user.Password.plaintext)
    }

    if user.Password.hash == nil {
        panic("missing password hash for user")
    }
}


type UserModel struct {
    DB *sql.DB
}

func (m UserModel) Insert(ctx context.Context, user *User) error {
    query := `
        INSERT INTO users (name, email, password_hash, activated) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id, created_at, version`

    args := []any{user.Name, user.Email, user.Password.hash, user.Activated}

    err:= m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
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

func (m UserModel) GetByEmail(ctx context.Context, email string) (*User, error) {
    query := `
        SELECT id, created_at, name, email, password_hash, activated, version
        FROM users
        WHERE email = $1`

    var user User
    
    err := m.DB.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
        &user.CreatedAt,
        &user.Name,
        &user.Email,
        &user.Password.hash,
        &user.Activated,
        &user.Version,
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

func (m UserModel) Update(ctx context.Context, user *User) error {
    query := `
        UPDATE users
        SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
        WHERE id = $5 AND version = $6
        RETURNING version`

    args := []any{
        user.Name,
        user.Email,
        user.Password.hash,
        user.Activated,
        user.ID,
        user.Version,
    }

    err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
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
