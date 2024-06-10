package dto

import (
	"github.com/bagaking/memorianexus/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespSuccess defines the response structure for a successful operation.
type (
	RespSuccess[T any] struct {
		Message string `json:"message"`
		Data    T      `json:"data,omitempty"`
	}

	Updater[T any] struct {
		From    T              `json:"from,omitempty"`
		To      T              `json:"to,omitempty"`
		Updates map[string]any `json:"updates,omitempty"`
	}

	// RespSuccessPage defines the response structure for a successful operation.
	RespSuccessPage[T any] struct {
		Message string `json:"message"`
		Data    []T    `json:"data"`

		*utils.Pager
	}

	RespIDList = RespSuccessPage[utils.UInt64]
)

func (resp *RespSuccess[T]) With(t T) *RespSuccess[T] {
	resp.Data = t
	return resp
}

func (resp *RespSuccess[T]) Response(c *gin.Context, msgAppend ...string) *RespSuccess[T] {
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
	return resp
}

func (resp *RespSuccessPage[T]) WithPager(pager *utils.Pager) *RespSuccessPage[T] {
	resp.Pager = pager
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

func (resp *RespSuccessPage[T]) Response(c *gin.Context, msgAppend ...string) *RespSuccessPage[T] {
	for _, msg := range msgAppend {
		if resp.Message != "" {
			resp.Message += " "
		}
		resp.Message += msg
	}
	if resp.Message == "" {
		resp.Message = "success"
	}

	if resp.Data == nil {
		resp.Data = make([]T, 0)
	}

	// Respond with a generic success message.
	c.JSON(http.StatusOK, resp)
	return resp
}
