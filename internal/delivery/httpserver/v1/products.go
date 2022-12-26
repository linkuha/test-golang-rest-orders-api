package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/usecase/product"
	"net/http"
)

type createProductInput struct {
	Product entity.Product
	Prices  []entity.Price
}

// @Summary Create product
// @Security ApiKeyAuth
// @Tags product
// @Description Create product
// @ID product-create
// @Accept  json
// @Produce  json
// @Param input body createProductInput true "product data"
// @Success 200 {string} string "id"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /products [post]
func (ctrl *Controller) createProduct(c *gin.Context) {
	var input createProductInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	uc := product.NewProductUseCase(ctrl.repos.Products)
	id, err := uc.CreateWithPrices(input.Product, input.Prices)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Get product
// @Security ApiKeyAuth
// @Tags product
// @Description get product
// @ID product-get
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.Product
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /products/:id [get]
func (ctrl *Controller) getProductByID(c *gin.Context) {
	id := c.Param("id")

	uc := product.NewProductUseCase(ctrl.repos.Products)
	p, err := uc.GetByID(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
	products, err := uc.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
// @Param input body entity.ProductUpdateInput true "product updating data"
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /products/:id [put]
func (ctrl *Controller) updateProductByID(c *gin.Context) {
	var input entity.ProductUpdateInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id := c.Param("id")

	uc := product.NewProductUseCase(ctrl.repos.Products)
	if err := uc.Update(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
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
// @Success 200 {object} statusResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /products/:id [delete]
func (ctrl *Controller) deleteProductByID(c *gin.Context) {
	id := c.Param("id")

	uc := product.NewProductUseCase(ctrl.repos.Products)
	if err := uc.Remove(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{true})
}
