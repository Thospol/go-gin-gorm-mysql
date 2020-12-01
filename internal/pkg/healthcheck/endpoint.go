package healthcheck

import (
	"go-gin-gorm-mysql/internal/core/config"

	"github.com/gin-gonic/gin"
)

// Endpoint endpoint health check interface
type Endpoint interface {
	HealthCheck(c *gin.Context)
}

type endpoint struct {
	config *config.Configs
	result *config.ReturnResult
}

// NewEndpoint new endpoint health check
func NewEndpoint(config *config.Configs, result *config.ReturnResult) Endpoint {
	return &endpoint{
		config: config,
		result: result,
	}
}

// HealthCheck godoc
// @Tags Health
// @Summary HealthCheck
// @Description HealthCheck Service API
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)"
// @Success 200 {object} config.SwaggerInfoResult
// @Failure 400 {object} config.SwaggerInfoResult
// @Router /v1/healthcheck [get]
func (ep *endpoint) HealthCheck(c *gin.Context) {
	c.JSON(config.RR.Internal.Success.HTTPStatusCode(), ep.result.Internal.Success.WithLocale(c))
}
