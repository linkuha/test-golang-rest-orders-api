package order

import (
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

func (uc *UseCase) GetByID(orderID string) (*entity.Order, error) {
	res, err := uc.repo.Get(orderID)
	if err != nil {
		return nil, errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}
	return res, nil
}

func (uc *UseCase) GetAllByUserID(userID string) (*[]entity.Order, error) {
	res, err := uc.repo.GetAllByUserID(userID)
	if err != nil {
		return nil, errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}
	return res, nil
}

func (uc *UseCase) GetAllOrderProducts(orderID string) (*[]entity.OrderProductView, error) {
	res, err := uc.repo.GetProducts(orderID)
	if err != nil {
		return nil, errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}
	return res, nil
}

func (uc *UseCase) Create(order entity.Order) (string, error) {
	if err := order.Validate(); err != nil {
		return "", errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}

	res, err := uc.repo.Store(&order)
	if err != nil {
		return "", errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}
	return res, nil
}

func (uc *UseCase) AddProduct(p *entity.Product, op *entity.OrderProduct) error {
	if err := p.Validate(); err != nil {
		return errs.NewErrorWrapper(errs.Validation, err, "order validation error")
	}

	if op.Amount > p.LeftInStock {
		return errs.NewErrorWrapper(errs.Logic, errs.LogicalError, "not enough amount in stock")
	}

	if err := uc.repo.AddProduct(op); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}
	return nil
}

func (uc *UseCase) Remove(id string) error {
	if err := uc.repo.Remove(id); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}
	return nil
}

func (uc *UseCase) Update(order entity.Order) error {
	if err := order.Validate(); err != nil {
		return errs.NewErrorWrapper(errs.Validation, err, "order validation error")
	}

	if err := uc.repo.Update(&order); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}
	return nil
}

func (uc *UseCase) RemoveProduct(orderID, productID string) error {
	if err := uc.repo.RemoveProduct(orderID, productID); err != nil {
		return errs.NewErrorWrapper(errs.Database, err, "error from orders repo")
	}
	return nil
}
