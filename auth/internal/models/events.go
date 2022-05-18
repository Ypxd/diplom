package models

type Events struct {
	Title    string `json:"title"`
	Address  string `json:"address"`
	Tags     string `json:"tags"`
	PNG      string `json:"png"`
	Age      int64  `json:"age_id" db:"age_id"`
	Selected string `json:"selected"`
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
