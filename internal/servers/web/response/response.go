package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
)

// Response 错误信息
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// HTTPError example
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// Ok 返回成功
// @msg: 返回消息
// @data: 返回成功的数据
func Ok(c *gin.Context, data ...interface{}) {
	resp := Response{
		Code:    http.StatusOK,
		Message: "success",
	}

	if len(data) > 0 {
		resp.Data = data[0]
	}

	c.JSON(http.StatusOK, resp)
}

// Err 返回失败
// @param: httpCode 错误码
// @param: msg 错误消息
// @param: data 额外信息（可选）
func Err(c *gin.Context, httpCode int, msg string, data ...interface{}) {

	resp := Response{
		Code:    httpCode,
		Message: msg,
	}

	if len(data) > 0 {
		resp.Data = data[0]
	}

	c.JSON(httpCode, resp)
}

func ErrFromSwagger(c *gin.Context, httpCode int, msg string) {

	httputil.NewError(c, httpCode, errors.New(msg))
}
