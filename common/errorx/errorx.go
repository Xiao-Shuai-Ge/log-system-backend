package errorx

import (
	"fmt"
)

// CodeError 基础错误结构
type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// CodeErrorResponse 基础错误响应
type CodeErrorResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// NewCodeError 创建一个包含错误码的错误
func NewCodeError(code int, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

// NewDefaultError 创建一个包含默认错误码的错误
func NewDefaultError(msg string) error {
	return NewCodeError(1001, msg)
}

const (
	CodeParamError = 400
	CodeAuthError  = 401
	CodeForbidden  = 403
	CodeNotFound   = 404
	CodeInternal   = 500
)

func (e *CodeError) Error() string {
	return fmt.Sprintf("错误代码:%d, 错误消息:%s", e.Code, e.Msg)
}

// Data 返回错误响应对象
func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: e.Code,
		Msg:  e.Msg,
	}
}
