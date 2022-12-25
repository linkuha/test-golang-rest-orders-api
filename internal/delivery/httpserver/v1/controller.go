package v1

import (
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository"
)

type Controller struct {
	repos repository.Repository
}

func NewController(repos repository.Repository) *Controller {
	return &Controller{repos}
}
