package product

import (
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/errs"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository/product"
)

type UseCase struct {
	repo product.Repository
}

func NewProductUseCase(repo product.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetByID(productID string) (*entity.Product, error) {
	res, err := uc.repo.Get(productID)
	if err != nil {
		return nil, errs.NewErrorWrapper(errs.Database, err, "error from product repo")
	}
	return res, nil
}

func (uc *UseCase) GetAll() (*[]entity.Product, error) {
	res, err := uc.repo.GetAll()
	if err != nil {
		return nil, errs.NewErrorWrapper(errs.Database, err, "error from product repo")
	}
	return res, nil
}

func (uc *UseCase) Create(product entity.Product) (string, error) {
	if err := product.Validate(); err != nil {
		return "", errs.NewErrorWrapper(errs.Validation, err, "product validation error")
	}

	res, err := uc.repo.Store(&product)
	if err != nil {
		return "", errs.NewErrorWrapper(errs.Database, err, "error from product repo")
	}
	return res, nil
}

func (uc *UseCase) CreateWithPrices(product entity.Product) (string, error) {
	if err := product.Validate(); err != nil {
		return "", errs.NewErrorWrapper(errs.Validation, err, "product validation error")
	}

	res, err := uc.repo.StoreWithPrices(&product)
	if err != nil {
		return "", errs.NewErrorWrapper(errs.Database, err, "error from product repo")
	}
	return res, nil
}

func (uc *UseCase) Remove(id string) error {
	if err := uc.repo.Remove(id); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from product repo")
	}
	return nil
}

func (uc *UseCase) Update(id string, input entity.ProductUpdateInput) error {
	if err := uc.repo.Update(id, &input); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from product repo")
	}
	return nil
}
