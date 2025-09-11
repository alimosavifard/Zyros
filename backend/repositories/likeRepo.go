// repositories/likeRepo.go

package repositories

import (
    "github.com/alimosavifard/zyros-backend/models"
    "gorm.io/gorm"
)

type LikeRepository struct {
    db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) *LikeRepository {
    return &LikeRepository{db: db}
}

func (r *LikeRepository) Create(like *models.PostLike) error {
    return r.db.Create(like).Error
}

func (r *LikeRepository) Delete(userID, postID uint) error {
    return r.db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&models.PostLike{}).Error
}

func (r *LikeRepository) IsLiked(userID, postID uint) (bool, error) {
    var like models.PostLike
    err := r.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error
    return err == nil, err
}

func (r *LikeRepository) CountLikes(postID uint) (int64, error) {
    var count int64
    err := r.db.Model(&models.PostLike{}).Where("post_id = ?", postID).Count(&count).Error
    return count, err
}

func (r *LikeRepository) FindPostsByUserID(userID uint) ([]models.Post, error) {
    var posts []models.Post
    err := r.db.Joins("JOIN post_likes ON posts.id = post_likes.post_id").
        Where("post_likes.user_id = ?", userID).
        Find(&posts).Error
    return posts, err
}