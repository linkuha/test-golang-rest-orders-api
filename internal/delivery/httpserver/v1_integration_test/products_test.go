package v1_integration_test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	v1 "github.com/linkuha/test-golang-rest-orders-api/internal/delivery/httpserver/v1"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/entity"
	"github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository"
	mockProducts "github.com/linkuha/test-golang-rest-orders-api/internal/domain/repository/product/mocks"
	"github.com/pkg/errors"
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
				Price:    1.0,
				Currency: "USD",
			},
		},
	}

	// Create dummy repos, for don't use other
	repos := repository.Repository{}

	repoProducts := mockProducts.NewMockRepository(ctrl)
	repoProducts.EXPECT().Get(reqID).Return(exp, nil).Times(1)

	repos.Products = repoProducts
	handler := v1.NewController(repos)

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
		`{"ID":"c401f9dc-1e68-4b44-82d9-3a93b09e3fe7","Name":"milk","Description":"","left_in_stock":1,"Prices":[{"Currency":"USD","Price":1}]}`

	require.Equal(t, rec.Code, http.StatusOK)
	require.Equal(t, expected, data)
}

func TestGetProductServiceError(t *testing.T) {
	reqID := "c401f9dc-1e68-4b44-82d9-3a93b09e3fe7"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create dummy repos, for don't use other
	repos := repository.Repository{}

	repoProducts := mockProducts.NewMockRepository(ctrl)
	repoProducts.EXPECT().Get(reqID).Return(nil, errors.New("db is down")).Times(1)

	repos.Products = repoProducts
	handler := v1.NewController(repos)

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
		"{\"ok\":false,\"message\":\"" + v1.ServiceError.Error() + "\"}"

	require.Equal(t, rec.Code, http.StatusInternalServerError)
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
				Price:    1.0,
				Currency: "USD",
			},
		},
	}
	productJson, err := json.Marshal(product)
	if err != nil {
		t.Errorf("Can't encode json request: %s", err.Error())
	}

	// Create dummy repos, for don't use other
	repos := repository.Repository{}

	repoProducts := mockProducts.NewMockRepository(ctrl)
	repoProducts.EXPECT().StoreWithPrices(product).Return(exp, nil).Times(1)

	repos.Products = repoProducts
	handler := v1.NewController(repos)

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

	require.Equal(t, rec.Code, http.StatusOK)
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
				Price:    1.0,
				Currency: "USD",
			},
		},
	}
	productJson, err := json.Marshal(product)
	if err != nil {
		t.Errorf("Can't encode json request: %s", err.Error())
	}

	// Create dummy repos, for don't use other
	repos := repository.Repository{}

	repoProducts := mockProducts.NewMockRepository(ctrl)
	repoProducts.EXPECT().StoreWithPrices(product).Return("", errors.New("db is down")).Times(1)

	repos.Products = repoProducts
	handler := v1.NewController(repos)

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
		"{\"ok\":false,\"message\":\"" + v1.ServiceError.Error() + "\"}"

	require.Equal(t, rec.Code, http.StatusInternalServerError)
	require.Equal(t, expected, data)
}

func TestCreateProductBadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	productJson := `{"name":"milk","left_in_stock":1,"prices":["price":1.0,"currency":"usd"]`

	// Create dummy repos, for don't use other
	repos := repository.Repository{}
	handler := v1.NewController(repos)

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
		"{\"ok\":false,\"message\":\"" + v1.MalformedRequest.Error() + "\"}"

	require.Equal(t, rec.Code, http.StatusBadRequest)
	require.Equal(t, expected, data)
}
