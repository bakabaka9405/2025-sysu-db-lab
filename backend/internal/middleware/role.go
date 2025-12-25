package middleware

import (
	v1 "backend/api/v1"
	"backend/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole returns a middleware that checks if the user has one of the allowed roles.
// This middleware must be used after StrictAuth middleware.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, exists := ctx.Get("claims")
		if !exists {
			v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}

		userClaims, ok := claims.(*jwt.MyCustomClaims)
		if !ok {
			v1.HandleError(ctx, http.StatusUnauthorized, v1.ErrUnauthorized, nil)
			ctx.Abort()
			return
		}

		// Check if user's role is in allowed roles
		for _, role := range allowedRoles {
			if userClaims.Role == role {
				ctx.Next()
				return
			}
		}

		// User role not in allowed roles
		v1.HandleError(ctx, http.StatusForbidden, v1.ErrForbidden, nil)
		ctx.Abort()
	}
}
