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


// HasPermission checks if a user has a specific permission.
func (r *UserRepository) HasPermission(ctx context.Context, userID uint, permissionName string) (bool, error) {
	var count int64

	// Join with user_roles and role_permissions tables to check for permission
	err := r.db.WithContext(ctx).
		Table("users").
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Joins("JOIN role_permissions ON user_roles.role_id = role_permissions.role_id").
		Joins("JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("users.id = ? AND permissions.name = ? AND users.deleted_at IS NULL", userID, permissionName).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}