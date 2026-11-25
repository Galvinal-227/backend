package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"pkl-attendance/backend/config"
	"pkl-attendance/backend/repositories"
	"pkl-attendance/backend/utils"
)

func Auth(cfg *config.Config, tokens *repositories.TokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			utils.Error(c, 401, "missing or invalid authorization header")
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := utils.ParseToken(cfg.JWTSecret, tokenStr)
		if err != nil {
			utils.Error(c, 401, "invalid or expired token")
			c.Abort()
			return
		}
		if tokens.IsRevoked(claims.ID) {
			utils.Error(c, 401, "token revoked")
			c.Abort()
			return
		}
		var expUnix int64
		if claims.ExpiresAt != nil {
			expUnix = claims.ExpiresAt.Unix()
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("jti", claims.ID)
		c.Set("exp", expUnix)
		c.Next()
	}
}