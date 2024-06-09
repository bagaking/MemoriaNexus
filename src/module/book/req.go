package book

import "github.com/bagaking/memorianexus/internal/utils"

type (

	// ReqGetBookItemsQuery encapsulates the request parameters for fetching items.
	ReqGetBookItemsQuery struct {
		Page  int `form:"page"`
		Limit int `form:"limit"`
	}

	ReqAddItems struct {
		ItemIDs []utils.UInt64 `json:"item_ids"`
	}

	ReqCreateOrUpdateBook struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Tags        []string `json:"tags,omitempty"`
	}
)
