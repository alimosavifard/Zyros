package controllers

import (
	"github.com/alimosavifard/zyros-backend/models"
	"github.com/alimosavifard/zyros-backend/requests"
	"github.com/alimosavifard/zyros-backend/services"
	"github.com/alimosavifard/zyros-backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req requests.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if err := req.Validate(); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Validation failed", err)
		return
	}

	user := &models.User{
		Username: req.Username,
		Password: req.Password,
	}

	token, err := c.authService.Register(user)
	if err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Failed to register", err)
		return
	}

	// Set JWT token as an HTTP-only cookie
	ctx.SetCookie("token", token, 3600*24, "/", "", true, true)
	utils.SendSuccess(ctx, "Registration successful", nil, nil)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req requests.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Invalid input", err)
		return
	}

	if err := req.Validate(); err != nil {
		utils.SendError(ctx, http.StatusBadRequest, "Validation failed", err)
		return
	}

	token, err := c.authService.Login(&req)
	if err != nil {
		if err.Error() == utils.ErrInvalidCredentials.Error() {
			utils.SendError(ctx, http.StatusUnauthorized, "Invalid username or password", nil)
			return
		}
		utils.SendError(ctx, http.StatusInternalServerError, "Failed to login", err)
		return
	}

	// Set JWT token as an HTTP-only cookie
	ctx.SetCookie("token", token, 3600*24, "/", "", true, true)
	utils.SendSuccess(ctx, "Login successful", nil, nil)
}

func (c *AuthController) GetCSRFToken(ctx *gin.Context) {
	// The CSRF token is already set in an HTTP-only cookie by the CSRF middleware.
	// You can just send a success response.
	utils.SendSuccess(ctx, "CSRF token is set in cookie", nil, nil)
}