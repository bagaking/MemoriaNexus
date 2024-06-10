package utils

import "fmt"

type Pager struct {
	Offset int   `json:"offset"`
	Limit  int   `json:"limit"`
	Total  int64 `json:"total,omitempty"`
}

func (resp *Pager) String() string {
	return fmt.Sprintf("pager(%d,%d)", resp.Offset, resp.Limit)
}

// Page start from 1
func (resp *Pager) Page() int {
	return resp.Offset/resp.Limit + 1
}

func (resp *Pager) SetTotal(total int64) *Pager {
	resp.Total = total
	return resp
}

func (resp *Pager) SetPageAndLimit(page, limit int) *Pager {
	resp.Limit = limit
	resp.Offset = (page - 1) * limit
	return resp
}

func (resp *Pager) SetOffsetAndLimit(offset int, limit int) *Pager {
	resp.Offset = offset
	resp.Limit = limit
	return resp
}

func (resp *Pager) SetFirstCount(count int) *Pager {
	resp.Limit = count
	return resp
}
