package router

import (
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitParcelRouter(
	deps RouterDeps,
	r *gin.RouterGroup,
) {
	// No auth routes (public)
	noAuthRouter := r.Group("/")
	{
		noAuthRouter.POST("/parcels/pickup", deps.ParcelHandler.PickupParcel)
		noAuthRouter.GET("/parcels/by-pickup-code/:code", deps.ParcelHandler.GetParcelByPickupCode)
		noAuthRouter.GET("/parcels/by-tracking/:number", deps.ParcelHandler.GetParcelByTrackingNumber)
		noAuthRouter.GET("/parcels", deps.ParcelHandler.GetParcelList)
		noAuthRouter.GET("/dashboard/statistics", deps.ParcelHandler.GetDashboardStatistics)
	}

	// Staff/Admin routes (require login + staff or admin role)
	staffRouter := r.Group("/").Use(
		middleware.StrictAuth(deps.JWT, deps.Logger),
		middleware.RequireRole("staff", "admin"),
	)
	{
		staffRouter.POST("/parcels/receive", deps.ParcelHandler.ReceiveParcel)
		staffRouter.POST("/parcels/:id/shelve", deps.ParcelHandler.ShelveParcel)
	}

	// Strict auth routes (require login)
	strictAuthRouter := r.Group("/").Use(middleware.StrictAuth(deps.JWT, deps.Logger))
	{
		strictAuthRouter.GET("/my-parcels", deps.ParcelHandler.GetMyParcels)
	}
}
