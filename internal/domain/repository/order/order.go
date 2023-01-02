package order

import (
	"database/sql"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
)

type Repository interface {
	Get(id string) (*entity.Order, error)
	GetAllByUserID(userId string) (*[]entity.Order, error)
	GetProducts(id string) (*[]entity.OrderProductView, error)

	Store(order *entity.Order) (string, error)
	Update(order *entity.Order) error
	Remove(id string) error
	AddProduct(p *entity.OrderProduct) error
	RemoveProduct(orderID, productID string) error
}

func NewRepository(db *sql.DB) Repository {
	return newOrderPostgresRepository(db)
}
