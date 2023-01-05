package product

import (
	"context"
	"database/sql"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
)

type Repository interface {
	Get(ctx context.Context, id string) (*entity.Product, error)
	GetAll(ctx context.Context) (*[]entity.Product, error)
	GetPrices(ctx context.Context, id string) (*[]entity.Price, error)

	Store(ctx context.Context, product *entity.Product) (string, error)
	StoreWithPrices(ctx context.Context, product *entity.Product) (string, error)
	Update(ctx context.Context, id string, input *entity.ProductUpdateInput) error
	Remove(ctx context.Context, id string) error
	AddPrice(ctx context.Context, productId string, price *entity.Price) error
}

func NewRepository(db *sql.DB) Repository {
	return newProductPostgresRepository(db)
}
