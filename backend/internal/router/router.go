package router

import (
	"backend/internal/handler"
	"backend/pkg/jwt"
	"backend/pkg/log"
	"github.com/spf13/viper"
)

type RouterDeps struct {
	Logger          *log.Logger
	Config          *viper.Viper
	JWT             *jwt.JWT
	UserHandler     *handler.UserHandler
	ParcelHandler   *handler.ParcelHandler
	ShipmentHandler *handler.ShipmentHandler
	ShelfHandler    *handler.ShelfHandler
	OverdueHandler  *handler.OverdueHandler
}
