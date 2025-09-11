package repositories

import (
	"context"
	"github.com/alimosavifard/zyros-backend/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) FindByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) AssignRoleToUser(ctx context.Context, userID, roleID uint) error {
	userRole := models.UserRole{UserID: userID, RoleID: roleID}
	return r.db.WithContext(ctx).Create(&userRole).Error
}

func (r *RoleRepository) GetUserRoles(ctx context.Context, userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.WithContext(ctx).Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).Find(&roles).Error
	return roles, err
}