package product

import (
	"database/sql"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
)

type Reader interface {
	Get(id string) (*entity.Product, error)
	GetAll() (*[]entity.Product, error)
	GetPrices(id string) (*[]entity.Price, error)
}

type Writer interface {
	Store(product *entity.Product) (string, error)
	StoreWithPrices(product *entity.Product, prices []entity.Price) (string, error)
	Update(id string, input *entity.ProductUpdateInput) error
	Remove(id string) error
	AddPrice(productId string, price *entity.Price) error
}

type Repository interface {
	Reader
	Writer
}

func NewRepository(db *sql.DB) Repository {
	return newProductPostgresRepository(db)
}
