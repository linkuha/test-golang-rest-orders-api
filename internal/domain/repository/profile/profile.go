package profile

import (
	"context"
	"database/sql"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
)

type Repository interface {
	GetByUserID(ctx context.Context, userID string) (*entity.Profile, error)

	Store(ctx context.Context, profile *entity.Profile) (int, error)
	Update(ctx context.Context, profile *entity.Profile) error
	RemoveByUserID(ctx context.Context, userID string) error
}

func NewRepository(db *sql.DB) Repository {
	return newProfilePostgresRepository(db)
}
