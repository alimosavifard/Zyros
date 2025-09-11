package repositories

import (
	"context"
	"github.com/alimosavifard/zyros-backend/models"
	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, post *models.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *PostRepository) CreateWithTx(tx *gorm.DB, post *models.Post) error {
	return tx.Create(post).Error
}

func (r *PostRepository) GetByLang(ctx context.Context, lang string, postType string, page, limit int) ([]models.Post, error) {
	var posts []models.Post
	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).
		Where("lang = ? AND type = ? AND deleted_at IS NULL", lang, postType).
		Preload("User"). // Preload User برای نمایش username در frontend
		Offset(offset).Limit(limit).Find(&posts).Error
	return posts, err
}

func (r *PostRepository) FindByID(ctx context.Context, id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.WithContext(ctx).
		Preload("User"). // Preload User
		First(&post, id).Error
	return &post, err
}