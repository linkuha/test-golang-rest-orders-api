package v1

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/linkuha/test-golang-rest-orders-api/config"
	docs "github.com/linkuha/test-golang-rest-orders-api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

// ConfigureRoutes -.
// Swagger spec:
// @title       Go Orders API Test Issue
// @description REST API example
// @version     1.0
// @BasePath    /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func (ctrl *Controller) ConfigureRoutes(cfg *config.Config) *gin.Engine {
	// Swagger
	docs.SwaggerInfo.Host = "localhost:" + cfg.EnvParams.Port

	router := gin.New()

	router.Use(requestid.New(), ctrl.customLogRequest)
	//router.Use(gin.LoggerWithWriter(log.Logger, "/status", "/healthz"))

	if cfg.EnvParams.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
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

		api := v1.Group("/", ctrl.userIdentity)
		{
			profile := api.Group("/profiles")
			{
				profile.POST("/my", ctrl.createMyProfile)
				profile.GET("/my", ctrl.getMyProfile)
				profile.GET("/:id", ctrl.getProfile)
				profile.PUT("/:id", ctrl.updateProfile) // e.g. under admin ACL
			}

			followers := api.Group("/followers")
			{
				followers.POST("/", ctrl.addFollower)
			}

			products := api.Group("/products")
			{
				products.POST("/", ctrl.CreateProduct)
				products.GET("/", ctrl.getAllProducts)
				products.GET("/:id", ctrl.GetProductByID)
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
					orderProducts.GET("/", ctrl.getAllOrderProducts)
					orderProducts.DELETE("/:productID", ctrl.deleteOrderProduct)
				}
			}
		}
	}

	return router
}
