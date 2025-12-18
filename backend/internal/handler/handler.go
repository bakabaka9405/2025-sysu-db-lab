package handler

import (
	"backend/internal/service"
	"backend/pkg/jwt"
	"backend/pkg/log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger      *log.Logger
	userService service.UserService
}

func NewHandler(
	logger *log.Logger,
	userService service.UserService,
) *Handler {
	return &Handler{
		logger:      logger,
		userService: userService,
	}
}
func GetUserIdFromCtx(ctx *gin.Context) int64 {
	v, exists := ctx.Get("claims")
	if !exists {
		return 0
	}
	userId := v.(*jwt.MyCustomClaims).UserId
	// UserId 在 JWT claims 中是字符串，需要转换为 int64
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return 0
	}
	return id
}
