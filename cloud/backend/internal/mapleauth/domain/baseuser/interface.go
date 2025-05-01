// github.com/Maple-Open-Tech/monorepo/cloud/backend/internal/mapleauth/domain/baseuser/interface.go
package baseuser

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository Interface for federatediam.
type Repository interface {
	Create(ctx context.Context, m *BaseUser) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*BaseUser, error)
	GetByEmail(ctx context.Context, email string) (*BaseUser, error)
	GetByVerificationCode(ctx context.Context, verificationCode string) (*BaseUser, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	DeleteByEmail(ctx context.Context, email string) error
	CheckIfExistsByEmail(ctx context.Context, email string) (bool, error)
	UpdateByID(ctx context.Context, m *BaseUser) error
	ListAll(ctx context.Context) ([]*BaseUser, error)
	CountByFilter(ctx context.Context, filter *BaseUserFilter) (uint64, error)
	ListByFilter(ctx context.Context, filter *BaseUserFilter) (*BaseUserFilterResult, error)
}
