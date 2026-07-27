package pkg_http_response

type ErrorResponse struct {
	Error   string `json:"Error"`
	Message string `json:"Message"`
}
