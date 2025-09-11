package repositories

import (
	"context"
	"gorm.io/gorm"
	"github.com/alimosavifard/zyros-backend/models"
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
	err := r.db.WithContext(ctx).Where("lang = ?", lang).Where("type = ?", postType).Offset(offset).Limit(limit).Find(&posts).Error
	return posts, err
}

// تابع جدید برای پیدا کردن پست با شناسه
func (r *PostRepository) FindByID(id uint) (*models.Post, error) {
    var post models.Post
    err := r.db.First(&post, id).Error
    if err != nil {
        return nil, err
    }
    return &post, nil
}