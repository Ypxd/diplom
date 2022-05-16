package models

type AllTags struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Val   int    `json:"val,omitempty"`
}

type TagsRep struct {
	ID    int64 `json:"id"`
	Count int64 `json:"count"`
}
