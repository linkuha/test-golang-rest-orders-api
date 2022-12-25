package repository

import (
	"database/sql"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository/order"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository/product"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository/profile"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository/user"
)

type Repository struct {
	Orders   order.Repository
	Products product.Repository
	Users    user.Repository
	Profiles profile.Repository
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Orders:   order.NewRepository(db),
		Products: product.NewRepository(db),
		Users:    user.NewRepository(db),
		Profiles: profile.NewRepository(db),
	}
}
