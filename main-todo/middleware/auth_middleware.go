// middleware/auth.go
package middleware

import (
	"main-todo/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get token from header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			utils.HandleError(c, 401, "Authorization header is missing")
			return
		}

		// 2. Parse and validate token
		token, err := utils.ParseToken(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte("your-256-bit-secret"), nil // Replace with your actual secret
		})

		if err != nil {
			utils.HandleError(c, 401, "Invalid token: "+err.Error())
			return
		}

		// 3. Verify claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			utils.HandleError(c, 401, "Invalid token claims")
			return
		}

		// 4. Set user ID in context
		if userID, exists := claims["user_id"]; exists {
			c.Set("user_id", userID)
		} else {
			utils.HandleError(c, 401, "Missing user_id in token")
			return
		}

		c.Next()
	}
}
