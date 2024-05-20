package module

type (
	// ErrorResponse defines the standard error response structure.
	ErrorResponse struct {
		Message string `json:"message"`
	}

	// SuccessResponse defines the response structure for a successful operation.
	SuccessResponse struct {
		Message string `json:"message"`
	}
)
