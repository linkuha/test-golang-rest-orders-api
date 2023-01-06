package order

import (
	"context"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/errs"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository/order"
)

type UseCase struct {
	repo order.Repository
}

func NewOrderUseCase(repo order.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetByID(ctx context.Context, orderID string) (*entity.Order, error) {
	res, err := uc.repo.Get(ctx, orderID)
	if err != nil {
		return nil, errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}
	return res, nil
}

func (uc *UseCase) GetAllByUserID(ctx context.Context, userID string) (*[]entity.Order, error) {
	res, err := uc.repo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}
	return res, nil
}

func (uc *UseCase) GetAllOrderProducts(ctx context.Context, orderID string) (*[]entity.OrderProductView, error) {
	res, err := uc.repo.GetProducts(ctx, orderID)
	if err != nil {
		return nil, errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}
	return res, nil
}

func (uc *UseCase) Create(ctx context.Context, order entity.Order) (string, error) {
	if err := order.Validate(); err != nil {
		return "", errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}

	res, err := uc.repo.Store(ctx, &order)
	if err != nil {
		return "", errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}
	return res, nil
}

func (uc *UseCase) AddProduct(ctx context.Context, p *entity.Product, op *entity.OrderProduct) error {
	if err := op.Validate(); err != nil {
		return errs.NewErrorWrapper(errs.Validation, err, "order product validation error")
	}

	if op.Amount > p.LeftInStock {
		return errs.NewErrorWrapper(errs.Logic, errs.LogicalError, "not enough amount in stock")
	}

	if err := uc.repo.AddProduct(ctx, op); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}
	return nil
}

func (uc *UseCase) Remove(ctx context.Context, id string) error {
	if err := uc.repo.Remove(ctx, id); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}
	return nil
}

func (uc *UseCase) Update(ctx context.Context, order entity.Order) error {
	if err := order.Validate(); err != nil {
		return errs.NewErrorWrapper(errs.Validation, err, "order validation error")
	}

	if err := uc.repo.Update(ctx, &order); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}
	return nil
}

func (uc *UseCase) RemoveProduct(ctx context.Context, orderID, productID string) error {
	if err := uc.repo.RemoveProduct(ctx, orderID, productID); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}
	return nil
}
