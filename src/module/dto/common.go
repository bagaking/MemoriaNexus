package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespSuccess defines the response structure for a successful operation.
type RespSuccess[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

// RespSuccessPage defines the response structure for a successful operation.
type RespSuccessPage[T any] struct {
	Message string `json:"message"`
	Data    []T    `json:"data,omitempty"`

	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

func (resp *RespSuccessPage[T]) Append(items ...T) int {
	resp.Total += int64(len(items))
	resp.Data = append(resp.Data, items...)
	return len(resp.Data)
}

func (resp *RespSuccess[T]) With(t T) *RespSuccess[T] {
	resp.Data = t
	return resp
}

func (resp *RespSuccess[T]) Response(c *gin.Context, msgAppend ...string) {
	for _, msg := range msgAppend {
		if resp.Message != "" {
			resp.Message += " "
		}
		resp.Message += msg
	}
	if resp.Message == "" {
		resp.Message = "success"
	}

	// Respond with a generic success message.
	c.JSON(http.StatusOK, resp)
}
