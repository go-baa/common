package errors

import (
	"fmt"
)

// APIError API错误结构
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error ...
func (e *APIError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// New 创建一个新的错误结构
func New(code int, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}
