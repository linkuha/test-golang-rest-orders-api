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
// @Param input body entity.Order true "order data"
// @Success 200 {string} string "id"
// @Failure 400,403,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /orders/:id/products [post]
func (ctrl *Controller) addOrderProduct(c *gin.Context) {
	var input entity.OrderProduct

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	uc := order.NewOrderUseCase(ctrl.repos.Orders)
	o, err := uc.GetByID(input.OrderID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	if o.UserID != userId {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	puc := product.NewProductUseCase(ctrl.repos.Products)
	p, err := puc.GetByID(input.ProductID)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	if err := uc.AddProduct(p, &input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /orders/:id/products [get]
func (ctrl *Controller) getAllOrderProducts(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id := c.Param("id")

	uc := order.NewOrderUseCase(ctrl.repos.Orders)
	o, err := uc.GetByID(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if o.UserID != userId {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	orders, err := uc.GetAllOrderProducts(o.ID)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /orders/:id/products/:productID [delete]
func (ctrl *Controller) deleteOrderProduct(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id := c.Param("id")
	productID := c.Param("productID")

	uc := order.NewOrderUseCase(ctrl.repos.Orders)

	o, err := uc.GetByID(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if o.UserID != userId {
		newErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	if err := uc.RemoveProduct(id, productID); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{true})
}
