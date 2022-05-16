package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Ypxd/diplom/auth/internal/models"
	"github.com/Ypxd/diplom/auth/internal/repository"
	"github.com/Ypxd/diplom/auth/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type AuthService struct {
	repo *repository.Repository
	conn *sqlx.DB
}

func (a *AuthService) Auth(ctx context.Context, req models.AuthReq) (string, error) {
	userID, err := a.repo.Auth.Auth(ctx, req)
	if err != nil {
		return "", err
	}

	mes, err := NewJWT(utils.GetConfig().Auth.SessionTime, utils.GetConfig().Auth.TokenSecret, userID.String())
	if err != nil {
		return "", err
	}
	return mes, nil
}

func (a *AuthService) ChangePass(ctx context.Context, req models.ChangePassReq, userID string) error {
	return a.repo.Auth.ChangePass(ctx, req, userID)
}

func (a *AuthService) Register(ctx context.Context, req models.AuthReq) (string, error) {
	age, err := a.repo.Auth.Age(ctx, req)
	if err != nil {
		return "", err
	}

	m, err := a.repo.Auth.Register(ctx, req, age)
	if err != nil {
		return m, err
	}

	return m, nil
}

func (a *AuthService) UserInfo(ctx context.Context, userID string) (models.UserInfo, error) {
	usr, err := a.repo.Auth.UserInfo(ctx, userID)
	if err != nil {
		return models.UserInfo{}, err
	}

	if usr.FTags != "" {
		ft := strings.Split(usr.FTags, ";")
		if len(ft) == 0 {
			usr.FTags = ""
		} else {
			fTags, err := a.repo.Tags.AllUnfavoriteTagsTags(ctx, ft)
			if err != nil {
				return models.UserInfo{}, err
			}
			s := ""
			for _, t := range fTags {
				if s == "" {
					s = s + t.Title
				} else {
					s = s + ", " + t.Title
				}
			}
			s = s + "."
			usr.FTags = s
		}
	}

	if usr.UFTags != "" {
		uft := strings.Split(usr.UFTags, ";")
		if len(uft) == 0 {
			usr.UFTags = ""
		} else {
			ufTags, err := a.repo.Tags.AllUnfavoriteTagsTags(ctx, uft)
			if err != nil {
				return models.UserInfo{}, err
			}
			s := ""
			for _, t := range ufTags {
				if s == "" {
					s = s + t.Title
				} else {
					s = s + ", " + t.Title
				}
			}
			s = s + "."
			usr.UFTags = s
		}
	}

	return *usr, nil
}

func NewJWT(sessionTime time.Duration, secretKey string, u string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.JWTToken{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(sessionTime * time.Hour).UTC().Unix(),
		},
		UserID: u,
	})

	tokenJWT, err := token.SignedString([]byte(secretKey))

	return tokenJWT, err
}

func ParseJWT(secretKey string, tokenS string) (models.JWTToken, string, error) {
	var (
		jwtRes models.JWTToken
	)

	token, err := jwt.Parse(tokenS, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if !token.Valid {
		return jwtRes, "", errors.New("invalid token")
	}
	if err != nil {
		return jwtRes, "", err
	}

	tByte, err := json.Marshal(token.Claims)
	if err != nil {
		return jwtRes, "", err
	}
	err = json.Unmarshal(tByte, &jwtRes)
	if err != nil {
		return jwtRes, "", err
	}

	if time.Now().UTC().Unix()+((time.Hour*utils.GetConfig().Auth.SessionTime).Milliseconds()/1000)/2 >= jwtRes.ExpiresAt {
		tokenS, err = NewJWT(utils.GetConfig().Auth.SessionTime, secretKey, jwtRes.UserID)
	}

	return jwtRes, tokenS, err
}

func NewAuthService(repo *repository.Repository, conn *sqlx.DB) *AuthService {
	return &AuthService{
		repo: repo,
		conn: conn,
	}
}
