package user

import (
	"database/sql"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
)

type Repository interface {
	Get(id string) (*entity.User, error)
	GetByUsername(username string) (*entity.User, error)
	//GetFollowerIDs(id string, offset, limit int) ([]string, error)

	Store(user *entity.User) (string, error)
	Update(user *entity.User) error
	Remove(id string) error
	AddFollower(userID, followerID string) error
}

func NewRepository(db *sql.DB) Repository {
	return newUserPostgresRepository(db)
}
