package main

import (
	"gin-backend/internal/config"

	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	// Catch-all route
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Route not found",
			"path":  ctx.Request.URL.Path,
		})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router.Handler(),
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
