package repositories

import (
	"context"
	"github.com/alimosavifard/zyros-backend/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.DB.WithContext(ctx).Preload("Roles.Permissions").Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.DB.WithContext(ctx).Preload("Roles.Permissions").First(&user, id).Error
	return &user, err
}