package profile

import (
	"database/sql"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
)

type Reader interface {
	GetByUserID(userID string) (*entity.Profile, error)
}

type Writer interface {
	Store(profile *entity.Profile) (int, error)
	Update(profile *entity.Profile) error
	RemoveByUserID(userID string) error
}

type Repository interface {
	Reader
	Writer
}

func NewRepository(db *sql.DB) Repository {
	return newProfilePostgresRepository(db)
}
