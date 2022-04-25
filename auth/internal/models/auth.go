package models

type AuthReq struct {
	Login    *string `json:"login"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	Name     *string `json:"name"`
	Age      *int    `json:"age"`
}
