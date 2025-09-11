package models

import (
	"gorm.io/gorm"
	"time"
)


// User represents a user entity in the system.
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Roles    []Role `gorm:"many2many:user_roles;" json:"roles"`
}


type Role struct {
	gorm.Model
	Name        string       `gorm:"unique;not null"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}


type Permission struct {
    ID   uint   `gorm:"primaryKey" json:"id"`
    Name string `gorm:"not null;unique" json:"name"`
}


// UserRole represents the many-to-many relationship between users and roles.
type UserRole struct {
	UserID    uint           `gorm:"primaryKey" json:"user_id"`
	RoleID    uint           `gorm:"primaryKey" json:"role_id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type RolePermission struct {
    RoleID       uint           `gorm:"primaryKey" json:"role_id"`
    PermissionID uint           `gorm:"primaryKey" json:"permission_id"`
    CreatedAt time.Time `gorm:"-"` // غیرفعال کردن CreatedAt
    UpdatedAt time.Time `gorm:"-"` // غیرفعال کردن UpdatedAt
    DeletedAt gorm.DeletedAt `gorm:"-"` // غیرفعال کردن DeletedAt
}


type Post struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"not null" json:"title"`
	Content   string         `gorm:"not null" json:"content"`
	Type      string         `gorm:"not null" json:"type"` // "post" or "article"
	Lang      string         `gorm:"not null" json:"lang"` // "fa" or "en"
	ImageUrl  string         `gorm:"type:text" json:"imageUrl"` // اختیاری
	UserID    uint           `gorm:"not null" json:"user_id"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	User      User           `json:"user,omitempty"` // برای preload
}

// PostLike به عنوان جدول واسط برای لایک‌ها (حذف Like، فقط این نگه داشته شود)
type PostLike struct {
	UserID    uint           `gorm:"primaryKey" json:"user_id"`
	PostID    uint           `gorm:"primaryKey" json:"post_id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	User      User           `gorm:"foreignKey:UserID"`
	Post      Post           `gorm:"foreignKey:PostID"`
}