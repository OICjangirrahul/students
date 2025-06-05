// Package main provides the entry point for the Student-Teacher Management API
package main

import (
	"log/slog"
	"os"

	_ "github.com/OICjangirrahul/students/docs" // Import swagger docs
	"github.com/OICjangirrahul/students/internal"
	"github.com/OICjangirrahul/students/internal/config"
	"github.com/OICjangirrahul/students/internal/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Student-Teacher Management API
// @version         1.0
// @description     A Go-based REST API for managing students and teachers.
// @host            localhost:8082
// @BasePath        /api/v1

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/local.yaml")
	if err != nil {
		slog.Error("failed to load configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Initialize handlers
	handlers, err := internal.InitializeAppHandlers(cfg)
	if err != nil {
		slog.Error("failed to initialize handlers", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Initialize Gin router
	r := gin.Default()

	// Add CORS middleware
	r.Use(middleware.CorsMiddleware())

	// API v1 group
	v1 := r.Group("/api/v1")

	// Teacher routes
	teachers := v1.Group("/teachers")
	{
		teachers.POST("", handlers.Teacher.Create())
		teachers.POST("/login", handlers.Teacher.Login())

		teacherManagement := teachers.Group("/:id")
		{
			teacherManagement.GET("", handlers.Teacher.GetByID())
			teacherManagement.PUT("", handlers.Teacher.Update())
			teacherManagement.DELETE("", handlers.Teacher.Delete())

			studentManagement := teacherManagement.Group("/students")
			{
				studentManagement.GET("", handlers.Teacher.GetStudents())
				studentManagement.POST("/:studentId", handlers.Teacher.AssignStudent())
			}
		}
	}

	// Student routes
	students := v1.Group("/students")
	{
		students.POST("", handlers.Student.Create())
		students.POST("/login", handlers.Student.Login())
		students.GET("/:id", handlers.Student.GetByID())
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	if err := r.Run(cfg.HTTPServer.Addr); err != nil {
		slog.Error("failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
