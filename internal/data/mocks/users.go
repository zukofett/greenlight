package mocks

import (
	"context"
	"time"

	"github.com/zukofett/greenlight/internal/data"
)

type UserModel struct{}

var mockUser = data.User{
	ID:        1,
	Name:      "Alice",
	Email:     "alice@examle.com",
	CreatedAt: time.Now(),
	Activated: false,
}

func (m *UserModel) Insert(_ context.Context, user *data.User) error {
	switch user.Email {
	case "dupe@examle.com":
		return data.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) GetByEmail(_ context.Context, email string) (*data.User, error) {
	switch email {
	case "alice@examle.com":
		return &mockUser, nil
	default:
		return nil, data.ErrRecordNotFound
	}
}

func (m *UserModel) Update(_ context.Context, user *data.User) error {
	switch {
	case user.Email == "dupe@email.com":
		return data.ErrDuplicateEmail
	case user.ID != 1:
		return data.ErrEditConflict
	default:
		return nil
	}
}

func (m UserModel) GetForToken(ctx context.Context, tokenScope string, tokenPlaintext string) (*data.User, error) {
	return &mockUser, nil
}
