package order

import (
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository/order"
)

type UseCase struct {
	repo order.Repository
}

func NewOrderUseCase(repo order.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetByID(orderID string) (*entity.Order, error) {
	return uc.repo.Get(orderID)
}

func (uc *UseCase) GetAllByUserID(userID string) (*[]entity.Order, error) {
	return uc.repo.GetAllByUserID(userID)
}

func (uc *UseCase) GetAllOrderProducts(orderID string) (*[]entity.OrderProductView, error) {
	return uc.repo.GetProducts(orderID)
}

func (uc *UseCase) Create(order entity.Order) (int, error) {
	if err := order.Validate(); err != nil {
		return 0, err
	}

	return uc.repo.Store(&order)
}

func (uc *UseCase) AddProduct(p entity.OrderProduct) error {
	if err := p.Validate(); err != nil {
		return err
	}

	return uc.repo.AddProduct(&p)
}

func (uc *UseCase) Remove(id string) error {
	return uc.repo.Remove(id)
}

func (uc *UseCase) Update(order entity.Order) error {
	if err := order.Validate(); err != nil {
		return err
	}

	return uc.repo.Update(&order)
}

func (uc *UseCase) RemoveProduct(orderID, productID string) error {
	return uc.repo.RemoveProduct(orderID, productID)
}
