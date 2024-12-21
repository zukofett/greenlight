package mocks

import (
	"context"
	"time"

	"github.com/zukofett/greenlight/internal/data"
)

type TokenModel struct {}


var mockToken = data.Token{
    Plaintext: "pa55word",
    Hash: []byte{},
    UserID: 1,
    Expiry: time.Now().Add(2*time.Minute),
    Scope: data.ScopeAuthentication,
}

func generateMockToken(_ int64, _ time.Duration, _ string) (*data.Token, error) {
	return &mockToken, nil
}

func (m TokenModel) New(ctx context.Context, userID int64, ttl time.Duration, scope string) (*data.Token, error) {
	token, err := generateMockToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(ctx, token)
	return token, err
}

func (m TokenModel) Insert(_ context.Context, token *data.Token) error {
	return nil
}

func (m TokenModel) DeleteAllForUser(_ context.Context, scope string, userID int64) error {
	return nil
}
