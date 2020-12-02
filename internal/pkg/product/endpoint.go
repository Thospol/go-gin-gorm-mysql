package product

import (
	"encoding/json"
	"go-gin-gorm-mysql/internal/core/config"
	"go-gin-gorm-mysql/internal/core/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Endpoint endpoint product interface
type Endpoint interface {
	Create(c *gin.Context)
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
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

	if err := config.CF.Validator.Struct(request); err != nil {
		c.AbortWithStatusJSON(
			ep.result.Internal.BadRequest.HTTPStatusCode(),
			config.RR.CustomMessage(err.Error(), err.Error()).WithLocale(c),
		)
		return
	}

	response, err := ep.service.Create(database.Get(), request)
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
	response, err := ep.service.GetAll(database.Get())
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

// GetByID godoc
// @Tags Product
// @Summary Get Product
// @Description Get Product Service API
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)"
// @Param id path int true "query by product_id"
// @Success 200 {object} models.Product
// @Failure 400 {object} config.SwaggerInfoResult
// @Router /v1/products/{id} [get]
func (ep *endpoint) GetByID(c *gin.Context) {
	value := c.Param("id")
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(ep.result.InvalidRequest.HTTPStatusCode(), ep.result.InvalidRequest.WithLocale(c))
		return
	}

	response, err := ep.service.GetByID(database.Get(), uint(id))
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

// Update godoc
// @Tags Product
// @Summary Update Product
// @Description Update Product Service API
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)"
// @Param id path int true "query by product_id"
// @Param request body product.createRequest true "request body"
// @Success 200 {object} models.Product
// @Failure 400 {object} config.SwaggerInfoResult
// @Router /v1/products/{id} [patch]
func (ep *endpoint) Update(c *gin.Context) {
	var request updateRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		c.AbortWithStatusJSON(ep.result.Internal.BadRequest.HTTPStatusCode(), ep.result.Internal.BadRequest.WithLocale(c))
		return
	}

	defer func() {
		c.Request.Body.Close()
	}()

	if err := config.CF.Validator.Struct(request); err != nil {
		c.AbortWithStatusJSON(
			ep.result.Internal.BadRequest.HTTPStatusCode(),
			config.RR.CustomMessage(err.Error(), err.Error()).WithLocale(c),
		)
		return
	}

	value := c.Param("id")
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(ep.result.InvalidRequest.HTTPStatusCode(), ep.result.InvalidRequest.WithLocale(c))
		return
	}

	response, err := ep.service.Update(database.Get(), uint(id), request)
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

// Delete godoc
// @Tags Product
// @Summary Delete Product
// @Description Delete Product Service API
// @Accept json
// @Produce json
// @Param Accept-Language header string false "(en, th)"
// @Param id path int true "query by product_id"
// @Success 200 {object} config.SwaggerInfoResult
// @Failure 400 {object} config.SwaggerInfoResult
// @Router /v1/products/{id} [delete]
func (ep *endpoint) Delete(c *gin.Context) {
	value := c.Param("id")
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(ep.result.InvalidRequest.HTTPStatusCode(), ep.result.InvalidRequest.WithLocale(c))
		return
	}

	err = ep.service.Delete(database.Get(), uint(id))
	if err != nil {
		errMsg := config.RR.Internal.ConnectionError
		if locErr, ok := err.(config.Result); ok {
			errMsg = locErr
		}
		c.AbortWithStatusJSON(errMsg.HTTPStatusCode(), errMsg)
		return
	}

	c.JSON(ep.result.Internal.Success.HTTPStatusCode(), ep.result.Internal.Success.WithLocale(c))
}
