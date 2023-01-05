package v1_integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	v1 "github.com/linkuha/test-golang-rest-orders-api/internal/delivery/httpserver/v1"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/errs"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository"
	mockProducts "github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository/product/mocks"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProductByID(t *testing.T) {
	reqID := "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	exp := &entity.Product{
		ID:          reqID,
		Name:        "milk",
		LeftInStock: 1,
		Prices: []entity.Price{
			{
				Price:    "1.0",
				Currency: "USD",
			},
		},
	}

	ctx := context.Background()

	// Create dummy repos, for don't use other
	repos := repository.Repository{}

	repoProducts := mockProducts.NewMockRepository(ctrl)
	repoProducts.EXPECT().Get(ctx, reqID).Return(exp, nil).Times(1)

	repos.Products = repoProducts
	handler := v1.NewController(ctx, repos)

	// Init endpoint
	r := gin.New()
	r.GET("/products/:id", handler.GetProductByID)

	// Create request
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodGet,
		"/products/"+reqID,
		nil,
	)

	// Make request
	r.ServeHTTP(rec, req)

	data := rec.Body.String()

	expected :=
		`{"id":"c401f9dc-1e68-4b44-82d9-3a93b09e3fe7","name":"milk","description":"","left_in_stock":1,"prices":[{"currency":"USD","price":"1.0"}]}`

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, expected, data)
}

func TestGetProductServiceError(t *testing.T) {
	reqID := "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	// Create dummy repos, for don't use other
	repos := repository.Repository{}

	repoProducts := mockProducts.NewMockRepository(ctrl)
	repoProducts.EXPECT().Get(ctx, reqID).Return(nil, errs.HandleErrorDB(sql.ErrConnDone)).Times(1)

	repos.Products = repoProducts
	handler := v1.NewController(ctx, repos)

	// Init endpoint
	r := gin.New()
	r.GET("/products/:id", handler.GetProductByID)

	// Create request
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodGet,
		"/products/"+reqID,
		nil,
	)

	// Make request
	r.ServeHTTP(rec, req)

	data := rec.Body.String()

	expected :=
		"{\"ok\":false,\"message\":\"" + v1.ErrServiceUnavailableText + "\"}"

	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
	require.Equal(t, expected, data)
}

func TestCreateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	exp := "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"
	product := &entity.Product{
		Name:        "milk",
		LeftInStock: 1,
		Prices: []entity.Price{
			{
				Price:    "1.0",
				Currency: "USD",
			},
		},
	}
	productJson, err := json.Marshal(product)
	if err != nil {
		t.Errorf("Can't encode json request: %s", err.Error())
	}

	ctx := context.Background()

	// Create dummy repos, for don't use other
	repos := repository.Repository{}

	repoProducts := mockProducts.NewMockRepository(ctrl)
	repoProducts.EXPECT().StoreWithPrices(ctx, product).Return(exp, nil).Times(1)

	repos.Products = repoProducts
	handler := v1.NewController(ctx, repos)

	// Init endpoint
	r := gin.New()
	r.POST("/products", handler.CreateProduct)

	// Create request
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		bytes.NewBuffer(productJson),
	)

	// Make request
	r.ServeHTTP(rec, req)

	data := rec.Body.String()

	expected := "{\"id\":\"" + exp + "\"}"

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, expected, data)
}

func TestCreateProductServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	product := &entity.Product{
		Name:        "milk",
		LeftInStock: 1,
		Prices: []entity.Price{
			{
				Price:    "1.0",
				Currency: "USD",
			},
		},
	}
	productJson, err := json.Marshal(product)
	if err != nil {
		t.Errorf("Can't encode json request: %s", err.Error())
	}

	ctx := context.Background()

	// Create dummy repos, for don't use other
	repos := repository.Repository{}

	repoProducts := mockProducts.NewMockRepository(ctrl)
	repoProducts.EXPECT().StoreWithPrices(ctx, product).Return("", errs.HandleErrorDB(sql.ErrConnDone)).Times(1)

	repos.Products = repoProducts
	handler := v1.NewController(ctx, repos)

	// Init endpoint
	r := gin.New()
	r.POST("/products", handler.CreateProduct)

	// Create request
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		bytes.NewBuffer(productJson),
	)

	// Make request
	r.ServeHTTP(rec, req)

	data := rec.Body.String()

	expected :=
		"{\"ok\":false,\"message\":\"" + v1.ErrServiceUnavailableText + "\"}"

	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
	require.Equal(t, expected, data)
}

func TestCreateProductBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productJson := `{"name":"milk","left_in_stock":1,"prices":["price":"1.0"","currency":"usd"]`

	// Create dummy repos, for don't use other
	repos := repository.Repository{}
	handler := v1.NewController(context.Background(), repos)

	// Init endpoint
	r := gin.New()
	r.POST("/products", handler.CreateProduct)

	// Create request
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		bytes.NewBufferString(productJson),
	)

	// Make request
	r.ServeHTTP(rec, req)

	data := rec.Body.String()

	expected :=
		"{\"ok\":false,\"message\":\"" + v1.ErrInputJSONText + "\"}"

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, expected, data)
}

func TestCreateProductValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productJson := `{"name":"milk","left_in_stock":-1,"prices":[{"price":"1.0","currency":"usd"}]}`

	// Create dummy repos, for don't use other
	repos := repository.Repository{}
	handler := v1.NewController(context.Background(), repos)

	// Init endpoint
	r := gin.New()
	r.POST("/products", handler.CreateProduct)

	// Create request
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		bytes.NewBufferString(productJson),
	)

	// Make request
	r.ServeHTTP(rec, req)

	data := rec.Body.String()

	expected :=
		"{\"ok\":false,\"message\":\"" + v1.ErrValidationText + "\"}"

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	require.Equal(t, expected, data)
}

func TestCreateProductInnerValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productJson := `{"name":"milk","left_in_stock":1,"prices":[{"price":"1.0","currency":"us"}]}`

	// Create dummy repos, for don't use other
	repos := repository.Repository{}
	handler := v1.NewController(context.Background(), repos)

	// Init endpoint
	r := gin.New()
	r.POST("/products", handler.CreateProduct)

	// Create request
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/products",
		bytes.NewBufferString(productJson),
	)

	// Make request
	r.ServeHTTP(rec, req)

	data := rec.Body.String()

	expected :=
		"{\"ok\":false,\"message\":\"" + v1.ErrValidationText + "\"}"

	require.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	require.Equal(t, expected, data)
}
