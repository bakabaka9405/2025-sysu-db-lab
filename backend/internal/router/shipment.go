package router

import (
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitShipmentRouter(
	deps RouterDeps,
	r *gin.RouterGroup,
) {
	// No auth routes (public)
	noAuthRouter := r.Group("/")
	{
		noAuthRouter.POST("/shipments/calculate-freight", deps.ShipmentHandler.CalculateFreight)
		noAuthRouter.GET("/courier-companies", deps.ShipmentHandler.GetCourierCompanies)
	}

	// Strict auth routes (require login)
	strictAuthRouter := r.Group("/").Use(middleware.StrictAuth(deps.JWT, deps.Logger))
	{
		strictAuthRouter.POST("/shipments", deps.ShipmentHandler.CreateShipment)
		strictAuthRouter.GET("/shipments/:id", deps.ShipmentHandler.GetShipmentDetail)
		strictAuthRouter.GET("/my-shipments", deps.ShipmentHandler.GetMyShipments)
	}
}
