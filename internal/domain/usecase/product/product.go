package product

import (
	"context"
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

func (uc *UseCase) GetByID(ctx context.Context, productID string) (*entity.Product, error) {
	res, err := uc.repo.Get(ctx, productID)
	if err != nil {
		return nil, errs.NewErrorWrapper(errs.Database, err, "error from product repo")
	}
	return res, nil
}

func (uc *UseCase) GetAll(ctx context.Context) (*[]entity.Product, error) {
	res, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, errs.NewErrorWrapper(errs.Database, err, "error from product repo")
	}
	return res, nil
}

func (uc *UseCase) Create(ctx context.Context, product entity.Product) (string, error) {
	if err := product.Validate(); err != nil {
		return "", errs.NewErrorWrapper(errs.Validation, err, "product validation error")
	}

	res, err := uc.repo.Store(ctx, &product)
	if err != nil {
		return "", errs.NewErrorWrapper(errs.Database, err, "error from product repo")
	}
	return res, nil
}

func (uc *UseCase) CreateWithPrices(ctx context.Context, product entity.Product) (string, error) {
	if err := product.Validate(); err != nil {
		return "", errs.NewErrorWrapper(errs.Validation, err, "product validation error")
	}

	res, err := uc.repo.StoreWithPrices(ctx, &product)
	if err != nil {
		return "", errs.NewErrorWrapper(errs.Database, err, "error from product repo")
	}
	return res, nil
}

func (uc *UseCase) Remove(ctx context.Context, id string) error {
	if err := uc.repo.Remove(ctx, id); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from product repo")
	}
	return nil
}

func (uc *UseCase) Update(ctx context.Context, id string, input entity.ProductUpdateInput) error {
	if err := uc.repo.Update(ctx, id, &input); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from product repo")
	}
	return nil
}
