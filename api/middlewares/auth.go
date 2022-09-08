package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"tinderutf/internal/auth"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var header struct {
			Token string `header:"Authorization"`
		}

		if err := c.ShouldBindHeader(&header); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Auth header provided.",
			})
			return
		}

		tokenData := strings.Split(header.Token, "Bearer ")

		if len(tokenData) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Token provided.",
			})
			return
		}

		token := tokenData[1]

		user, err := auth.JWT.GetUserDataFromToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Token provided.",
			})
			return
		}

		c.Set("userId", user.Id)
		c.Next()
	}
}
