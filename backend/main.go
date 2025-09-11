package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/alimosavifard/zyros-backend/config"
	"github.com/alimosavifard/zyros-backend/controllers"
	"github.com/alimosavifard/zyros-backend/middleware"
	"github.com/alimosavifard/zyros-backend/repositories"
	"github.com/alimosavifard/zyros-backend/services"
	"github.com/alimosavifard/zyros-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func main() {
	if err := godotenv.Load(); err != nil {
		utils.InitLogger().Fatal().Err(err).Msg("Error loading .env file")
	}

	// Load all configs from environment variables
	cfg := config.NewConfig()

	logger := utils.InitLogger()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	
	
	// Connect to DB and Redis using the config
	db := config.ConnectDB(cfg)
	redisClient := config.ConnectRedis(cfg)

	// Initialize repositories with context-aware methods
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	postRepo := repositories.NewPostRepository(db)
	likeRepo := repositories.NewLikeRepository(db)

	// اصلاح ترتیب: likeService را اول تعریف کنید
	likeService := services.NewLikeService(likeRepo)
	authService := services.NewAuthService(userRepo, roleRepo, redisClient, cfg)
	postService := services.NewPostService(postRepo, redisClient, likeService) // حالا likeService تعریف شده
	
	
	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	postController := controllers.NewPostController(postService)
	articleController := controllers.NewArticleController(postService)
	likeController := controllers.NewLikeController(likeService)

	// Pass config values to middlewares
	r.Use(middleware.CORSMiddleware(cfg.ALLOWED_ORIGINS))
	r.Use(gin.Logger())
	r.Use(middleware.RateLimitMiddleware(redisClient, cfg.RATE_LIMIT))
	r.Use(gin.Recovery()) // جدید: برای مدیریت panic و جلوگیری از کرش سرور

	// CSRF middleware is now initialized with a secret
	r.GET("/api/v1/health", controllers.HealthCheck)
	r.POST("/api/v1/register", authController.Register)
	r.POST("/api/v1/login", authController.Login)
	r.GET("/api/v1/csrf-token", authController.GetCSRFToken)
	r.GET("/api/v1/posts", postController.GetPosts)
	r.GET("/api/v1/posts/:id", postController.GetPostByID)

	api := r.Group("/api/v1")
	api.Use(middleware.CSRFMiddleware(cfg.CSRF_SECRET), middleware.AuthMiddleware(authService))
	{
		api.POST("/posts", middleware.PermissionMiddleware(authService, "create_post"), postController.CreatePost)
		api.POST("/articles", middleware.PermissionMiddleware(authService, "create_article"), articleController.CreateArticle)
		api.POST("/upload-image", middleware.PermissionMiddleware(authService, "upload_image"), postController.UploadImage)
		api.POST("/posts/:id/like", middleware.PermissionMiddleware(authService, "like_post"), likeController.LikePost)
		api.DELETE("/posts/:id/like", middleware.PermissionMiddleware(authService, "unlike_post"), likeController.UnlikePost)
	}

	r.Static("/uploads", "./uploads")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":" + cfg.PORT)
}