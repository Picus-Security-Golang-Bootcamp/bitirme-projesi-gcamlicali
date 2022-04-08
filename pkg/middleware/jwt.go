package mw

import (
	jwtHelper "github.com/gcamlicali/tradeshopExample/pkg/jwt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "" {
			decodedClaims := jwtHelper.VerifyToken(c.GetHeader("Authorization"), secretKey)
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
				log.Println("Bos Decoded")
				c.Abort()
				return
			}
			c.Set("userId", decodedClaims.UserId)
			c.Set("isAdmin", decodedClaims.IsAdmin)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized!"})
		}
		//log.Println("Aborta geldi")
		c.Abort()
		return
	}
}
