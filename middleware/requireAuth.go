package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"example.com/gin/database"
	"example.com/gin/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret = []byte(os.Getenv("AUTHO_KEY"))
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		var user model.User
		stmt, err := database.Db.Prepare("SELECT id, username, password FROM users WHERE id = ?")
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		defer stmt.Close()
		err = stmt.QueryRow(claims["user_id"]).Scan(&user)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		c.Set("user", user)
		c.Next()

	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
	}

}
