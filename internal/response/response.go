package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{Code: 0, Message: "success", Data: data})
}

func Error(c *gin.Context, httpCode int, bizCode int, message string) {
	c.JSON(httpCode, Response{Code: bizCode, Message: message, Data: nil})
}
