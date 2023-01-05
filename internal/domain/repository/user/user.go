package user

import (
	"context"
	"database/sql"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
)

type Repository interface {
	Get(ctx context.Context, id string) (*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	//GetFollowerIDs(id string, offset, limit int) ([]string, error)

	Store(ctx context.Context, user *entity.User) (string, error)
	Update(ctx context.Context, user *entity.User) error
	Remove(ctx context.Context, id string) error
	AddFollower(ctx context.Context, userID, followerID string) error
}

func NewRepository(db *sql.DB) Repository {
	return newUserPostgresRepository(db)
}
