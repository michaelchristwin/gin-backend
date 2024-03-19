package main

import (
	"net/http"

	"example.com/gin/database"
	"example.com/gin/model"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDatabase()
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))
	r.GET("ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.POST("/stocks", model.AddStock)
	r.GET("/stocks", model.GetStock)
	r.PUT("/stocks/:id", model.UpdateStock)
	r.DELETE("/stocks/:id", model.DeleteStock)
	r.POST("/login", model.Login)
	protected := r.Group("/admin", model.AuthMiddleware)
	protected.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the admin dashboard",
		})
	})
	r.Run(":8080")

}
