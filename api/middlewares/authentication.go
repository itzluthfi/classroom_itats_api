package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := jwt.Parse(c.GetHeader("token"), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
			}

			return []byte(viper.GetString("SECRET_KEY")), nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed", "error": err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if jwt.NewNumericDate(time.Unix((int64(claims["exp"].(float64))), 0)).After(time.Now()) {
				c.Next()
				return
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "User Token Expired"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Login First Before Access This Feature"})
			c.Abort()
			return
		}
	}
}

func ApiKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("x_api_key") == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Unauthorized Request"})
			c.Abort()
			return
		}

		if viper.GetString("X_API_KEY") == c.GetHeader("x_api_key") {
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "User Key Missmatch"})
		c.Abort()
		return
	}
}
