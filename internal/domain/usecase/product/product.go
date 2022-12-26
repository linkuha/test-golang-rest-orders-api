package product

import (
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository/product"
)

type UseCase struct {
	repo product.Repository
}

func NewProductUseCase(repo product.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetByID(productID string) (*entity.Product, error) {
	return uc.repo.Get(productID)
}

func (uc *UseCase) GetAll() (*[]entity.Product, error) {
	return uc.repo.GetAll()
}

func (uc *UseCase) Create(product entity.Product) (string, error) {
	if err := product.Validate(); err != nil {
		return "", err
	}

	return uc.repo.Store(&product)
}

func (uc *UseCase) CreateWithPrices(product entity.Product, prices []entity.Price) (string, error) {
	if err := product.Validate(); err != nil {
		return "", err
	}

	for _, price := range prices {
		if err := price.Validate(); err != nil {
			return "", err
		}
	}

	return uc.repo.Store(&product)
}

func (uc *UseCase) Remove(id string) error {
	return uc.repo.Remove(id)
}

func (uc *UseCase) Update(id string, input entity.ProductUpdateInput) error {
	return uc.repo.Update(id, &input)
}
