package config

import (
	"github.com/gin-gonic/gin"
)

const (
	// LangKey context key lang
	LangKey = "lang"
)

// Set value into context with key
func Set(c *gin.Context, key string, value interface{}) {
	c.Set(key, value)
}

// Get value from context with key
func Get(c *gin.Context, key string) interface{} {
	value, exists := c.Get(key)
	if !exists {
		return nil
	}

	return value
}
