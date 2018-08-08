package util

import (
	"errors"

	"github.com/hohice/gin-web/server/ex"

	"github.com/gin-gonic/gin"
)

func GetPathParams(c *gin.Context, names []string) (values map[string]string, err error) {
	for _, name := range names {
		values[name] = c.Param(name)
	}
	for _, value := range values {
		if len(value) == 0 {
			err = errors.New("")
			c.JSON(ex.ReturnBadRequest())
			break
		}
	}
	return
}
