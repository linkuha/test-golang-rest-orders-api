package profile

import (
	"database/sql"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
)

type Repository interface {
	GetByUserID(userID string) (*entity.Profile, error)

	Store(profile *entity.Profile) (int, error)
	Update(profile *entity.Profile) error
	RemoveByUserID(userID string) error
}

func NewRepository(db *sql.DB) Repository {
	return newProfilePostgresRepository(db)
}
