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
		fmt.Println("tokenstring not found", err.Error())
		return
	}
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		fmt.Println("Parsing failed", err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if token == nil {
		// Token is nil, redirect to the home page ("/") with status code 302 (Found)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return

		}
		var user model.User
		stmt, err := database.Db.Prepare("SELECT id, username, password FROM users WHERE id = $1")
		if err != nil {
			wemo := fmt.Errorf("statement creation failed: %s", err.Error())
			fmt.Println(wemo)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer stmt.Close()
		err = stmt.QueryRow(claims["user_id"]).Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			fmt.Println("Query row failed", err.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Set("user", user)
		c.Next()

	} else {
		c.AbortWithStatus(http.StatusInternalServerError)
		fmt.Println("claims doesn't exist")
		return
	}
}
