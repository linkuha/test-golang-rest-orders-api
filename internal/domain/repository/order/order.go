package order

import (
	"context"
	"database/sql"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
)

type Repository interface {
	Get(ctx context.Context, id string) (*entity.Order, error)
	GetAllByUserID(ctx context.Context, userId string) (*[]entity.Order, error)
	GetProducts(ctx context.Context, id string) (*[]entity.OrderProductView, error)

	Store(ctx context.Context, order *entity.Order) (string, error)
	Update(ctx context.Context, order *entity.Order) error
	Remove(ctx context.Context, id string) error
	AddProduct(ctx context.Context, p *entity.OrderProduct) error
	RemoveProduct(ctx context.Context, orderID, productID string) error
}

func NewRepository(db *sql.DB) Repository {
	return newOrderPostgresRepository(db)
}
