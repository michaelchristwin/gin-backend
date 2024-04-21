package main

import (
	"example.com/gin/controllers"
	"example.com/gin/database"
	"example.com/gin/middleware"
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
	protected := r.Group("/admin")
	protected.Use(middleware.RequireAuth)
	protected.GET("/")
	protected.POST("/stocks", model.AddStock)
	protected.GET("/stocks", model.GetStock)
	protected.PUT("/stocks/:id", model.UpdateStock)
	protected.DELETE("/stocks/:id", model.DeleteStock)
	r.POST("/login", controllers.Login)
	r.POST("/signup", controllers.Signup)
	r.Run(":3001")

}
