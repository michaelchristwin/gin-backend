package main

import (
	"database/sql"
	"gin-backend/internal/config"
	sqlc "gin-backend/internal/db"
	"gin-backend/internal/handler"
	"gin-backend/internal/middleware"
	"gin-backend/internal/repository"
	"gin-backend/internal/service"
	migrator "gin-backend/sqlc"
	"path/filepath"

	_ "modernc.org/sqlite" // pure-Go sqlite driver (no CGO needed)

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
	db, err := sql.Open("sqlite", cfg.DBDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := os.MkdirAll(filepath.Dir(cfg.DBPath), 0755); err != nil {
		log.Fatalf("creating db directory: %v", err)
	}
	if err := migrator.RunMigrations(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	queries := sqlc.New(db)
	userRepo := repository.NewUserRepository(queries)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	sessionRepo := repository.NewSessionRepository(queries)
	sessionService := service.NewSessionService(sessionRepo)
	authService := service.NewAuthService(userRepo, sessionRepo)
	authHandler := handler.NewAuthHandler(userService, authService)
	authMiddleware := middleware.NewAuthMiddleware(sessionService, userService)

	router := gin.Default()
	router.Use(middleware.ErrorHandler())
	api := router.Group("/api/v1")

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Gin Server")
	})
	userHandler.RegisterRoutes(api, authMiddleware)
	authHandler.RegisterRoutes(api, authMiddleware)
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
