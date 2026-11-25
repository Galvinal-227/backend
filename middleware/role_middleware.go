package middleware

import (
	"github.com/gin-gonic/gin"

	"pkl-attendance/backend/utils"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		roleVal, _ := c.Get("role")
		role, _ := roleVal.(string)
		if !allowed[role] {
			utils.Error(c, 403, "forbidden: insufficient role")
			c.Abort()
			return
		}
		c.Next()
	}
}