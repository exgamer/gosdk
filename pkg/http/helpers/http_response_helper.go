package helpers

import (
	"errors"
	"fmt"
	"github.com/exgamer/gosdk/pkg/exception"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

func FormattedTextErrorResponse(c *gin.Context, statusCode int, message string, context map[string]any) {
	TextErrorResponse(c, statusCode, message, context)
	FormattedResponse(c)
}

func TextErrorResponse(c *gin.Context, statusCode int, message string, context map[string]any) {
	AppExceptionResponse(c, exception.NewAppException(statusCode, errors.New(message), context))
}

func FormattedErrorResponse(c *gin.Context, statusCode int, err error, context map[string]any) {
	ErrorResponse(c, statusCode, err, context)
	FormattedResponse(c)
}

func ErrorResponse(c *gin.Context, statusCode int, err error, context map[string]any) {
	AppExceptionResponse(c, exception.NewAppException(statusCode, err, context))
}

func FormattedAppExceptionResponse(c *gin.Context, exception *exception.AppException) {
	AppExceptionResponse(c, exception)
	FormattedResponse(c)
}

func AppExceptionResponse(c *gin.Context, exception *exception.AppException) {
	c.Set("exception", exception)
	c.Status(exception.Code)
}

func SuccessResponse(c *gin.Context, data any) {
	c.Set("data", data)
}

func SuccessCreatedResponse(c *gin.Context, data any) {
	c.Set("data", data)
	c.Set("status_code", http.StatusCreated)
}

func SuccessDeletedResponse(c *gin.Context, data any) {
	c.Set("data", data)
	c.Set("status_code", http.StatusNoContent)
}

func FormattedSuccessResponse(c *gin.Context, data any) {
	SuccessResponse(c, data)
	FormattedResponse(c)
}

func FormattedResponse(c *gin.Context) {
	for _, err := range c.Errors {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "data": err.Error(), "service_code": 0})

		return
	}

	appExceptionObject, exists := c.Get("exception")
	fmt.Printf("%+v\n", appExceptionObject)

	if !exists {
		data, _ := c.Get("data")
		c.Writer.Status()

		statusCode, ex := c.Get("status_code")

		if !ex {
			c.JSON(http.StatusOK, gin.H{"success": true, "data": data, "service_code": 0})
		} else {
			c.JSON(statusCode.(int), gin.H{"success": true, "data": data, "service_code": 0})
		}

		return
	}

	appException := exception.AppException{}
	mapstructure.Decode(appExceptionObject, &appException)
	fmt.Printf("%+v\n", appException)

	c.JSON(appException.Code, gin.H{"success": false, "message": appException.Error.Error(), "details": appException.Context, "service_code": appException.ServiceCode})
}
