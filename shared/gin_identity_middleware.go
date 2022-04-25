package shared

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type (
	UserInfo struct {
		Exist      bool       `json:"exist,omitempty"`
		Email      string     `json:"email,omitempty"`
		FullName   string     `json:"full_name,omitempty"`
		Groups     []string   `json:"groups,omitempty"`
		Prefix     string     `json:"prefix,omitempty"`
		IsExternal bool       `json:"is_external,omitempty"`
		UserID     uuid.UUID  `json:"user_id,omitempty"`
		OrgID      *uuid.UUID `json:"org_id,omitempty"`
	}

	Auth struct {
		debug      bool
		cookieName string
		rdb        *redis.Client
	}
)

var (
	orgID     = uuid.MustParse("123e4567-e89e-12d3-a456-426614174000")
	debugUser = UserInfo{
		Email:    "test@name.ru",
		FullName: "test name",
		Prefix:   "orgPrefix",
		Groups:   []string{"dp_admin"},
		UserID:   uuid.MustParse("123e4567-e89e-12d3-a456-426614174000"),
		OrgID:    &orgID,
	}
)

func (a *Auth) GetUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var identity *UserInfo
		if a.debug {
			identity = &debugUser
			c.Set("identity", identity)
			c.Next()
			return
		}
		var (
			sessionID string
			err       error
		)

		sessionID = c.GetHeader(a.cookieName)
		if sessionID == "" {
			sessionID, err = c.Cookie(a.cookieName)
			if err != nil {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}

		identity, err = a.getUserFromRequest(c.Request.Context(), sessionID)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if identity != nil {
			c.Set("identity", identity)
			c.Next()
			return
		}
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func (a *Auth) getUserFromRequest(ctx context.Context, session string) (*UserInfo, error) {
	var userInfo UserInfo

	resp, err := a.rdb.Get(ctx, session).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(resp), &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func NewJWT(cookieName string, debug bool, rdb *redis.Client) *Auth {
	return &Auth{
		debug:      debug,
		cookieName: cookieName,
		rdb:        rdb,
	}
}
