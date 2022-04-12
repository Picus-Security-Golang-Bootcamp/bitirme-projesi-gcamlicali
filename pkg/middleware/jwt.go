package mw

import (
	jwtHelper "github.com/gcamlicali/tradeshopExample/pkg/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "" {
			decodedClaims, err := jwtHelper.VerifyToken(c.GetHeader("Authorization"), secretKey)
			//if decodedClaims != nil {
			//	for _, role := range decodedClaims.Roles {
			//		if role == "admin" {
			//			c.Next()
			//			c.Abort()
			//			return
			//		}
			//	}
			//}
			if decodedClaims == nil {
				c.JSON(http.StatusUnauthorized, gin.H{"Authorization error": err.Error()})
				c.Abort()
				return
			}
			c.Set("userId", decodedClaims.UserId)
			c.Set("isAdmin", decodedClaims.IsAdmin)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"Authorization error": "You are not authorized!"})
		}
		c.Abort()
		return
	}
}
