package repositories

import (
	"github.com/alimosavifard/zyros-backend/models"
	"gorm.io/gorm"
	"context"
)

type LikeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
	return &LikeRepository{db: db}
}

func (r *LikeRepository) Create(ctx context.Context, like *models.PostLike) error {
	return r.db.WithContext(ctx).Create(like).Error
}

func (r *LikeRepository) Delete(ctx context.Context, userID, postID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND post_id = ?", userID, postID).Delete(&models.PostLike{}).Error
}

func (r *LikeRepository) IsLiked(ctx context.Context, userID, postID uint) (bool, error) {
	var like models.PostLike
	err := r.db.WithContext(ctx).Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return err == nil, err
}

func (r *LikeRepository) CountLikes(ctx context.Context, postID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.PostLike{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}

func (r *LikeRepository) FindPostsByUserID(ctx context.Context, userID uint) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.WithContext(ctx).
		Joins("JOIN post_likes ON posts.id = post_likes.post_id").
		Where("post_likes.user_id = ?", userID).
		Find(&posts).Error
	return posts, err
}