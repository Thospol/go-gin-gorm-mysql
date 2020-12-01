package product

import (
	"encoding/json"
	"go-gin-gorm-mysql/internal/core/config"
	"go-gin-gorm-mysql/internal/core/database"

	"github.com/gin-gonic/gin"
)

// Endpoint endpoint product interface
type Endpoint interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
}

type endpoint struct {
	config  *config.Configs
	result  *config.ReturnResult
	service Service
}

// NewEndpoint new endpoint product
func NewEndpoint(config *config.Configs, result *config.ReturnResult) Endpoint {
	return &endpoint{
		config:  config,
		result:  result,
		service: NewService(config, result),
	}
}

// Create godoc
// @Tags Product
// @Summary Create Product
// @Description Create Product Service API
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)"
// @Param request body product.createRequest true "request body"
// @Success 200 {object} models.Product
// @Failure 400 {object} config.SwaggerInfoResult
// @Security ApiKeyAuth
// @Router /v1/products [post]
func (ep *endpoint) Create(c *gin.Context) {
	var request createRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		c.AbortWithStatusJSON(ep.result.Internal.BadRequest.HTTPStatusCode(), ep.result.Internal.BadRequest.WithLocale(c))
		return
	}

	defer func() {
		c.Request.Body.Close()
	}()

	response, err := ep.service.Create(database.Database, request)
	if err != nil {
		errMsg := config.RR.Internal.ConnectionError
		if locErr, ok := err.(config.Result); ok {
			errMsg = locErr
		}
		c.AbortWithStatusJSON(errMsg.HTTPStatusCode(), errMsg)
		return
	}

	c.JSON(ep.result.Internal.Success.HTTPStatusCode(), response)
}

// GetAll godoc
// @Tags Product
// @Summary Get Products
// @Description Get Products Service API
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)"
// @Success 200 {array} models.Product
// @Failure 400 {object} config.SwaggerInfoResult
// @Router /v1/products [get]
func (ep *endpoint) GetAll(c *gin.Context) {
	response, err := ep.service.GetAll(database.Database)
	if err != nil {
		errMsg := config.RR.Internal.ConnectionError
		if locErr, ok := err.(config.Result); ok {
			errMsg = locErr
		}
		c.AbortWithStatusJSON(errMsg.HTTPStatusCode(), errMsg)
		return
	}

	c.JSON(ep.result.Internal.Success.HTTPStatusCode(), response)
}
