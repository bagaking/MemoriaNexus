package book

import (
	"time"

	"github.com/bagaking/memorianexus/internal/utils"
)

// ReqCreateBook 定义创建书册的请求体
type ReqCreateBook struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type RespBook struct {
	ID          utils.UInt64 `json:"id"`
	UserID      utils.UInt64 `json:"user_id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type RespBooks struct {
	Books []RespBook `json:"books"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
}
