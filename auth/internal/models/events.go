package models

type EventsRequest struct {
	JWTToken string   `json:"jwt_token"`
	Tags     []string `json:"tags,omitempty"`
}

type EventsResponse struct {
	Title   string `json:"title"`
	Address string `json:"address"`
	PNG     []byte `json:"png"`
}

type Events struct {
	Title   string `json:"title"`
	Address string `json:"address"`
	PNG     string `json:"png"`
}
