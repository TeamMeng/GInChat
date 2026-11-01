package utils

import (
	"errors"
	"net/http"

	"GinChat/logger"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "ok",
		Data: data,
	})
}

func SuccessWithMsg(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  msg,
	})
}

func Error(ctx *gin.Context, code int, msg string) {
	ctx.JSON(httpStatusOf(code), Response{
		Code: code,
		Msg:  msg,
	})
}

func Fail(ctx *gin.Context, err error) {
	if biz, ok := errors.AsType[*BizError](err); ok {
		Error(ctx, 400, biz.Msg)
		return
	}
	logger.Errorf("[API] %s %s Internal Server Error: %v", ctx.Request.Method, ctx.Request.URL.Path, err)
	Error(ctx, 500, "Internal Server Error")
}

func httpStatusOf(code int) int {
	switch {
	case code == CodeSuccess:
		return http.StatusOK
	case code >= 40100 && code <= 40199:
		return http.StatusUnauthorized
	case code >= 400 && code <= 499:
		return http.StatusBadRequest
	case code >= 500:
		return http.StatusInternalServerError
	default:
		return http.StatusOK
	}
}
