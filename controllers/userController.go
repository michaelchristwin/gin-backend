package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"example.com/gin/database"
	"example.com/gin/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtSecret = []byte(os.Getenv("AUTHO_KEY"))
)

func Login(ctx *gin.Context) {
	var user model.User

	// Directly bind JSON to the user struct for efficiency
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Use prepared statement for security and clarity
	stmt, err := database.Db.Prepare("SELECT id, username, password FROM users WHERE username = $1")
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		println("failed to create statement", err.Error())
		return
	}
	defer stmt.Close()                  // Ensure proper resource cleanup
	row := stmt.QueryRow(user.Username) // Use username as query parameter
	var userID int
	var username string
	var hashedPassword string
	err = row.Scan(&userID, &username, &hashedPassword)
	if err != nil {
		fmt.Println("Failed to scan row", err.Error())
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Authentication successful, proceed with actions for logged-in user
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Set token expiration
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		println("Failed to get tokenString", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie("Authorization", tokenString, 3600*24, "/", "https://stock-manager-z.vercel.app/", true, true)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cookie has been set",
	})

}

func Signup(c *gin.Context) {
	var body struct {
		Username string
		Password string
	}
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
	}
	user := model.User{Username: body.Username, Password: string(hash)}

	stmt, err := database.Db.Prepare("INSERT INTO users(username, password) VALUES ($1, $2)")
	if err != nil {
		fmt.Println("failed to prepare sql statment", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to prepare sql stmt ",
		})
		return

	}
	defer stmt.Close()
	_, err = stmt.Exec(user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
