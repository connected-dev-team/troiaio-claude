package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func makeDbaseConnection() (*pgx.Conn, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		CONF.DBASE_USER,
		CONF.DBASE_PASSWD,
		CONF.DBASE_HOST,
		CONF.DBASE_PORT,
		CONF.DBASE_NAME,
	)
	return pgx.Connect(context.Background(), connStr)
}

func corsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

func main() {
	router := gin.Default()

	// Enable CORS
	router.Use(corsMiddleware())

	// Serve static files (dashboard)
	router.Static("/dashboard", "./dashboard")
	router.StaticFile("/", "./dashboard/index.html")

	// API routes
	api := router.Group("/api")
	{
		// Public routes (login)
		api.POST("/login", handleLogin)

		// Protected routes (require auth)
		protected := api.Group("")
		protected.Use(authMiddleware())
		{
			// Verify token
			protected.GET("/verify", handleVerify)

			// Cities CRUD
			protected.GET("/cities", handleGetCities)
			protected.POST("/cities", handleAddCity)
			protected.PUT("/cities/:id", handleUpdateCity)
			protected.DELETE("/cities/:id", handleDeleteCity)

			// Schools CRUD
			protected.GET("/schools", handleGetSchools)
			protected.POST("/schools", handleAddSchool)
			protected.PUT("/schools/:id", handleUpdateSchool)
			protected.DELETE("/schools/:id", handleDeleteSchool)

			// Posts moderation
			protected.GET("/posts", handleGetAllPosts)
			protected.GET("/posts/pending", handleGetPendingPosts)
			protected.GET("/posts/reported", handleGetReportedPosts)
			protected.PUT("/posts/:id/approve", handleApprovePost)
			protected.PUT("/posts/:id/reject", handleRejectPost)
			protected.PUT("/posts/:id/status", handleSetPostStatus)
			protected.DELETE("/posts/:id", handleDeletePost)

			// Spotted moderation
			protected.GET("/spotted", handleGetAllSpotted)
			protected.GET("/spotted/pending", handleGetPendingSpotted)
			protected.GET("/spotted/reported", handleGetReportedSpotted)
			protected.PUT("/spotted/:id/approve", handleApproveSpotted)
			protected.PUT("/spotted/:id/reject", handleRejectSpotted)
			protected.PUT("/spotted/:id/status", handleSetSpottedStatus)
			protected.DELETE("/spotted/:id", handleDeleteSpotted)
		}
	}

	fmt.Printf("Moderator Dashboard running on http://%s:%d\n", CONF.HOST, CONF.PORT)
	fmt.Printf("Dashboard UI: http://%s:%d/dashboard\n", CONF.HOST, CONF.PORT)
	router.Run(fmt.Sprintf("%s:%d", CONF.HOST, CONF.PORT))
}
