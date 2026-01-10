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
			// Verify token (accessible to all authenticated users)
			protected.GET("/verify", handleVerify)

			// Users management (accessible to all authenticated users)
			protected.GET("/users/search", handleSearchUsers)
			protected.GET("/users/:id", handleGetUser)
			protected.PUT("/users/:id/role", handleSetUserRole)
		}

		// Full access routes (blocked for users_only role)
		fullAccess := api.Group("")
		fullAccess.Use(authMiddleware())
		fullAccess.Use(requireFullAccess())
		{
			// Statistics
			fullAccess.GET("/statistics", handleGetStatistics)

			// Cities CRUD
			fullAccess.GET("/cities", handleGetCities)
			fullAccess.POST("/cities", handleAddCity)
			fullAccess.PUT("/cities/:id", handleUpdateCity)
			fullAccess.DELETE("/cities/:id", handleDeleteCity)

			// Schools CRUD
			fullAccess.GET("/schools", handleGetSchools)
			fullAccess.POST("/schools", handleAddSchool)
			fullAccess.PUT("/schools/:id", handleUpdateSchool)
			fullAccess.DELETE("/schools/:id", handleDeleteSchool)

			// Posts moderation
			fullAccess.GET("/posts", handleGetAllPosts)
			fullAccess.GET("/posts/pending", handleGetPendingPosts)
			fullAccess.GET("/posts/reported", handleGetReportedPosts)
			fullAccess.PUT("/posts/:id/approve", handleApprovePost)
			fullAccess.PUT("/posts/:id/reject", handleRejectPost)
			fullAccess.PUT("/posts/:id/status", handleSetPostStatus)
			fullAccess.DELETE("/posts/:id", handleDeletePost)

			// Spotted moderation
			fullAccess.GET("/spotted", handleGetAllSpotted)
			fullAccess.GET("/spotted/pending", handleGetPendingSpotted)
			fullAccess.GET("/spotted/reported", handleGetReportedSpotted)
			fullAccess.PUT("/spotted/:id/approve", handleApproveSpotted)
			fullAccess.PUT("/spotted/:id/reject", handleRejectSpotted)
			fullAccess.PUT("/spotted/:id/status", handleSetSpottedStatus)
			fullAccess.DELETE("/spotted/:id", handleDeleteSpotted)
		}
	}

	fmt.Printf("Moderator Dashboard running on http://%s:%d\n", CONF.HOST, CONF.PORT)
	fmt.Printf("Dashboard UI: http://%s:%d/dashboard\n", CONF.HOST, CONF.PORT)
	router.Run(fmt.Sprintf("%s:%d", CONF.HOST, CONF.PORT))
}
