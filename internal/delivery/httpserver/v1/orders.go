package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/usecase/order"
	"net/http"
)

// @Summary Create order
// @Security ApiKeyAuth
// @Tags order
// @Description Create order
// @ID order-create
// @Accept  json
// @Produce  json
// @Param input body entity.Order true "order data"
// @Success 200 {string} string "id"
// @Failure 400,403,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /orders [post]
func (ctrl *Controller) createOrder(c *gin.Context) {
	var input entity.Order

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, newJSONBindingErrorWrapper(err))
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	if input.UserID != userId {
		newErrorResponse(c, forbiddenError)
		return
	}

	uc := order.NewOrderUseCase(ctrl.repos.Orders)
	id, err := uc.Create(ctrl.ctx, input)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Get order
// @Security ApiKeyAuth
// @Tags order
// @Description get order
// @ID order-get
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} entity.Order
// @Failure 400,403,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /orders/{id} [get]
func (ctrl *Controller) getOrderByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newErrorResponse(c, emptyParameterID)
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	uc := order.NewOrderUseCase(ctrl.repos.Orders)
	o, err := uc.GetByID(ctrl.ctx, id)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	if o.UserID != userId {
		newErrorResponse(c, forbiddenError)
		return
	}

	c.JSON(http.StatusOK, o)
}

// @Summary Get all orders
// @Security ApiKeyAuth
// @Tags order
// @Description get all orders
// @ID order-get-all
// @Accept  json
// @Produce  json
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /orders [get]
func (ctrl *Controller) getAllOrders(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	uc := order.NewOrderUseCase(ctrl.repos.Orders)
	orders, err := uc.GetAllByUserID(ctrl.ctx, userId)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	newDataResponse(c, *orders)
}

// @Summary Update order
// @Security ApiKeyAuth
// @Tags order
// @Description update order
// @ID order-update
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Param input body entity.Order true "order updating data"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /orders/{id} [put]
func (ctrl *Controller) updateOrderByID(c *gin.Context) {
	var input entity.Order

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, newJSONBindingErrorWrapper(err))
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	id := c.Param("id")
	if id == "" {
		newErrorResponse(c, err)
		return
	}
	input.ID = id

	uc := order.NewOrderUseCase(ctrl.repos.Orders)

	o, err := uc.GetByID(ctrl.ctx, id)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	if o.UserID != userId {
		newErrorResponse(c, forbiddenError)
		return
	}

	if err = uc.Update(ctrl.ctx, input); err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{true})
}

// @Summary Delete order
// @Security ApiKeyAuth
// @Tags order
// @Description delete product
// @ID order-delete
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /orders/{id} [delete]
func (ctrl *Controller) deleteOrderByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newErrorResponse(c, emptyParameterID)
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	uc := order.NewOrderUseCase(ctrl.repos.Orders)

	o, err := uc.GetByID(ctrl.ctx, id)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	if o.UserID != userId {
		newErrorResponse(c, forbiddenError)
		return
	}

	if err = uc.Remove(ctrl.ctx, id); err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{true})
}
