package router

import (
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitOverdueRouter(
	deps RouterDeps,
	r *gin.RouterGroup,
) {
	// Public routes (viewing only)
	noAuthRouter := r.Group("/")
	{
		noAuthRouter.GET("/overdue", deps.OverdueHandler.GetOverdueList)
	}

	// Staff/Admin routes (require login + staff or admin role)
	staffRouter := r.Group("/").Use(
		middleware.StrictAuth(deps.JWT, deps.Logger),
		middleware.RequireRole("staff", "admin"),
	)
	{
		staffRouter.POST("/overdue/:id/remind", deps.OverdueHandler.SendReminder)
		staffRouter.POST("/overdue/batch-return", deps.OverdueHandler.BatchReturn)
		staffRouter.POST("/overdue/check", deps.OverdueHandler.CheckOverdue)
	}
}
