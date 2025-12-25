package router

import (
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitShelfRouter(
	deps RouterDeps,
	r *gin.RouterGroup,
) {
	// Public routes (no auth required for viewing)
	noAuthRouter := r.Group("/")
	{
		noAuthRouter.GET("/shelves", deps.ShelfHandler.GetShelfList)
		noAuthRouter.GET("/shelves/statistics", deps.ShelfHandler.GetShelfStatistics)
		noAuthRouter.GET("/shelves/:id", deps.ShelfHandler.GetShelf)
	}

	// Staff/Admin routes (require login + staff or admin role)
	staffRouter := r.Group("/").Use(
		middleware.StrictAuth(deps.JWT, deps.Logger),
		middleware.RequireRole("staff", "admin"),
	)
	{
		staffRouter.POST("/shelves", deps.ShelfHandler.CreateShelf)
		staffRouter.PUT("/shelves/:id", deps.ShelfHandler.UpdateShelf)
		staffRouter.DELETE("/shelves/:id", deps.ShelfHandler.DeleteShelf)
	}
}
