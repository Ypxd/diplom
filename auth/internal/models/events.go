package models

type Events struct {
	Title   string `json:"title"`
	Address string `json:"address"`
	Tags    string `json:"tags"`
	PNG     string `json:"png"`
}

type EventsResponse struct {
	AllEvents int      `json:"all_events"`
	AllTags   int      `json:"all_tags"`
	Events    []Events `json:"events"`
}

type MyEvents struct {
	Title   string `json:"title"`
	Address string `json:"address"`
	Tags    string `json:"tags"`
	PNG     string `json:"png"`
	Val     int    `json:"val"`
}
