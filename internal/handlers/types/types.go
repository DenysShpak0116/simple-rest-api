package types

type ApiResponse[T any] struct {
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
