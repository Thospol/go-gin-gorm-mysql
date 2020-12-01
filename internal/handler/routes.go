package handler

import (
	"net/http"

	"go-gin-gorm-mysql/internal/core/config"
	"go-gin-gorm-mysql/internal/pkg/healthcheck"

	"github.com/gin-gonic/gin"
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
	r.v1 = []route{
		{
			Name:        "Healthcheck",
			Description: "response success message",
			Method:      http.MethodGet,
			Pattern:     "healthcheck",
			Endpoint:    healthcheckEndpoint.HealthCheck,
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

// Healthcheck health check service
func Healthcheck(c *gin.Context) {
	c.JSON(config.RR.Internal.Success.HTTPStatusCode(), config.RR.Internal.Success)
}
