package models

import "github.com/dgrijalva/jwt-go"

type AuthReq struct {
	Login    *string `json:"login"`
	Password *string `json:"password"`
	Email    *string `json:"email"`
	Name     *string `json:"name"`
	Age      *string `json:"age"`
}

type UserInfo struct {
	Login  string `json:"login"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Age    string `json:"age"`
	FTags  string `json:"f_tags" db:"f_tags"`
	UFTags string `json:"uf_tags" db:"uf_tags"`
}

type ChangePassReq struct {
	OldPassword *string `json:"old_password"`
	NewPassword *string `json:"new_password"`
}

type JWTToken struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
}
