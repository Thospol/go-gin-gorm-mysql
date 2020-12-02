package handler

import (
	"net/http"

	"go-gin-gorm-mysql/internal/core/config"
	"go-gin-gorm-mysql/internal/pkg/healthcheck"
	"go-gin-gorm-mysql/internal/pkg/product"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"   // gin-swagger middleware
	"github.com/swaggo/gin-swagger/swaggerFiles" // swagger embed files
)

type route struct {
	Name        string
	Description string
	Method      string
	Pattern     string
	Endpoint    gin.HandlerFunc
}

// Routes holds configurations related to API of this project
type Routes struct {
	v1 []route
}

// Init adds API routes to route object
func (r Routes) Init(config *config.Configs, result *config.ReturnResult) http.Handler {
	healthcheckEndpoint := healthcheck.NewEndpoint(config, result)
	productEndpoint := product.NewEndpoint(config, result)

	r.v1 = []route{
		{
			Name:        "Healthcheck",
			Description: "response success message",
			Method:      http.MethodGet,
			Pattern:     "healthcheck",
			Endpoint:    healthcheckEndpoint.HealthCheck,
		},
		{
			Name:        "Swagger UI",
			Description: "swagger ui",
			Method:      http.MethodGet,
			Pattern:     "swagger/*any",
			Endpoint:    ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("doc.json")),
		},
		{
			Name:        "Create",
			Description: "create product",
			Method:      http.MethodPost,
			Pattern:     "products",
			Endpoint:    productEndpoint.Create,
		},
		{
			Name:        "Get all",
			Description: "get all products",
			Method:      http.MethodGet,
			Pattern:     "products",
			Endpoint:    productEndpoint.GetAll,
		},
		{
			Name:        "Get by id",
			Description: "get product by id",
			Method:      http.MethodGet,
			Pattern:     "products/:id",
			Endpoint:    productEndpoint.GetByID,
		},
		{
			Name:        "Update",
			Description: "update product by id",
			Method:      http.MethodPatch,
			Pattern:     "products/:id",
			Endpoint:    productEndpoint.Update,
		},
		{
			Name:        "Delete",
			Description: "delete product by id",
			Method:      http.MethodDelete,
			Pattern:     "products/:id",
			Endpoint:    productEndpoint.Delete,
		},
	}

	route := gin.New()
	route.Use(Request())
	g1 := route.Group("api/v1")
	for _, e := range r.v1 {
		g1.Handle(e.Method, e.Pattern, e.Endpoint)
	}

	return route
}
