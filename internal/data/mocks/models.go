package mocks

import (
	"github.com/zukofett/greenlight/internal/data"
)

func NewModels() data.Models {
	return data.Models{
		Movies:      &MovieModel{},
		Permissions: &PermissionModel{},
		Tokens:      &TokenModel{},
		Users:       &UserModel{},
	}
}
