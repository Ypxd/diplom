package models

type HttpResponse struct {
	ErrorText string      `json:"error_text"`
	HasError  bool        `json:"has_error"`
	Message   interface{} `json:"message"`
	Count     *int        `json:"count,omitempty"`
}
