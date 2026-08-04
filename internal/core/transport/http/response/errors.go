package core_transport_http_response

type ErrorResponse struct {
	Error   string `json:"error" example:"full error text"`
	Message string `json:"message" example:"short msg"`
}
