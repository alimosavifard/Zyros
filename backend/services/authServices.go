package services

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/alimosavifard/zyros-backend/config"
	"github.com/alimosavifard/zyros-backend/models"
	"github.com/alimosavifard/zyros-backend/repositories"
	"github.com/alimosavifard/zyros-backend/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo    *repositories.UserRepository
	roleRepo    *repositories.RoleRepository
	redisClient *redis.Client
	jwtSecret   string
	jwtExp      time.Duration
}

func NewAuthService(userRepo *repositories.UserRepository, roleRepo *repositories.RoleRepository, redisClient *redis.Client, cfg *config.Config) *AuthService {
	jwtExp, err := time.ParseDuration(cfg.JWT_EXPIRATION)
	if err != nil {
		utils.InitLogger().Fatal().Err(err).Msg("Invalid JWT_EXPIRATION format")
	}

	return &AuthService{
		userRepo:    userRepo,
		roleRepo:    roleRepo,
		redisClient: redisClient,
		jwtSecret:   cfg.JWT_SECRET,
		jwtExp:      jwtExp,
	}
}

func (s *AuthService) Register(ctx context.Context, user *models.User) (string, error) {
	if s.jwtSecret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPassword)

	defaultRole, err := s.roleRepo.FindByName(ctx, "user")
	if err != nil {
		return "", errors.New("default role not found")
	}
	user.Roles = append(user.Roles, *defaultRole)

	if err := s.userRepo.Create(ctx, user); err != nil {
		return "", err
	}
	return s.generateToken(user.ID, user.Username)
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	if s.jwtSecret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}
	
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", utils.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", utils.ErrInvalidCredentials
	}

	return s.generateToken(user.ID, user.Username)
}

func (s *AuthService) generateToken(userID uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"userID":   userID,
		"username": username,
		"exp":      time.Now().Add(s.jwtExp).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["userID"].(float64)
		if !ok {
			return 0, errors.New("invalid userID")
		}
		return uint(userID), nil
	}
	return 0, errors.New("invalid token")
}

func (s *AuthService) HasPermission(ctx context.Context, userID uint, permission string) (bool, error) {
	cacheKey := s.getPermissionsCacheKey(userID)
	cachedData, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var data map[string]interface{}
		if json.Unmarshal([]byte(cachedData), &data) == nil {
			if perms, ok := data["permissions"].([]interface{}); ok {
				for _, perm := range perms {
					if p, ok := perm.(string); ok && p == permission {
						return true, nil
					}
				}
			}
		}
	}

	hasPermission, err := s.userRepo.HasPermission(ctx, userID, permission)
	if err != nil {
		return false, err
	}
	return hasPermission, nil
}

func (s *AuthService) getPermissionsCacheKey(userID uint) string {
	return "user_permissions:" + strconv.FormatUint(uint64(userID), 10)
}