package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/usecase/product"
	"net/http"
)

// CreateProduct
// @Summary Create product
// @Security ApiKeyAuth
// @Tags product
// @Description Create product
// @ID product-create
// @Accept  json
// @Produce  json
// @Param input body entity.Product true "product data"
// @Success 200 {string} string "id"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /products [post]
func (ctrl *Controller) CreateProduct(c *gin.Context) {
	var input entity.Product

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, newJSONBindingErrorWrapper(err))
		return
	}

	uc := product.NewProductUseCase(ctrl.repos.Products)
	id, err := uc.CreateWithPrices(ctrl.ctx, input)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// GetProductByID
// @Summary Get product
// @Security ApiKeyAuth
// @Tags product
// @Description get product
// @ID product-get
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} entity.Product
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /products/{id} [get]
func (ctrl *Controller) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newErrorResponse(c, emptyParameterID)
		return
	}

	uc := product.NewProductUseCase(ctrl.repos.Products)
	p, err := uc.GetByID(ctrl.ctx, id)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, p)
}

// @Summary Get all products
// @Security ApiKeyAuth
// @Tags product
// @Description get all products
// @ID product-get-all
// @Accept  json
// @Produce  json
// @Success 200 {object} dataResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /products [get]
func (ctrl *Controller) getAllProducts(c *gin.Context) {
	uc := product.NewProductUseCase(ctrl.repos.Products)
	products, err := uc.GetAll(ctrl.ctx)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	newDataResponse(c, *products)
}

// @Summary Update product
// @Security ApiKeyAuth
// @Tags product
// @Description update product
// @ID product-update
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Param input body entity.ProductUpdateInput true "product updating data"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /products/{id} [put]
func (ctrl *Controller) updateProductByID(c *gin.Context) {
	var input entity.ProductUpdateInput

	id := c.Param("id")
	if id == "" {
		newErrorResponse(c, emptyParameterID)
		return
	}

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, newJSONBindingErrorWrapper(err))
		return
	}

	uc := product.NewProductUseCase(ctrl.repos.Products)
	if err := uc.Update(ctrl.ctx, id, input); err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{true})
}

// @Summary Delete product
// @Security ApiKeyAuth
// @Tags product
// @Description delete product
// @ID product-delete
// @Accept  json
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /products/{id} [delete]
func (ctrl *Controller) deleteProductByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		newErrorResponse(c, emptyParameterID)
		return
	}

	uc := product.NewProductUseCase(ctrl.repos.Products)
	if err := uc.Remove(ctrl.ctx, id); err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, statusResponse{true})
}
