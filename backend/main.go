package main

import (
	"lateslip/controllers"
	"lateslip/events"

	"lateslip/initialializers"
	"lateslip/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initialializers.LoadEnvVariables()
	initialializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	r.Use(middleware.RequestIDMiddleware())

	userRoutes := r.Group("/")
	{
		userRoutes.POST("/student/register", controllers.Register)
		userRoutes.POST("/student/login", controllers.Login)
		userRoutes.POST("/admin/register", controllers.AdminRegister)
		userRoutes.POST("/admin/login", controllers.AdminLogin)

	}

	// Add WebSocket routes
	studentRoutes := r.Group("/student").Use(middleware.AuthMiddleware(), middleware.RequireRole("student"))
	{
		studentRoutes.POST("/requestLateSlip", controllers.RequestLateSlip)
		// Replace SSE with WebSocket endpoint for students
		studentRoutes.GET("/ws", events.WebSocketHandler)
	}

	adminRoutes := r.Group("/admin").Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		adminRoutes.PUT("/lateslips/approve", controllers.ApproveLateSlip)
		adminRoutes.GET("/lateslips", controllers.GetAllLateSlips)
		adminRoutes.POST("/uploadStudentData", controllers.UploadStudentData)
		adminRoutes.GET("/lateslips/pending", controllers.GetAllPendingLateSlip)
		adminRoutes.PUT("/lateslips/reject", controllers.RejectLateSlip)
		adminRoutes.POST("/uploadScheduleData", controllers.UploadScheduleData)
		// Replace SSE with WebSocket endpoint for admins
		adminRoutes.GET("/ws", events.WebSocketHandler)
	}

	r.Run(":8000")
}