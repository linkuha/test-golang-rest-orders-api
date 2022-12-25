package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/linkuha/test-golang-rest-orders-api/config"
	_ "github.com/linkuha/test-golang-rest-orders-api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

// ConfigureRoutes -.
// Swagger spec:
// @title       Go Orders API Test Issue
// @description REST API example
// @version     1.0
// @host        localhost:3000
// @BasePath    /v1
func (ctrl *Controller) ConfigureRoutes(cfg *config.Config) *gin.Engine {
	router := gin.New()

	if cfg.EnvParams.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Swagger
	//docs.SwaggerInfo.Host = "localhost:" + cfg.EnvParams.Port
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	router.GET("/swagger/*any", swaggerHandler)

	// Load status
	router.GET("/status", ctrl.status)

	// K8s probe
	router.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	v1 := router.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/sign-up", ctrl.signUp)
			auth.POST("/sign-in", ctrl.signIn)
		}

		api := v1.Group("/api", ctrl.userIdentity)
		{
			profile := api.Group("/profile")
			{
				profile.POST("/", ctrl.getProfile)
				profile.GET("/", ctrl.getProfile)
				profile.PUT("/", ctrl.updateProfile)
			}

			products := api.Group("/products")
			{
				products.POST("/", ctrl.createProduct)
				products.GET("/", ctrl.getAllProducts)
				products.GET("/:id", ctrl.getProductByID)
				products.PUT("/:id", ctrl.updateProductByID)
				products.DELETE("/:id", ctrl.deleteProductByID)
			}

			orders := api.Group("/orders")
			{
				orders.POST("/", ctrl.createOrder)
				orders.GET("/", ctrl.getAllOrders)
				orders.GET("/:id", ctrl.getOrderByID)
				orders.PUT("/:id", ctrl.updateOrderByID)
				orders.DELETE("/:id", ctrl.deleteOrderByID)

				orderProducts := orders.Group(":id/products")
				{
					orderProducts.POST("/", ctrl.addOrderProduct)
					orderProducts.GET("/", ctrl.getAllOrders)
					orderProducts.DELETE("/:productID", ctrl.deleteOrderProduct)
				}
			}
		}
	}

	return router
}
