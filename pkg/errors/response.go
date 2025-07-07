package errors

// ErrorResponse стандартный формат ошибки API
// @description Используется для возврата ошибок клиенту
type ErrorResponse struct {
	// HTTP-код ошибки
	// @example 400
	Code int `json:"code"`

	// Сообщение об ошибке
	// @example "Invalid request parameters"
	Message string `json:"message"`

	// Детали ошибки (опционально)
	Details interface{} `json:"details,omitempty"`
}
