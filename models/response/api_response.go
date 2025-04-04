package response

// APIResponse represents a standardized structure for API responses containing data, errors, status, and messages.
type APIResponse struct {
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
}
