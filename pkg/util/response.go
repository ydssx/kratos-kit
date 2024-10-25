package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Reason string      `json:"reason"`
}

const (
	ERROR      = -1
	SUCCESS    = 0
	ErrorMsg   = "fail"
	SuccessMsg = "success"
)

func result(c *gin.Context, httpCode int, code int, data interface{}, msg, reason string) {
	c.JSON(httpCode, Response{
		Code:   code,
		Msg:    msg,
		Data:   data,
		Reason: reason,
	})
}

func OK(c *gin.Context) {
	result(c, http.StatusOK, SUCCESS, nil, SuccessMsg, "")
}

func OKWithData(c *gin.Context, data interface{}) {
	result(c, http.StatusOK, SUCCESS, data, SuccessMsg, "")
}

func FailWithMsg(c *gin.Context, msg string) {
	result(c, http.StatusOK, ERROR, nil, msg, "")
}

func FailWithError(c *gin.Context, err error) {
	se := errors.FromError(err)
	if se.Code != errors.UnknownCode {
		result(c, http.StatusOK, int(se.Code), nil, se.Message, se.Reason)
	} else {
		result(c, http.StatusOK, ERROR, nil, wrapValidateErrMsg(err), "")
	}
}

func wrapValidateErrMsg(err error) (msg string) {
	switch v := err.(type) {
	case *json.UnmarshalTypeError:
		msg = fmt.Sprintf("请求参数`%s`类型错误，应为%s类型", v.Field, v.Type.Name())
	case validator.ValidationErrors:
		for _, e := range v {
			msg += fmt.Sprintf("缺少必要参数：`%s`", strings.ToLower(e.Field()))
		}
	default:
		msg = err.Error()
	}
	return
}
