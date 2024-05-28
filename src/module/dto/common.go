package dto

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
