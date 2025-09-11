package migrations

import (
	"fmt"

	"github.com/alimosavifard/zyros-backend/models"
	"github.com/alimosavifard/zyros-backend/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RunMigrations executes database migrations and seeds initial data.
func RunMigrations(db *gorm.DB) error {
	// Drop all tables
	if err := db.Migrator().DropTable(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.Post{},
	); err != nil {
		return fmt.Errorf("failed to drop tables: %w", err)
	}
	utils.InitLogger().Info().Msg("All tables dropped successfully")

	// AutoMigrate tables
	if err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.Post{},
	); err != nil {
		return fmt.Errorf("failed to migrate tables: %w", err)
	}
	utils.InitLogger().Info().Msg("Database tables migrated successfully")

	// Seed admin user
	admin := &models.User{
		Username: "admin",
		Password: hashPassword("admin123"),
	}
	if err := db.Create(admin).Error; err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}
	utils.InitLogger().Info().Msg("Admin user seeded successfully")

	// Seed roles and permissions
	if err := seedRolesAndPermissions(db, admin); err != nil {
		return fmt.Errorf("failed to seed roles and permissions: %w", err)
	}

	return nil
}

// seedRolesAndPermissions seeds roles, permissions, and their relationships.
func seedRolesAndPermissions(db *gorm.DB, admin *models.User) error {
	// Seed permissions
	createPostPerm := &models.Permission{Name: "create_post"}
	if err := db.Where("name = ?", createPostPerm.Name).FirstOrCreate(createPostPerm).Error; err != nil {
		utils.InitLogger().Error().Err(err).Msg("Failed to seed create_post permission")
		return fmt.Errorf("failed to seed create_post permission: %w", err)
	}

	// Seed roles
	userRole := &models.Role{Name: "user"}
	if err := db.Where("name = ?", userRole.Name).FirstOrCreate(userRole).Error; err != nil {
		return fmt.Errorf("failed to seed user role: %w", err)
	}

	adminRole := &models.Role{Name: "admin"}
	if err := db.Where("name = ?", adminRole.Name).FirstOrCreate(adminRole).Error; err != nil {
		return fmt.Errorf("failed to seed admin role: %w", err)
	}

	// Assign permissions to roles
	if err := db.Model(userRole).Association("Permissions").Append(createPostPerm); err != nil {
		return fmt.Errorf("failed to assign create_post permission to user role: %w", err)
	}
	if err := db.Model(adminRole).Association("Permissions").Append(createPostPerm); err != nil {
		return fmt.Errorf("failed to assign create_post permission to admin role: %w", err)
	}

	// Assign roles to admin user
	if err := db.Model(admin).Association("Roles").Append([]*models.Role{userRole, adminRole}); err != nil {
		return fmt.Errorf("failed to assign roles to admin user: %w", err)
	}

	return nil
}

// hashPassword hashes a password using bcrypt.
func hashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.InitLogger().Error().Err(err).Msg("Failed to hash password")
	}
	return string(hashed)
}