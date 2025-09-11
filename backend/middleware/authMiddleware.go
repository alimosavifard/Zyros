package middleware

import (
	"net/http"
	"strings"
	
	"github.com/alimosavifard/zyros-backend/services"
	"github.com/alimosavifard/zyros-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.SendError(ctx, http.StatusUnauthorized, "Authorization header required", nil)
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.SendError(ctx, http.StatusUnauthorized, "Invalid authorization header", nil)
			ctx.Abort()
			return
		}
		
		userID, err := authService.ValidateToken(ctx, parts[1])
		if err != nil {
			utils.SendError(ctx, http.StatusUnauthorized, "Invalid token", err)
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}

func CORSMiddleware(allowedOrigins string) gin.HandlerFunc {
    origins := strings.Split(allowedOrigins, ",")
    if len(origins) == 0 || origins[0] == "" {
        origins = []string{"https://news.asrnegar.ir"} // Default value
    }

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func PermissionMiddleware(authService *services.AuthService, permName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			utils.SendError(ctx, http.StatusUnauthorized, "Unauthorized", nil)
			ctx.Abort()
			return
		}

		hasPermission, err := authService.HasPermission(ctx, userID.(uint), permName)
		if err != nil {
			utils.SendError(ctx, http.StatusInternalServerError, "Failed to check permission", err)
			ctx.Abort()
			return
		}

		if !hasPermission {
			utils.SendError(ctx, http.StatusForbidden, "Forbidden", nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}