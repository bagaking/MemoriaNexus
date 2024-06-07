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

type Updater[T any] struct {
	From    T              `json:"from,omitempty"`
	To      T              `json:"to,omitempty"`
	Updates map[string]any `json:"updates,omitempty"`
}

// RespSuccessPage defines the response structure for a successful operation.
type RespSuccessPage[T any] struct {
	Message string `json:"message"`
	Data    []T    `json:"data,omitempty"`

	Offset int   `json:"offset"`
	Limit  int   `json:"limit"`
	Total  int64 `json:"total,omitempty"`
}

// Page start from 1
func (resp *RespSuccessPage[T]) Page() int {
	return resp.Offset/resp.Limit + 1
}

func (resp *RespSuccessPage[T]) SetTotal(total int64) *RespSuccessPage[T] {
	resp.Total = total
	return resp
}

func (resp *RespSuccessPage[T]) SetPageAndLimit(page int, limit int) *RespSuccessPage[T] {
	resp.Offset = (page - 1) * resp.Limit
	resp.Limit = limit
	return resp
}

func (resp *RespSuccessPage[T]) SetOffsetAndLimit(offset int, limit int) *RespSuccessPage[T] {
	resp.Offset = offset
	resp.Limit = limit
	return resp
}

func (resp *RespSuccessPage[T]) Append(items ...T) *RespSuccessPage[T] {
	lenItems := len(items)
	if resp.Data == nil {
		// deep copy
		resp.Data = make([]T, lenItems)
		copy(resp.Data, items)
	} else {
		resp.Data = append(resp.Data, items...)
	}
	return resp
}

func (resp *RespSuccessPage[T]) Response(c *gin.Context, msgAppend ...string) {
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
