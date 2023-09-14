package middleware

import (
	"fmt"
	"jwt/go/initializers"
	"jwt/go/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if float64(time.Now().Unix()) > exp {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user models.User
		result := initializers.DB.First(&user, claims["sub"])
		if result.Error != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", user)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}
