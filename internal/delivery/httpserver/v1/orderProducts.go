package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/usecase/order"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/usecase/product"
	"net/http"
)

// @Summary Add product to order
// @Security ApiKeyAuth
// @Tags order
// @Description Add product to order
// @ID order-product-add
// @Accept  json
// @Produce  json
// @Param input body entity.OrderProductView true "product data"
// @Param id path string true "Order ID"
// @Success 200 {object} statusResponse
// @Failure 400,403,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /orders/{id}/products [post]
func (ctrl *Controller) addOrderProduct(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		newErrorResponse(c, emptyParameterID)
		return
	}

	var input entity.OrderProductView
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, newJSONBindingErrorWrapper(err))
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	ouc := order.NewOrderUseCase(ctrl.repos.Orders)
	o, err := ouc.GetByID(ctrl.ctx, orderID)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	if o.UserID != userId {
		newErrorResponse(c, forbiddenError)
		return
	}

	puc := product.NewProductUseCase(ctrl.repos.Products)
	p, err := puc.GetByID(ctrl.ctx, input.ID)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	op := entity.OrderProduct{
		OrderID:   orderID,
		ProductID: input.ID,
		Amount:    input.Amount,
	}
	if err = ouc.AddProduct(ctrl.ctx, p, &op); err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{true})
}

// @Summary Get all order products
// @Security ApiKeyAuth
// @Tags order
// @Description get all order products
// @ID order-products-get-all
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /orders/{id}/products [get]
func (ctrl *Controller) getAllOrderProducts(c *gin.Context) {
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

	orders, err := uc.GetAllOrderProducts(ctrl.ctx, o.ID)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	newDataResponse(c, *orders)
}

// @Summary Delete product from order
// @Security ApiKeyAuth
// @Tags order
// @Description Delete product from order
// @ID order-product-delete
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Param productID path string true "Product ID"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /orders/{id}/products/{productID} [delete]
func (ctrl *Controller) deleteOrderProduct(c *gin.Context) {
	id := c.Param("id")
	productID := c.Param("productID")
	if id == "" || productID == "" {
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

	if err := uc.RemoveProduct(ctrl.ctx, id, productID); err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{true})
}
