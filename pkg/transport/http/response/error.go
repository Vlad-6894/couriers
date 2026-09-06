package pkg_http_response

type ErrorResponse struct {
	Error   string `json:"Error"         example:"full error text"`
	Message string `json:"Message"       example:"message about error"`
}
