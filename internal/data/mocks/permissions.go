package mocks

import (
	"context"

	"github.com/zukofett/greenlight/internal/data"
)

var mockPermissions = data.Permissions{"movies:read", "movies:write"}

type PermissionModel struct{}

func (m PermissionModel) GetAllForUser(_ context.Context, userID int64) (data.Permissions, error) {
	switch userID {
	case 1:
		return mockPermissions, nil
	default:
		return nil, nil
	}
}

func (m PermissionModel) AddForUser(_ context.Context, _ int64, _ ...string) error {
	return nil
}
