package v1

import (
	"context"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository"
)

type Controller struct {
	ctx   context.Context
	repos repository.Repository
}

func NewController(ctx context.Context, repos repository.Repository) *Controller {
	return &Controller{ctx, repos}
}
